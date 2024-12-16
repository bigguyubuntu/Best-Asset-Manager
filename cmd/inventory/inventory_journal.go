package inventory

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
	"strconv"
)

// if movements per item are balanced then we:
//  1. create an inventory_journal_entry
//  2. update storehouse_item_quantity tables with new quantities,
//     negatives quantities are allowed in this table
//  3. update the inventory_movement table
//  4. update the items count, negatives are NOT allowed, fractions are okay
//     since we allow partial quantitites
func CreateInventoryJournalEntry(ij mdls.InventoryJournalEntry,
	ms []mdls.InventoryMovement) (string, cmn.ErrorCode) {
	storehouseCache := make(map[int]mdls.Storehouse)
	itemsCache := make(map[int]mdls.Item)

	txKey, err := dataaccess.ObtainTransactionKey()
	defer dataaccess.RollbackTransaction(txKey)
	if err != nil {
		cmn.Log("Failed to obtain transaction key while creating a new inventory journal entry",
			cmn.LogLevels.Error)
		return "", cmn.ErrorCodes.CreateFailed
	}
	ijId, err := dataaccess.CreateInventoryJournalEntry(txKey, ij)
	if ijId == 0 || err != nil {
		return "", cmn.ErrorCodes.CreateFailed
	}
	ij.Id = ijId

	ms, err = balanceMovements(txKey, ms)
	if err != nil {
		cmn.Log("failed to balance movements", cmn.LogLevels.Error)
		return "", cmn.ErrorCodes.CreateFailed
	}

	// todo maybe optimize by creating a local cache of fetched storehouse and items
	// and write their updates only once. Just like we do with updating and creating
	// financial journal entries
	for _, m := range ms {
		iId := m.ItemId
		sId := m.StorehouseId
		if iId == 0 || sId == 0 {
			cmn.Log("found zero id, either storehouse or item id is zero",
				cmn.LogLevels.Error)
			return "", cmn.ErrorCodes.CreateFailed
		}
		m.InventoryJournalEntryId = ijId
		err = dataaccess.CreateInventoryMovement(txKey, m)
		if err != nil {
			return "", cmn.ErrorCodes.CreateFailed
		}
		store, err := readStorehouseWithCache(txKey, sId, storehouseCache)
		if err != nil {
			return "", cmn.ErrorCodes.CreateFailed
		}
		err = createOrUpdateStorehouseItemQuantity(txKey, m, store)
		if err != nil {
			return "", cmn.ErrorCodes.CreateFailed
		}
		// update the items quantities in the cache
		i, ok := itemsCache[m.ItemId]
		if !ok {
			i = mdls.Item{Id: m.ItemId}
			i.UpdateCountAfterMovement(m, store)
			itemsCache[m.ItemId] = i
		} else {
			i.UpdateCountAfterMovement(m, store)
			itemsCache[m.ItemId] = i
		}
	}
	// create or update all the item caches
	err = postItemQuantityChangesToDb(txKey, itemsCache)
	if err != nil {
		return "", cmn.ErrorCodes.CreateFailed
	}
	err = dataaccess.CommitTransaction(txKey)
	if err != nil {
		return "", cmn.ErrorCodes.CreateFailed
	}
	return strconv.Itoa(ij.Id), cmn.ErrorCodes.NoError
}
