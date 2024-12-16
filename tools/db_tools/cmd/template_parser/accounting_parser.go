package template_parser

import (
	_ "embed"
	"encoding/json"
	cmn "inventory_money_tracking_software/cmd/common"
)

type GroupNode struct {
	Description string                 `json:"description"`
	Children    []map[string]GroupNode `json:"children"`
}

type AccountNode struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

//go:embed templates/accounting/groups.json
var groupsFileBytes []byte

//go:embed templates/accounting/accounts.json
var accountsFileBytes []byte

// returns the json template file for the groups.
// Each key is the group name, and the value is the
// group's data, it has the group description and children (subgroups)
// so data is in the following shape {"group_name": {description:"", children: [{"group_name": {description: "", children[]}}, {...}] }}
// There are 5 parent groups, these are always: Assets, Liabilities, Equity, Exepneses, Income. You can infer sub group
// types form the parent group.
func ParseGroupsTemplate() map[string]GroupNode {
	cmn.Log("Reading groups template json file", cmn.LogLevels.Operation)
	// the root groupsnodes will always be 5, for the 5 root groups
	// Assets, Liabilities, Equity, Exepneses, Income.
	var groups map[string]GroupNode
	json.Unmarshal(groupsFileBytes, &groups)
	return groups
}

// the key is the group this account belongs to, and the value is the list of accounts
// in that group. We can infer the account type from the group type.
func ParseAccountssTemplate() map[string][]AccountNode {
	cmn.Log("Reading Accounts template json file", cmn.LogLevels.Operation)
	var accounts map[string][]AccountNode
	json.Unmarshal(accountsFileBytes, &accounts)
	return accounts
}
