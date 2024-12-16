package inventory

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
	"strconv"
)

func CreateStorehouse(s mdls.Storehouse) (string, cmn.ErrorCode) {
	if s.Name == "" {
		return "", cmn.ErrorCodes.NameWasEmpty
	} else if s.StorehouseType == "" {
		return "", cmn.ErrorCodes.InvalidStorehouseType
	}
	sId, err := dataaccess.CreateNewStorehouse("", s)
	if err != nil || sId == 0 {
		return "", cmn.ErrorCodes.CreateFailed
	}
	return strconv.Itoa(sId), cmn.ErrorCodes.NoError
}

func ReadStorehouse(sId int) (mdls.Storehouse, cmn.ErrorCode) {
	s := mdls.Storehouse{}
	if sId == 0 {
		return s, cmn.ErrorCodes.InvalidId
	}
	s, err := dataaccess.ReadStorehouse("", sId)
	if err != nil {
		return s, cmn.ErrorCodes.ReadFailed
	}
	return s, cmn.ErrorCodes.NoError
}

func UpdateStorehouse(store mdls.Storehouse) cmn.ErrorCode {
	return cmn.ErrorCodes.NoError
}

// will read a storehouse from cache, if it's not there, it will read from db and update cache
func readStorehouseWithCache(txKey string, sId int, storehouseCache map[int]mdls.Storehouse) (mdls.Storehouse, error) {
	store, ok := storehouseCache[sId]
	if ok {
		return store, nil
	}
	store, err := dataaccess.ReadStorehouse(txKey, sId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return mdls.Storehouse{}, err
	}
	storehouseCache[sId] = store
	return store, nil
}
