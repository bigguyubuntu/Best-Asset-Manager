package mdls

import (
	"errors"
	cmn "inventory_money_tracking_software/cmd/common"
)

type InventoryMovement struct {
	dateTracking
	Id                      int          `json:"id"`
	ItemId                  int          `json:"item_id"`
	StorehouseId            int          `json:"storehouse_id"`
	Quantity                float64      `json:"quantity"`
	MoveType                cmn.MoveType `json:"move_type"`
	InventoryJournalEntryId int          `json:"inventory_journal_entry_id"`
}

func (m InventoryMovement) ValidateMoveType() bool {
	if m.MoveType == cmn.InventoryMovementTypes.Decrease ||
		m.MoveType == cmn.InventoryMovementTypes.Increase {
		return true
	}
	return false
}

// if the quantity is negative it will make it positive and flip the move type
func (m *InventoryMovement) EnsurePositiveQuantity() {
	if m.Quantity < 0 {
		m.Quantity = -1 * m.Quantity
		if m.MoveType == cmn.InventoryMovementTypes.Increase {
			m.MoveType = cmn.InventoryMovementTypes.Decrease
			return
		} else if m.MoveType == cmn.InventoryMovementTypes.Decrease {
			m.MoveType = cmn.InventoryMovementTypes.Increase
			return
		}
	}
}

type StorehouseItemQuantity struct {
	dateTracking
	Id           int     `json:"id"`
	ItemId       int     `json:"item_id"`
	StorehouseId int     `json:"storehouse_id"`
	Quantity     float64 `json:"quantity"`
}

// updates the quantity of the StorehouseItemQuantity depending on the movement and storehouse
// types. Make sure to call the item.UpdateCount after calling this function.
// if the storehouse type is owned and the movement type is increase
// this function updates the quantity of the calling object. Will return true if
// update is successful otherwise will return false
func (q *StorehouseItemQuantity) UpdateQuantityAfterInventoryMove(m InventoryMovement,
	store Storehouse) bool {
	if store.Id == 0 || q.StorehouseId == 0 || q.ItemId == 0 {
		cmn.HandleError(errors.New("id zero is invalid"), cmn.ErrorLevels.Error)
		return false
	}
	if q.StorehouseId != m.StorehouseId || m.StorehouseId != store.Id {
		cmn.HandleError(errors.New("UpdateQuantityAfterInventoryMove tried to update quantity of the storehouse"),
			cmn.ErrorLevels.Error)
		return false
	}
	if !cmn.IsInventoryMovementTypeValid(m.MoveType) {
		cmn.HandleError(errors.New("UpdateQuantityAfterInventoryMove invalid inventoryMovement type"),
			cmn.ErrorLevels.Error)
		return false
	}
	m.EnsurePositiveQuantity()
	if m.MoveType == cmn.InventoryMovementTypes.Increase {
		q.Quantity += m.Quantity
		return true
	} else if m.MoveType == cmn.InventoryMovementTypes.Decrease {
		q.Quantity -= m.Quantity
		return true
	}
	return false
}

type InventoryJournalEntry struct {
	dateTracking
	Id          int    `json:"id"`
	Description string `json:"description"`
	EntryDate   string `json:"entryDate"`
}

type Storehouse struct {
	dateTracking
	Id             int                `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	StorehouseType cmn.StorehouseType `json:"storehouse_type"`
}

type Category struct {
	dateTracking
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type StorehouseGroup struct {
	dateTracking
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Item struct {
	dateTracking
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Imagelink     string  `json:"image_link"`
	InventoryUnit string  `json:"inventory_unit"`
	Count         float64 `json:"count"`
}

// item count is always relative to our ownership. We only care about
// storehouses that we own. We don't care about storehouses
// that we are not_owned
// if the storehouse is owned then if we add to it we increase item count
// if we decrease from that storehosue we reduce our item count.
// if a store is not owned, then changes to it will not impact
// the items count
func (i *Item) UpdateCountAfterMovement(m InventoryMovement, store Storehouse) {
	if store.StorehouseType != cmn.StorehouseTypes.Owned {
		return // exit early
	}
	if m.MoveType == cmn.InventoryMovementTypes.Increase {
		i.Count += m.Quantity
	} else if m.MoveType == cmn.InventoryMovementTypes.Decrease {
		i.Count -= m.Quantity
	}
}

var ImbalanceStorehouse Storehouse = Storehouse{
	Name:           cmn.ImbalanceStorehouseName,
	Description:    "An imbalance storehouse, it's auto generated and used to keep track of inventory movements that don't have a storehosue specified. Since it's of type not_owned any addition or subtraction from this storehouse will not affect the items' count",
	StorehouseType: cmn.StorehouseTypes.NotOwned,
}
