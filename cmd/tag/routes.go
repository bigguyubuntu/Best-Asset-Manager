package tag

import "inventory_money_tracking_software/cmd/api"

var base = api.Prefix + "/tag"

type tagRoutesStruct struct {
	read                              string // reads a tag by id, route must be /tag/1 for example
	create                            string
	update                            string
	delete                            string
	readTagsByFinancialJournalEntryId string
	readTagsByInventoryJournalEntryId string
}

var tagRoutes = tagRoutesStruct{
	read:   base + "/",
	create: base + "/create",
}
