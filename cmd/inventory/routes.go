package inventory

import "inventory_money_tracking_software/cmd/api"

var base = api.Prefix + "/inventory"

type inventoryReadRoutes struct {
	GetAllItems string
}
type inventoryCreateRoutes struct {
	CreateItem string
}

var InventoryReadRoutes = inventoryReadRoutes{
	GetAllItems: base + "/get_all_items",
}

var InventoryCreateRoutes = inventoryCreateRoutes{
	CreateItem: base + "/create_item",
}
