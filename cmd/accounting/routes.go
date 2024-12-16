package accounting

import "inventory_money_tracking_software/cmd/api"

var base string = api.Prefix + "/accounting"

type accountRoutesStruct struct {
	findByIds            string
	findByGroupId        string
	findByNameAndGroupId string
	findAll              string
	create               string
	edit                 string
	delete               string
}

var accountRoutes = accountRoutesStruct{
	findByIds:            base + "/account/find_by_ids",
	findByGroupId:        base + "/account/find_by_group",
	findByNameAndGroupId: base + "/account/find_by_name_and_group",
	findAll:              base + "/account/find_all",
	create:               base + "/account/create",
}

type groupRoutesStruct struct {
	findByIds             string
	findByType            string
	findByParentId        string
	findByNameAndParentId string
	findAll               string
	create                string
	delete                string
	update                string
}

var groupRoutes = groupRoutesStruct{
	findByIds:             base + "/group/find_by_ids",
	findByType:            base + "/group/find_by_type",
	findByParentId:        base + "/group/find_subgroups",
	findByNameAndParentId: base + "/group/find_by_name_and_parent",
	findAll:               base + "/group/find_all",
	create:                base + "/group/create",
	delete:                base + "/group/delete",
	update:                base + "/group/update",
}

type journalEntryRoutesStruct struct {
	findByIds                      string
	findAll                        string
	compositeFindAll               string
	compositeFindById              string
	compositeFindByAccountId       string
	create                         string
	updateTransactions             string // updates the transactions of the journal entry
	update                         string // updates the journal entry itself
	updateBothEntryAndTransactions string // updates the transactions and the entry
}

var journalEntryRoutes = journalEntryRoutesStruct{
	findByIds:                      base + "/journal_entry/find_by_ids",
	findAll:                        base + "/journal_entry/find_all",
	compositeFindAll:               base + "/journal_entry/composite_find_all",
	compositeFindById:              base + "/journal_entry/composite_find_by_id",
	compositeFindByAccountId:       base + "/journal_entry/composite_find_by_financial_account_id",
	create:                         base + "/journal_entry/create",
	update:                         base + "/journal_entry/update",
	updateTransactions:             base + "/journal_entry/update_transactions",
	updateBothEntryAndTransactions: base + "/journal_entry/update_both_transactions_and_entry",
}

type transactionRoutesStruct struct {
	findByAccountId      string
	findByJournalEntryId string
	findByJournalEntries string
}

var transactionRoutes = transactionRoutesStruct{
	findByAccountId:      base + "/transaction/find_by_account",
	findByJournalEntryId: base + "/transaction/find_by_journal_entry",
	findByJournalEntries: base + "/transaction/find_by_journal_entries",
}
