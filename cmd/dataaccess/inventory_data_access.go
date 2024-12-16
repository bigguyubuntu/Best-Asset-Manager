package dataaccess

import (
	"errors"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
)

func CreateNewItem(item mdls.Item) (int, error) {
	return insertItem(item)
}

func ReadItem(txKey string, itemId int) (mdls.Item, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return mdls.Item{}, errors.New("")
		}
		return selectItem(tx, itemId)
	}
	return selectItem(db, itemId)

}

// reads storehouse by ID, if txKey is empty string will use a normal db
// connection. otherwise will use an sql transaction
func ReadStorehouse(txKey string, id int) (mdls.Storehouse, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return mdls.Storehouse{}, errors.New("")
		}
		return selectStorehouse(tx, id)
	}
	return selectStorehouse(db, id)
}
func ReadStorehouseByName(txKey string, name string) (mdls.Storehouse, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return mdls.Storehouse{}, errors.New("")
		}
		return selectStorehouseByName(tx, name)
	}
	return selectStorehouseByName(db, name)
}

func CreateNewStorehouse(txKey string, s mdls.Storehouse) (int, error) {
	if txKey != "" {
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return 0, errors.New("")
		}
		return insertStorehouse(tx, s)
	}
	return insertStorehouse(db, s)
}

func CreateInventoryMovement(txKey string, move mdls.InventoryMovement) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return errors.New("")
	}
	return insertInventoryMovement(tx, move)
}
func CreateInventoryJournalEntry(txKey string,
	ij mdls.InventoryJournalEntry) (int, error) {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return 0, errors.New("")
	}
	return insertInventoryJournalEntry(tx, ij)
}

// if passed txKey is nil then we will use normal db connection. Otherwise will use a transaction
func ReadStorehouseItemQuantity(txKey string, itemId, storehouseId int) (mdls.StorehouseItemQuantity, error) {
	if txKey != "" {
		cmn.Log("ReadStorehouseItemQuantity within a db transaction for storehouse id %d and item id %d",
			cmn.LogLevels.Operation,
			storehouseId, itemId)
		tx, ok := transactionMap.Load(txKey)
		if !ok {
			cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
			return mdls.StorehouseItemQuantity{}, errors.New("")
		}
		return selectStorehouseItemQuantity(tx, itemId, storehouseId)
	}
	return selectStorehouseItemQuantity(db, itemId, storehouseId)
}

// make sure to not allow users to directly update quantity
// a quanitty is only updated as a part of an inventory journal entry.
func CreateStorehouseItemQuantity(txKey string, q mdls.StorehouseItemQuantity) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return errors.New("")
	}
	_, err := insertStorehouseItemQuantity(tx, q)
	return err
}

func UpdateStorehouseItemQuantity(txKey string, q mdls.StorehouseItemQuantity) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return errors.New("")
	}
	err := updateStorehouseItemQuantity(tx, q)
	return err
}
func ReadInventoryMovement(id int) (mdls.InventoryMovement, error) {
	return selectInventoryMovement(id)
}
func ReadInventoryMovementByInventoryJournalEntryId(id int) ([]mdls.InventoryMovement, error) {
	return selectInventoryMovementByInventoryJournalEntryId(id)
}
func UpdateItemCount(txKey string, itemId int, newCount float64) error {
	tx, ok := transactionMap.Load(txKey)
	if !ok {
		cmn.Log("couldnt get transaction using transaction key", cmn.LogLevels.Error)
		return errors.New("")
	}
	return updateItemCount(tx, itemId, newCount)
}
func UpdateItem(txKey string, item mdls.Item) error {
	return updateItem(item)
}
