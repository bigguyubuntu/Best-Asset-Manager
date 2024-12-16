package mdls

import (
	cmn "inventory_money_tracking_software/cmd/common"
	"testing"
)

func TestEnsurePositiveQuantity(test *testing.T) {
	m := InventoryMovement{
		Quantity: 10,
		MoveType: cmn.InventoryMovementTypes.Increase,
	}
	m.EnsurePositiveQuantity()
	if m.Quantity != 10 || m.MoveType != cmn.InventoryMovementTypes.Increase {
		test.Errorf("expected quantity of 10 and move type of increase")
	}
	m.Quantity = -10
	m.EnsurePositiveQuantity()
	if m.Quantity != 10 || m.MoveType != cmn.InventoryMovementTypes.Decrease {
		test.Errorf("expected quantity of 10 and move type of decrease")
	}
	m = InventoryMovement{
		Quantity: -10,
		MoveType: cmn.InventoryMovementTypes.Decrease,
	}
	m.EnsurePositiveQuantity()
	if m.Quantity != 10 || m.MoveType != cmn.InventoryMovementTypes.Increase {
		test.Errorf("expected quantity of 10 and move type of increase")
	}
}

func TestUpdateQuantityAfterInventoryMove(test *testing.T) {
	i := Item{Id: 1}
	s := Storehouse{Id: 5,
		StorehouseType: cmn.StorehouseTypes.Owned}
	m := InventoryMovement{
		Id:           1,
		MoveType:     cmn.InventoryMovementTypes.Increase,
		Quantity:     10,
		ItemId:       i.Id,
		StorehouseId: s.Id,
	}
	// if movement is increaseing we should increase the quantity
	q := StorehouseItemQuantity{
		ItemId:       i.Id,
		StorehouseId: s.Id,
	}
	q.UpdateQuantityAfterInventoryMove(m, s)
	if q.Quantity != m.Quantity {
		test.Errorf("expected that StorehouseItemQuantity after a move will be %f but instead it was %f",
			m.Quantity, q.Quantity)
	}

	// experiment with negative move type
	q = StorehouseItemQuantity{
		ItemId:       i.Id,
		StorehouseId: s.Id,
		Quantity:     0,
	}
	m.MoveType = cmn.InventoryMovementTypes.Decrease
	q.UpdateQuantityAfterInventoryMove(m, s)
	if q.Quantity != -1*m.Quantity {
		test.Errorf("expected that StorehouseItemQuantity after a move will be %f but instead it was %f",
			-1*m.Quantity, q.Quantity)
	}

	// experiment with partial move
	q = StorehouseItemQuantity{
		ItemId:       i.Id,
		StorehouseId: s.Id,
		Quantity:     0,
	}
	m.MoveType = cmn.InventoryMovementTypes.Decrease
	m.Quantity = 0.5
	q.UpdateQuantityAfterInventoryMove(m, s)
	if q.Quantity != -1*m.Quantity {
		test.Errorf("expected that StorehouseItemQuantity after a move will be %f but instead it was %f",
			-1*m.Quantity, q.Quantity)
	}

	// experiment with multiple moves both positive and negative
	q = StorehouseItemQuantity{
		ItemId:       i.Id,
		StorehouseId: s.Id,
		Quantity:     0,
	}
	m.MoveType = cmn.InventoryMovementTypes.Decrease
	m.Quantity = 10
	q.UpdateQuantityAfterInventoryMove(m, s)
	m.MoveType = cmn.InventoryMovementTypes.Increase
	m.Quantity = 15.5
	q.UpdateQuantityAfterInventoryMove(m, s)
	expected := 15.5 - 10
	if q.Quantity != expected {
		test.Errorf("expected that StorehouseItemQuantity after a move will be %f but instead it was %f",
			expected, q.Quantity)
	}
}

func TestUpdateItemCountAfterMove(t *testing.T) {
	i := Item{Id: 1}
	s := Storehouse{Id: 5,
		StorehouseType: cmn.StorehouseTypes.Owned}
	m := InventoryMovement{
		Id:           1,
		MoveType:     cmn.InventoryMovementTypes.Increase,
		Quantity:     10.0,
		ItemId:       i.Id,
		StorehouseId: s.Id,
	}
	i.UpdateCountAfterMovement(m, s)
	if i.Count != m.Quantity {
		t.Errorf("expected item count to be %f but it was %f",
			m.Quantity, i.Count)
	}
	// now it shouldn't udpate if we have a no owner storehouse
	i2 := Item{Id: 2}
	s2 := Storehouse{Id: 99,
		StorehouseType: cmn.StorehouseTypes.NotOwned}
	m2 := InventoryMovement{
		Id:           2,
		MoveType:     cmn.InventoryMovementTypes.Increase,
		Quantity:     10.0,
		ItemId:       i.Id,
		StorehouseId: s2.Id,
	}
	i2.UpdateCountAfterMovement(m2, s2)
	if i2.Count != 0 {
		t.Errorf("Expected item count to not change and stay zero if we changed a not-owned storehouse, instead we found count of %f",
			i2.Count)
	}
}
