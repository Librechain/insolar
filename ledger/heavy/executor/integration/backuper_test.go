//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package intergration

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"
)

type testKey struct {
	id uint64
}

func (t *testKey) ID() []byte {
	bs := make([]byte, 8)
	binary.PutUvarint(bs, t.id)
	return bs
}

func (t *testKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func addLastBackupFile(t *testing.T, to string, lastBackupedVersion uint64) {
	backupInfo := executor.LastBackupInfo{
		LastBackupedVersion: lastBackupedVersion,
	}
	rawInfo, err := json.MarshalIndent(backupInfo, "", "    ")
	require.NoError(t, err)

	err = ioutil.WriteFile(to, rawInfo, 0600)
	println("============ addLastBackupFile: ", to)
	require.NoError(t, err)
}

func makeBackuperConfig(t *testing.T, prefix string, badgerDir string) configuration.Ledger {

	cwd, err := os.Getwd()
	if err != nil {
		require.NoError(t, err)
	}

	lastBackupedVersionFile := badgerDir + "/last_version.json"
	addLastBackupFile(t, lastBackupedVersionFile, 0)

	tmpDir := "/tmp/BKP/"

	cfg := configuration.Backup{
		ConfirmFile:          "BACKUPED",
		MetaInfoFile:         "META.json",
		TargetDirectory:      tmpDir + "/TARGET/" + prefix,
		TmpDirectory:         tmpDir + "/TMP",
		DirNameTemplate:      "pulse-%d",
		BackupWaitPeriod:     10,
		BackupFile:           "incr.bkp",
		Enabled:              true,
		PostProcessBackupCmd: []string{"bash", "-c", cwd + "/post_process_backup.sh"},
		LastBackupInfoFile:   lastBackupedVersionFile,
	}

	err = os.MkdirAll(cfg.TargetDirectory, 0777)
	require.NoError(t, err)
	err = os.MkdirAll(cfg.TmpDirectory, 0777)
	require.NoError(t, err)

	return configuration.Ledger{
		Backup: cfg,
		Storage: configuration.Storage{
			DataDirectory: badgerDir,
		},
	}
}

func clearData(t *testing.T, cfg configuration.Ledger) {
	err := os.RemoveAll(cfg.Backup.TargetDirectory)
	require.NoError(t, err)
	err = os.RemoveAll(cfg.Backup.TmpDirectory)
	require.NoError(t, err)
}

func TestBackuper(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	savedKeys := make(map[store.Key]insolar.PulseNumber, 0)

	var stopWriting uint32
	sgWriteStopped := sync.WaitGroup{}
	sgWriteStopped.Add(1)

	testPulse := insolar.GenesisPulse.PulseNumber + insolar.PulseNumber(rand.Int()%20000+1)
	// writing data to db
	go func() {
		for i := 0; i < 2000000; i++ {
			if atomic.LoadUint32(&stopWriting) != 0 {
				break
			}
			key := &testKey{id: uint64(i)}
			value := testPulse + insolar.PulseNumber(i)
			err := db.Set(key, value.Bytes())
			require.NoError(t, err)
			savedKeys[key] = value
			time.Sleep(time.Duration(rand.Int()%10) * time.Millisecond)
		}
		sgWriteStopped.Done()
	}()

	wgBackup := sync.WaitGroup{}
	numIterations := 15

	wgBackup.Add(numIterations)
	// doing backups
	go func() {
		for i := 0; i < numIterations; i++ {
			err := bm.MakeBackup(context.Background(), testPulse+insolar.PulseNumber(i))
			require.NoError(t, err)
			wgBackup.Done()
			time.Sleep(time.Duration(rand.Int()%1000) * time.Millisecond)
		}
	}()

	// wait for all backups done
	wgBackup.Wait()
	// stop writing to db
	atomic.StoreUint32(&stopWriting, 1)
	// wait for stopping
	sgWriteStopped.Wait()

	require.NotEqual(t, 0, len(savedKeys))

	// final backup to collect all rest records
	err = bm.MakeBackup(context.Background(), testPulse+insolar.PulseNumber(numIterations))
	require.NoError(t, err)

	// check backup hashes
	checkBackupMetaInfo(t, cfg.Backup, numIterations, testPulse)

	// load all backups and check all records
	{
		recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
		require.NoError(t, err)
		defer os.RemoveAll(recovTmpDir)

		createDirForBackup(t, recovTmpDir)

		for i := 0; i < numIterations+1; i++ {
			bkpFileName := filepath.Join(
				cfg.Backup.TargetDirectory,
				fmt.Sprintf(cfg.Backup.DirNameTemplate, testPulse+insolar.PulseNumber(i)),
				cfg.Backup.BackupFile,
			)

			loadIncrementalBackup(t, recovTmpDir, bkpFileName)
		}

		recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
		require.NoError(t, err)
		defer recoveredDB.Stop(context.Background())

		for k, v := range savedKeys {
			gotRawValue, err := recoveredDB.Get(k)
			require.NoError(t, err)
			gotPulseNumber := insolar.NewPulseNumber(gotRawValue)
			require.Equal(t, v, gotPulseNumber)
		}
	}

	// check last backuped version file
	{
		require.NotEqual(t, 0, loadLastBackupedVersion(t, cfg.Backup.LastBackupInfoFile))
	}
}

func loadLastBackupedVersion(t *testing.T, fileName string) uint64 {
	raw, err := ioutil.ReadFile(fileName)
	require.NoError(t, err)
	var backupInfo executor.LastBackupInfo
	err = json.Unmarshal(raw, &backupInfo)
	require.NoError(t, err)

	return backupInfo.LastBackupedVersion
}

var binaryPath string

func init() {
	var ok bool

	binaryPath, ok = os.LookupEnv("BIN_DIR")
	if !ok {
		wd, err := os.Getwd()
		binaryPath = filepath.Join(wd, "..", "..", "..", "..", "bin")

		if err != nil {
			panic(err.Error())
		}
	}
}

// prepareBackup uses backupmanager utility to prepare backup for usage
func prepareBackup(t *testing.T, dbDir string, lastBackupInfoFile string) {
	println("=====> Start preparing backup")
	cmd := exec.Command(binaryPath+"/backupmanager", "prepare_backup", "-d", dbDir, "-l", lastBackupInfoFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)
	println("<===== Finish preparing backup")
}

// createDirForBackup uses backupmanager utility to create empty badger
func createDirForBackup(t *testing.T, dbDir string) {
	println("=====> Start creating db for backup")
	cmd := exec.Command(binaryPath+"/backupmanager", "create", "-d", dbDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)
	println("<===== Finish creating db for backup")
}

// loadIncrementalBackup uses backupmanager utility to roll backups
func loadIncrementalBackup(t *testing.T, dbDir string, backupFile string) {
	println("=====> Start loading backup")
	cmd := exec.Command(binaryPath+"/backupmanager", "merge", "-t", dbDir, "-n", backupFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)
	err = cmd.Wait()
	require.NoError(t, err)
	println("<===== Finish loading backup")
}

func makeCurrentBkpDir(cfg configuration.Backup, pulse insolar.PulseNumber) string {
	return filepath.Join(cfg.TargetDirectory, fmt.Sprintf(cfg.DirNameTemplate, pulse))
}

func calculateFileHash(t *testing.T, fileName string) string {
	f, err := os.Open(fileName)
	require.NoError(t, err)
	defer f.Close()
	hasher := sha256.New()
	_, err = io.Copy(hasher, f)
	require.NoError(t, err)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func checkBackupMetaInfo(t *testing.T, cfg configuration.Backup, numIterations int, testPulse insolar.PulseNumber) {
	for i := 0; i < numIterations+1; i++ {
		currentPulse := testPulse + insolar.PulseNumber(i)
		currentBkpDir := makeCurrentBkpDir(cfg, currentPulse)
		metaInfo := filepath.Join(currentBkpDir, cfg.MetaInfoFile)
		raw, err := ioutil.ReadFile(metaInfo)
		require.NoError(t, err)

		bi := executor.BackupInfo{}
		err = json.Unmarshal(raw, &bi)
		require.NoError(t, err)

		// check file hash
		bkpFile := filepath.Join(currentBkpDir, cfg.BackupFile)
		md5sum := calculateFileHash(t, bkpFile)
		require.Equal(t, md5sum, bi.SHA256)

		// check pulse
		require.Equal(t, currentPulse, bi.Pulse)
	}
}

func TestBackup_UsageOfLastBackupedVersion(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)
	addLastBackupFile(t, cfg.Backup.LastBackupInfoFile, 40000)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	key := &testKey{id: uint64(3)}

	err = db.Set(key, []byte{})
	require.NoError(t, err)

	_, err = db.Get(key)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 100000)
	require.NoError(t, err)

	recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(recovTmpDir)
	createDirForBackup(t, recovTmpDir)

	bkpFileName := filepath.Join(
		cfg.Backup.TargetDirectory,
		fmt.Sprintf(cfg.Backup.DirNameTemplate, 100000),
		cfg.Backup.BackupFile,
	)
	loadIncrementalBackup(t, recovTmpDir, bkpFileName)
	recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	defer recoveredDB.Stop(context.Background())

	// we should not find 'key', since tried to backup from non existing high version
	_, err = recoveredDB.Get(key)
	require.EqualError(t, err, store.ErrNotFound.Error())
}

