package inventory

import (
	"errors"
	cmn "inventory_money_tracking_software/cmd/common"
	"inventory_money_tracking_software/cmd/dataaccess"
	mdls "inventory_money_tracking_software/cmd/models"
)

func ReadInventoryMovement(moveId int) (mdls.InventoryMovement, cmn.ErrorCode) {
	return mdls.InventoryMovement{}, cmn.ErrorCodes.NoError
}

func UpdateInventoryMovement(m mdls.InventoryMovement) cmn.ErrorCode {
	return cmn.ErrorCodes.NoError
}

// makes sure that item movements are balanced, in other words if a storehouse increased by 4 items
// then there must be other storehouses that decrease by 4 items, if the movements aren't balanced
// then we will create an imbalance storehouse and move that item in or out to balance
func balanceMovements(txKey string, ms []mdls.InventoryMovement) ([]mdls.InventoryMovement, error) {
	balancedMoves := []mdls.InventoryMovement{}
	// key is itemId, value is a slice of all movements that belongs to that item
	specificMoves := make(map[int][]mdls.InventoryMovement)
	for _, m := range ms {
		if m.StorehouseId == 0 || m.ItemId == 0 {
			err := errors.New("storehouseId and itemId cant be zero")
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return balancedMoves, err
		}
		key := m.ItemId
		_, ok := specificMoves[key]
		if !ok {
			specificMoves[key] = []mdls.InventoryMovement{m}
		} else {
			specificMoves[key] = append(specificMoves[key], m)
		}
	}
	// now loop over every value in the map, and if it's not balanced add the
	// balancing moves
	for _, v := range specificMoves {
		l, err := balanceSingleItem(txKey, v)
		if err != nil {
			return balancedMoves, err
		}
		balancedMoves = append(balancedMoves, l...)
	}
	return balancedMoves, nil
}

// will take a slice of inventory movements who all have the same inventory id
// and storehouse id and will return the inventory movement slice after balancing
// the increase and decrease moves.
func balanceSingleItem(txKey string, ms []mdls.InventoryMovement) ([]mdls.InventoryMovement, error) {
	if len(ms) == 0 {
		cmn.Log("tried to balance an empty inventoryMovement array... investigate more",
			cmn.LogLevels.Error)
		return []mdls.InventoryMovement{}, nil
	}
	increase := 0.0
	decrease := 0.0
	iId := ms[0].ItemId
	for _, m := range ms {
		// make sure all itemids and storehouseIds are the same in the slice
		if m.ItemId != iId || m.ItemId == 0 || m.StorehouseId == 0 {
			err := errors.New(
				"Expected that the slice will have all equal item ids, and that no id is equal to zero")
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return ms, err
		}
		// sum all the increase and decrease in the
		if m.MoveType == cmn.InventoryMovementTypes.Increase {
			increase += m.Quantity
		} else if m.MoveType == cmn.InventoryMovementTypes.Decrease {
			decrease += m.Quantity
		}
	}
	if increase == decrease { // change nothing and return as the moves are balanced
		return ms, nil
	} else { // moves aren't balanced so we need to create imbalanced storehouse
		// and assign the difference...
		imbalanceStorehouse, err := findOrCreateImbalanceStorehouse(txKey)
		if err != nil {
			return ms, err
		}
		balancingMove := mdls.InventoryMovement{
			StorehouseId: imbalanceStorehouse.Id,
			ItemId:       iId,
		}

		if increase > decrease {
			balancingMove.Quantity = increase - decrease
			balancingMove.MoveType = cmn.InventoryMovementTypes.Decrease
		} else {
			balancingMove.Quantity = decrease - increase
			balancingMove.MoveType = cmn.InventoryMovementTypes.Increase
		}
		return append(ms, balancingMove), nil
	}
}

// will query the imbalanceStorehouse if it exists otherwise will create it
func findOrCreateImbalanceStorehouse(txKey string) (mdls.Storehouse, error) {
	s, err := dataaccess.ReadStorehouseByName(txKey, cmn.ImbalanceStorehouseName)
	if err == nil {
		return s, nil
	}
	cmn.Log("imbalance storehouse doesnt exist, will create it", cmn.LogLevels.Info)
	s = mdls.ImbalanceStorehouse
	sId, err := dataaccess.CreateNewStorehouse(txKey, s)
	s.Id = sId
	return s, err

}
