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

package executor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"
)

func TestBackuper_BadConfig(t *testing.T) {
	existingDir, err := os.Getwd()
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber

	cfg := configuration.Backup{TmpDirectory: "-----", Enabled: true}

	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "check TmpDirectory returns error: stat -----: no such file or directory")

	cfg = configuration.Backup{TmpDirectory: existingDir, TargetDirectory: "+_+_+_+", Enabled: true}
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "check TargetDirectory returns error: stat +_+_+_+: no such file or directory")

	cfg.TargetDirectory = existingDir
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "ConfirmFile can't be empty")

	cfg.ConfirmFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "MetaInfoFile can't be empty")

	cfg.MetaInfoFile = "Test2"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "DirNameTemplate can't be empty")

	cfg.DirNameTemplate = "Test3"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "BackupWaitPeriod can't be 0")

	cfg.BackupWaitPeriod = 20
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "BackupFile can't be empty")

	cfg.BackupFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "PostProcessBackupCmd can't be empty")

	cfg.PostProcessBackupCmd = []string{"some command"}
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "LastBackupInfoFile can't be empty")

	tmpDir := "/tmp/BKP/"
	err = os.MkdirAll(tmpDir, 0777)
	defer os.RemoveAll(tmpDir)
	require.NoError(t, err)
	lastBackupedVersionFile := tmpDir + "/last_version.json"
	addLastBackupFile(t, lastBackupedVersionFile, 200)

	db := store.NewDBMock(t)
	db.GetMock.Return([]byte{}, nil)

	cfg.LastBackupInfoFile = lastBackupedVersionFile
	storageConfig := configuration.Storage{DataDirectory: tmpDir}
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg, Storage: storageConfig}, testPulse, db)
	require.NoError(t, err)
}

func addLastBackupFile(t *testing.T, to string, lastBackupedVersion uint64) {
	backupInfo := executor.LastBackupInfo{
		LastBackupedVersion: lastBackupedVersion,
	}
	rawInfo, err := json.MarshalIndent(backupInfo, "", "    ")
	require.NoError(t, err)

	err = ioutil.WriteFile(to, rawInfo, 0600)
	require.NoError(t, err)
}

func makeBackuperConfig(t *testing.T, prefix string, badgerDir string) configuration.Ledger {

	lastBackupedVersionFile := badgerDir + "/last_version.json"
	addLastBackupFile(t, lastBackupedVersionFile, 200)

	tmpDir := "/tmp/BKP/"
	cfg := configuration.Backup{
		ConfirmFile:          "BACKUPED",
		MetaInfoFile:         "META.json",
		TargetDirectory:      tmpDir + "/TARGET/" + prefix,
		TmpDirectory:         tmpDir + "/TMP",
		DirNameTemplate:      "pulse-%d",
		BackupWaitPeriod:     60,
		BackupFile:           "incr.bkp",
		Enabled:              true,
		PostProcessBackupCmd: []string{"ls"},
		LastBackupInfoFile:   lastBackupedVersionFile,
	}

	err := os.MkdirAll(cfg.TargetDirectory, 0777)
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

func TestBackuper_Disabled(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name(), os.TempDir())
	cfg.Backup.Enabled = false
	defer clearData(t, cfg)
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, 0, nil)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 1)
	require.Equal(t, err, executor.ErrBackupDisabled)
}

func TestBackuper_PostProcessCmdReturnError(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)

	cfg.Backup.BackupWaitPeriod = 1

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	cfg.Backup.PostProcessBackupCmd = []string{""}
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.Contains(t, err.Error(), "failed to start post process command")
}

func TestBackuper_HappyPath(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)

	cfg.Backup.BackupWaitPeriod = 1

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	confirmFile := filepath.Join(cfg.Backup.TargetDirectory, fmt.Sprintf(cfg.Backup.DirNameTemplate, testPulse+1), cfg.Backup.ConfirmFile)
	cfg.Backup.PostProcessBackupCmd = []string{"touch", confirmFile}
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.NoError(t, err)
}

func TestBackuper_BackupWaitPeriodExpired(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)
	cfg.Backup.BackupWaitPeriod = 1

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.Contains(t, err.Error(), "no backup confirmation")
}

func TestBackuper_Backup_OldPulse(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name(), os.TempDir())
	defer clearData(t, cfg)

	db := store.NewDBMock(t)
	db.GetMock.Return([]byte{}, nil)

	testPulse := insolar.GenesisPulse.PulseNumber
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse)
	require.Equal(t, err, executor.ErrAlreadyDone)

	err = bm.MakeBackup(context.Background(), testPulse-1)
	require.Equal(t, err, executor.ErrAlreadyDone)
}

func TestBackuper_TruncateHead(t *testing.T) {

	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, cfg)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	numElements := 10

	for i := 0; i < numElements; i++ {
		err = db.Set(executor.BackupStartKey(testPulse+insolar.PulseNumber(i)), []byte{})
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2

	err = bm.TruncateHead(context.Background(), testPulse+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		_, err = db.Get(executor.BackupStartKey(testPulse + insolar.PulseNumber(i)))
		require.NoError(t, err)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err = db.Get(executor.BackupStartKey(testPulse + insolar.PulseNumber(i)))
		require.EqualError(t, err, store.ErrNotFound.Error())
	}
}