func TestBackupSendDeleteRecords(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	key := &testKey{id: uint64(3)}
	deletedKey := &testKey{id: uint64(4)}

	err = db.Set(key, []byte{})
	require.NoError(t, err)

	err = db.Set(deletedKey, []byte{})
	require.NoError(t, err)

	err = db.Delete(deletedKey)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 100000)
	require.NoError(t, err)
	err = bm.MakeBackup(context.Background(), 100001)
	require.NoError(t, err)

	recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(recovTmpDir)

	createDirForBackup(t, recovTmpDir)
	for i := 0; i < 2; i++ {
		bkpFileName := filepath.Join(
			cfg.Backup.TargetDirectory,
			fmt.Sprintf(cfg.Backup.DirNameTemplate, 100000),
			cfg.Backup.BackupFile,
		)

		loadIncrementalBackup(t, recovTmpDir, bkpFileName)
	}
	recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	_, err = recoveredDB.Get(key)
	require.NoError(t, err)
	recoveredDB.Stop(context.Background())

	err = db.Delete(key)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 100002)
	require.NoError(t, err)

	bkpFileName := filepath.Join(
		cfg.Backup.TargetDirectory,
		fmt.Sprintf(cfg.Backup.DirNameTemplate, 100002),
		cfg.Backup.BackupFile,
	)
	loadIncrementalBackup(t, recovTmpDir, bkpFileName)
	recoveredDB, err = store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	defer recoveredDB.Stop(context.Background())

	_, err = recoveredDB.Get(key)
	require.EqualError(t, err, store.ErrNotFound.Error())

	_, err = recoveredDB.Get(deletedKey)
	require.EqualError(t, err, store.ErrNotFound.Error())
}

