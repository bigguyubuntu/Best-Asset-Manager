package inventory

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
	"strconv"
)

func CreateItem(item mdls.Item) (string, cmn.ErrorCode) {
	if item.Name == "" {
		return "", cmn.ErrorCodes.NameWasEmpty
	}
	iId, err := dataaccess.CreateNewItem(item)
	if err != nil || iId == 0 {
		return "", cmn.ErrorCodes.CreateFailed
	}
	itemId := strconv.Itoa(iId)
	return itemId, cmn.ErrorCodes.NoError
}

func ReadAllItems() []mdls.Item {
	l := []mdls.Item{}
	item := mdls.Item{
		Name: "test item no actual db was called",
	}
	l = append(l, item)
	return l
}

func ReadItemById(itemId int) (mdls.Item, cmn.ErrorCode) {
	item := mdls.Item{}
	if itemId == 0 {
		return item, cmn.ErrorCodes.InvalidId
	}
	item, err := dataaccess.ReadItem("", itemId)
	if err != nil {
		return item, cmn.ErrorCodes.ReadFailed
	}
	return item, cmn.ErrorCodes.NoError
}

// can update eveything except the quantity, to update quantity you need
// to do an inventory journal change
func UpdateItem(item mdls.Item) cmn.ErrorCode {
	return cmn.ErrorCodes.NoError
}

// goes over the cache and puts the updated quanitites in the database
func postItemQuantityChangesToDb(txKey string, itemsCache map[int]mdls.Item) error {
	for iId, item := range itemsCache {
		oldItem, err := dataaccess.ReadItem(txKey, iId)
		if err != nil {
			return err
		}
		item.Count = oldItem.Count + item.Count
		err = dataaccess.UpdateItemCount(txKey, iId, item.Count)
		if err != nil {
			return err
		}
	}
	return nil
}
