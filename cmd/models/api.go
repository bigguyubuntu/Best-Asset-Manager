package mdls

import (
	cmn "inventory_money_tracking_software/cmd/common"
)

type ApiError struct {
	Code        cmn.ErrorCode
	Description string
}