func TestBackup_FullCycle(t *testing.T) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)

	db, err := store.NewBadgerDB(badger.DefaultOptions(tmpdir))
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber, db)
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	testJet := insolar.ZeroJetID

	pulsesDB := pulse.NewDB(db)
	err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)
	err = pulsesDB.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetsDB := jet.NewDBStore(db)
	err = jetsDB.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)

	jetKeeper := executor.NewJetKeeper(jetsDB, db, pulsesDB)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	err = jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	err = bm.MakeBackup(ctx, testPulse)
	require.NoError(t, err)

	recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	defer os.RemoveAll(recovTmpDir)
	createDirForBackup(t, recovTmpDir)
	bkpFileName := filepath.Join(
		cfg.Backup.TargetDirectory,
		fmt.Sprintf(cfg.Backup.DirNameTemplate, testPulse),
		cfg.Backup.BackupFile,
	)
	loadIncrementalBackup(t, recovTmpDir, bkpFileName)
	prepareBackup(t, recovTmpDir, filepath.Base(cfg.Backup.LastBackupInfoFile))
	recoveredDB, err := store.NewBadgerDB(badger.DefaultOptions(recovTmpDir))
	require.NoError(t, err)
	defer recoveredDB.Stop(context.Background())

	recoveredJetKeeper := executor.NewJetKeeper(jet.NewDBStore(recoveredDB), recoveredDB, pulse.NewDB(recoveredDB))

	// pulse must be finalized when prepare_backup complete without error
	require.Equal(t, testPulse, recoveredJetKeeper.TopSyncPulse())
}
