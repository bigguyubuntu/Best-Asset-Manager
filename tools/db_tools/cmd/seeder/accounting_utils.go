package seeder

import (
	"bytes"
	"db_tools/cmd/template_parser"
	"encoding/json"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
	"io"
	"net/http"
	"os"
	"strconv"
)

type FullJournalEntry struct {
	J            mdls.FinancialJournalEntry  `json:"journal_entry"`
	Transactions []mdls.FinancialTransaction `json:"transactions"`
}

func handleBadResponse(resp *http.Response) {
	if resp.StatusCode != 200 {
		cmn.Log("error seeding", cmn.ErrorLevels.Error)
		cmn.Log("server response: ", cmn.ErrorLevels.Error)
		cmn.Log(readResponse(resp), cmn.ErrorLevels.Error)
		os.Exit(1)
	}
}
func readResponse(resp *http.Response) string {
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func createGroup(g mdls.FinancialGroup) mdls.FinancialGroup {
	url := baseUrl + "/accounting/group/create"
	body, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	bodyReader := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	handleBadResponse(resp)
	s := readResponse(resp)
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	g.Id = id
	writeToMap(allGroupsByName, g, g.Name)
	return g
}

func createAccount(a mdls.FinancialAccount) mdls.FinancialAccount {
	url := baseUrl + "/accounting/account/create"
	body, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	bodyReader := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	handleBadResponse(resp)
	s := readResponse(resp)
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	a.Id = id
	writeToMap[mdls.FinancialAccount](allAccountsByName, a, a.Name)
	return a
}

func createJournalEntry(j FullJournalEntry) {
	url := baseUrl + "/accounting/journal_entry/create"
	body, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	bodyReader := bytes.NewBuffer(body)
	resp, err := http.Post(url, "application/json", bodyReader)
	if err != nil {
		panic(err)
	}
	handleBadResponse(resp)
	resp.Body.Close()
}

// recursivly creates the group and all its line of children
func createGroupsAndSubGroups(groupName string, groupData template_parser.GroupNode, groupType string, parentId int) {
	g := mdls.FinancialGroup{
		Name:        groupName,
		Description: groupData.Description,
		ParentId:    parentId,
		GroupType:   groupType,
	}
	g = createGroup(g)
	parentId = g.Id
	for _, groupNode := range groupData.Children { // looks at the immediate children of this group
		for childGroupName, childGroupData := range groupNode {
			g = mdls.FinancialGroup{
				Name:        childGroupName,
				Description: childGroupData.Description,
				ParentId:    parentId,
			}
			createGroupsAndSubGroups(childGroupName, childGroupData, groupType, parentId)
		}
	}
}

func determineAccountTypeFromGroupType(groupType string) string {
	if groupType == "assets" ||
		groupType == "expenses" {
		return "debit_increased"
	}
	return "credit_increased"
}
