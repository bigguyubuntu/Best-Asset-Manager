package inventory

import (
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
)

func CreateStorehouseGroup(store mdls.StorehouseGroup) (string, cmn.ErrorCode) {
	return "", cmn.ErrorCodes.NoError
}
