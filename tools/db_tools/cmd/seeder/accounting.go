package seeder

import (
	"db_tools/cmd/template_parser"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"strings"
	"sync"
)

func seedFinancialGroups() {
	cmn.Log("seeding groups", cmn.LogLevels.Info)
	groups := template_parser.ParseGroupsTemplate()
	var wg sync.WaitGroup
	groupCreater := func(groupName string, groupData template_parser.GroupNode, groupType string, parentId int) {
		defer wg.Done()
		createGroupsAndSubGroups(groupName, groupData, groupType, parentId)
	}

	for groupName, groupData := range groups {
		parentId := -1
		// the group type is the same as the root group name in this case
		groupType := strings.ToLower(groupName)
		wg.Add(1)
		go groupCreater(groupName, groupData, groupType, parentId)
	}
	wg.Wait() // waits for all groups to be created
}

func seedFinancialAccounts() {
	cmn.Log("seeding accounts", cmn.LogLevels.Info)
	accounts := template_parser.ParseAccountssTemplate()
	var wg sync.WaitGroup
	accountCreater := func(a mdls.FinancialAccount) {
		defer wg.Done()
		createAccount(a)
	}
	for groupName, accounts := range accounts {
		group := readFromMap[mdls.FinancialGroup](allGroupsByName, groupName)
		accountType := determineAccountTypeFromGroupType(group.GroupType)
		gId := group.Id
		for _, a := range accounts {
			acc := mdls.FinancialAccount{
				Name:        a.Name,
				Description: a.Description,
				AccountType: accountType,
				GroupId:     gId,
			}
			wg.Add(1)
			go accountCreater(acc)
		}
	}
	wg.Wait()
}

func seedFinancialJournalEntries() {
	cmn.Log("seeding journal entries", cmn.LogLevels.Info)
	var js = []FullJournalEntry{
		// pay rent with cash
		{J: mdls.FinancialJournalEntry{EntryDate: "2023-11-14T22:44.905Z", Description: "Hey I am some random description"},
			Transactions: []mdls.FinancialTransaction{
				{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 2000 * 1000},
				{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Rent").Id, TransactionType: "debit", Amount: 2000 * 1000},
			}},
		// get cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2023-11-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 2000 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 2000 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2024-11-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 200 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 200 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2022-11-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 7400 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 7400 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2024-09-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 1400 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 1400 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2006-10-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 600 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 600 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2019-11-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 400 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 400 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2020-11-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 2010 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 2010 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2024-04-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 100 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 100 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2024-10-14T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 21200 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 21200 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2002-01-11T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 100 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 100 * 1000},
		}},
		// get more cash from equity
		{J: mdls.FinancialJournalEntry{EntryDate: "2022-01-11T22:44.905Z"}, Transactions: []mdls.FinancialTransaction{
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Cash").Id, TransactionType: "credit", Amount: 91200 * 1000},
			{FinancialAccountId: readFromMap[mdls.FinancialAccount](allAccountsByName, "Owner Equity").Id, TransactionType: "debit", Amount: 91200 * 1000},
		}},
		// todo add more accounts later
		// I got revenue
		// I got have to pay wages
	}

	var wg sync.WaitGroup
	entryCreater := func(j FullJournalEntry) {
		defer wg.Done()
		createJournalEntry(j)
	}

	for _, j := range js {
		wg.Add(1)
		entryCreater(j)
	}
	wg.Wait()
}
