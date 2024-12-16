package seeder

import (
	mdls "inventory_money_tracking_software/cmd/models"
	"sync"
)

var allGroupsByName = make(map[string]mdls.FinancialGroup)
var allAccountsByName = make(map[string]mdls.FinancialAccount)

type dataModel interface {
	mdls.FinancialAccount | mdls.FinancialGroup
}

var tags = [5]mdls.Tag{
	{Name: "Groceries", Description: "Bismiallah al rahman al raheem"},
	{Name: "Client 131", Description: "Rabi ehdini"},
	{Name: "T-shirts", Description: "Rabi qawini"},
	{Name: "Toys", Description: "Rabi erzukni"},
	{Name: "Canada", Description: "Rabi 2rda 3ani"},
}

var lock = sync.RWMutex{}

func readFromMap[V dataModel](mapObj map[string]V, key string) V {
	lock.RLock()
	defer lock.RUnlock()
	return mapObj[key]
}
func writeToMap[V dataModel](mapObj map[string]V, obj V, key string) {
	lock.Lock()
	defer lock.Unlock()
	mapObj[key] = obj
}
