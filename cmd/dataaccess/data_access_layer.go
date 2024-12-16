package dataaccess

import (
	"database/sql"
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"sync"
)

// this layer will decide whether to access
// the database or the redis cache
// or whereever data is stored.
// This is the only class that has access to the db object and can preform reads and writes.
//you should just pass data directly to the package involved, don't write logic here, just pass data

// some logic requires a transaction
// for example if we're performing a financial transfer
// between two accounts. We don't want to write the
// financial logic to the database pkg, so we
// allow the dataaccess layer to expose
// database transaction ability.
// we use sync.Mutex in case two process try
// to create a transaction at the same time
type trasactionsMapType struct {
	sync.Mutex
	m map[string]*sql.Tx
}

func (m *trasactionsMapType) Store(key string, tx *sql.Tx) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = tx
}
func (m *trasactionsMapType) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, key)
}
func (m *trasactionsMapType) Load(key string) (*sql.Tx, bool) {
	m.Lock()
	defer m.Unlock()
	v, ok := m.m[key]
	return v, ok
}

var transactionMap trasactionsMapType

func InitDataAccess() {
	initDatabase()
	transactionMap.m = make(map[string]*sql.Tx)
}

func ObtainTransactionKey() (string, error) {
	cmn.Log("DataAccessTransaction Begin", cmn.LogLevels.Operation)
	tx, err := db.Begin()
	if err != nil {
		cmn.Log("DataAccessTransaction failed", cmn.LogLevels.Operation)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		tx.Rollback()
		return "", err
	}
	key := fmt.Sprintf("%d", &tx)
	transactionMap.Store(key, tx)
	return key, nil
}
func RollbackTransaction(transactionKey string) {
	cmn.Log("DataAccessTransaction Rollback", cmn.LogLevels.Operation)
	tx, ok := transactionMap.Load(transactionKey)
	if ok {
		err := tx.Rollback()
		transactionMap.Delete(transactionKey)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Operation)
		}
	} else {
		cmn.Log("Cant find transaction key!", cmn.LogLevels.Critical)
	}
}
func CommitTransaction(transactionKey string) error {
	cmn.Log("DataAccessTransaction Commit", cmn.LogLevels.Operation)
	tx, ok := transactionMap.Load(transactionKey)
	if ok {
		err := tx.Commit()
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			transactionMap.Delete(transactionKey)
			return err
		}
	} else {
		cmn.Log("Cant find transaction key!", cmn.LogLevels.Critical)
	}
	return nil
}
