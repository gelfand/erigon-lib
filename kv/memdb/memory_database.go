/*
   Copyright 2021 Erigon contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package memdb

import (
	"context"
	"testing"

	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon-lib/kv/mdbx"
	"github.com/ledgerwatch/log/v3"
)

func New() kv.RwDB {
	logger := log.New() //TODO: move higher
	return mdbx.NewMDBX(logger).InMem().MustOpen()
}

func NewTestDB(t testing.TB) kv.RwDB {
	db := New()
	t.Cleanup(db.Close)
	return db
}
func NewTestPoolDB(t testing.TB) kv.RwDB {
	logger := log.New() //TODO: move higher
	db := mdbx.NewMDBX(logger).InMem().WithTablessCfg(func(defaultBuckets kv.TableCfg) kv.TableCfg { return kv.TxpoolTablesCfg }).MustOpen()
	t.Cleanup(db.Close)
	return db
}

func NewTestPoolTx(t testing.TB) (kv.RwDB, kv.RwTx) {
	db := NewTestPoolDB(t)
	tx, err := db.BeginRw(context.Background()) //nolint
	if err != nil {
		t.Fatal(err)
	}
	if t != nil {
		t.Cleanup(tx.Rollback)
	}
	return db, tx
}

func NewTestTx(t testing.TB) (kv.RwDB, kv.RwTx) {
	db := New()
	t.Cleanup(db.Close)
	tx, err := db.BeginRw(context.Background()) //nolint
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(tx.Rollback)
	return db, tx
}
