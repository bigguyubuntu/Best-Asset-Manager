package inventory

import (
	"fmt"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
)

// if transaction key is empty string then the read will happen as a normal
// database read. If we pass a transaction key then the read will be a part of
// the transaction
func ReadStorehouseItemQuantity(txKey string, itemId,
	storehouseId int) (mdls.StorehouseItemQuantity, cmn.ErrorCode) {
	q := mdls.StorehouseItemQuantity{ItemId: itemId,
		StorehouseId: storehouseId}
	if itemId == 0 || storehouseId == 0 {
		return q, cmn.ErrorCodes.ReadFailed
	}
	q, err := dataaccess.ReadStorehouseItemQuantity(txKey, itemId, storehouseId)
	if err != nil {
		return q, cmn.ErrorCodes.ReadFailed
	}
	return q, cmn.ErrorCodes.NoError
}

// will either create or update a quanity in the cache.
// Note that in case of update the function's param
// is the "delta", in other words by how much to change the already existing quantity.
// After we create or update a quantity we store the result in database
func createOrUpdateStorehouseItemQuantity(txKey string, m mdls.InventoryMovement,
	store mdls.Storehouse) error {
	iId := m.ItemId
	sId := m.StorehouseId
	q, errCode := ReadStorehouseItemQuantity(txKey, iId, sId)
	createNew := errCode != cmn.ErrorCodes.NoError
	if createNew { // it's not here, we need to create it
		q = mdls.StorehouseItemQuantity{
			ItemId:       iId,
			StorehouseId: sId,
		}
	}
	ok := q.UpdateQuantityAfterInventoryMove(m, store)
	if !ok {
		err := fmt.Errorf("failed at creating StorehosueItemQuantity with item id %d and storehouse id %d",
			iId, sId)
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return err
	}
	if createNew { // create the storehosueItemQuantity in the db
		err := dataaccess.CreateStorehouseItemQuantity(txKey, q)
		return err
	} else { // the storehouseItemQuantity already exists so just update the db record
		err := dataaccess.UpdateStorehouseItemQuantity(txKey, q)
		return err
	}
}
