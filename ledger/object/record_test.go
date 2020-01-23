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

// +build slowtest

package object

import (
	"context"
	"os"
	"sync"
	"testing"

	fuzz "github.com/google/gofuzz"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/ledger/heavy/migration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/tests/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var (
	poolLock     sync.Mutex
	globalPgPool *pgxpool.Pool
)

func setPool(pool *pgxpool.Pool) {
	poolLock.Lock()
	defer poolLock.Unlock()
	globalPgPool = pool
}

func getPool() *pgxpool.Pool {
	poolLock.Lock()
	defer poolLock.Unlock()
	return globalPgPool
}

// TestMain does the before and after setup
func TestMain(m *testing.M) {
	ctx := context.Background()
	log.Info("[TestMain] About to start PostgreSQL...")
	pgURL, stopPostgreSQL := common.StartPostgreSQL()
	log.Info("[TestMain] PostgreSQL started!")

	pool, err := pgxpool.Connect(ctx, pgURL)
	if err != nil {
		stopPostgreSQL()
		log.Panicf("[TestMain] pgxpool.Connect() failed: %v", err)
	}

	migrationPath := "../../migration"
	cwd, err := os.Getwd()
	if err != nil {
		stopPostgreSQL()
		panic(errors.Wrap(err, "[TestMain] os.Getwd failed"))
	}
	log.Infof("[TestMain] About to run PostgreSQL migration, cwd = %s, migration migrationPath = %s", cwd, migrationPath)
	ver, err := migration.MigrateDatabase(ctx, pool, migrationPath)
	if err != nil {
		stopPostgreSQL()
		panic(errors.Wrap(err, "Unable to migrate database"))
	}
	log.Infof("[TestMain] PostgreSQL database migration done, current schema version: %d", ver)

	setPool(pool)
	// Run all tests
	code := m.Run()

	log.Info("[TestMain] Cleaning up...")
	stopPostgreSQL()
	os.Exit(code)
}

func TestSetNilSignature(t *testing.T) {
	ctx := context.Background()
	db := NewRecordDB(getPool())
	f := fuzz.New()

	records := make([]record.Material, 1)
	records[0] = record.Material{
		Virtual:   record.Virtual{},
		ID:        gen.ID(),
		ObjectID:  gen.ID(),
		JetID:     gen.JetID(),
		Signature: nil,
	}
	f.Fuzz(&records[0].Polymorph)

	err := db.Set(ctx, records[0])
	require.Error(t, err)

	err = db.BatchSet(ctx, records)
	require.Error(t, err)

	// Make sure no records where ceated
	_, err = db.ForID(ctx, records[0].ID)
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	db := NewRecordDB(getPool())
	f := fuzz.New()
	// Make sure there is no record with such ID
	id := gen.ID()
	_, err := db.ForID(ctx, id)
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)

	rec1 := record.Material{
		Virtual:  record.Virtual{},
		ID:       id,
		ObjectID: gen.ID(),
		JetID:    gen.JetID(),
	}

	f.Fuzz(&rec1.Polymorph)
	f.NilChance(0).Fuzz(&rec1.Signature)

	err = db.Set(ctx, rec1)
	require.NoError(t, err)

	rec2, err := db.ForID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, rec1, rec2)
}

func TestBatchSet(t *testing.T) {
	ctx := context.Background()
	db := NewRecordDB(getPool())
	f := fuzz.New()

	var ids [3]insolar.ID
	for i := 0; i < len(ids); i++ {
		ids[i] = gen.ID()
		// Make sure there is no record with such ID
		_, err := db.ForID(ctx, ids[i])
		require.Error(t, err)
		require.Equal(t, ErrNotFound, err)
	}

	records := make([]record.Material, len(ids))
	for i := 0; i < len(records); i++ {
		records[i] = record.Material{
			Virtual:  record.Virtual{},
			ID:       ids[i],
			ObjectID: gen.ID(),
			JetID:    gen.JetID(),
		}
		f.Fuzz(&records[i].Polymorph)
		f.NilChance(0).Fuzz(&records[i].Signature)
	}

	err := db.BatchSet(ctx, records)
	require.NoError(t, err)

	for i := 0; i < len(ids); i++ {
		rec, err := db.ForID(ctx, ids[i])
		require.NoError(t, err)
		require.Equal(t, records[i], rec)
	}
}

func TestPosition(t *testing.T) {
	ctx := context.Background()
	db := NewRecordDB(getPool())
	f := fuzz.New()

	// Make sure there are no records with such pulse
	pn := gen.PulseNumber()
	_, err := db.LastKnownPosition(pn)
	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)

	for ctr := 1; ctr <= 3; ctr++ {
		// Make sure there is no record with such ID
		id := gen.IDWithPulse(pn)
		_, err = db.ForID(ctx, id)
		require.Error(t, err)
		require.Equal(t, ErrNotFound, err)

		_, err = db.AtPosition(pn, uint32(ctr))
		require.Error(t, err)
		require.Equal(t, ErrNotFound, err)

		rec1 := record.Material{
			Virtual:  record.Virtual{},
			ID:       id,
			ObjectID: gen.ID(),
			JetID:    gen.JetID(),
		}
		f.Fuzz(&rec1.Polymorph)
		f.NilChance(0).Fuzz(&rec1.Signature)

		err = db.Set(ctx, rec1)
		require.NoError(t, err)

		pos, err := db.LastKnownPosition(pn)
		require.NoError(t, err)
		require.Equal(t, uint32(ctr), pos)

		id2, err := db.AtPosition(id.Pulse(), uint32(ctr))
		require.NoError(t, err)
		require.Equal(t, id, id2)
	}
}
