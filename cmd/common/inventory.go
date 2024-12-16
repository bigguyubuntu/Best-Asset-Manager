package cmn

type MoveType string
type moveTypes struct {
	Increase MoveType
	Decrease MoveType
}

var InventoryMovementTypes moveTypes = moveTypes{
	Increase: "increase",
	Decrease: "decrease",
}

type StorehouseType string
type storehouseTypes struct {
	Owned    StorehouseType
	NotOwned StorehouseType
}

var StorehouseTypes storehouseTypes = storehouseTypes{
	Owned:    "owned",
	NotOwned: "not_owned",
}

var ImbalanceStorehouseName string = "ImbalanceStorehouse"
