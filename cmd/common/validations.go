package cmn

import (
	"strings"
)

// this file will have validations that are common to many packages

func IsStringEmpty(s string) bool {
	// checks if the string is empty or just has white space
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return true
	}
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	if s == "" {
		return true
	}
	return false
}

func IsDuplicateKeyConstraintViolated(e string) bool {
	// some database rows don't allow duplicates, so if a duplicate is detected
	// the error.Error() will contain a specific string. We are using that string to detect
	// if the uniquness constraint is violated
	if strings.Contains(e, "duplicate key value violates unique constraint") ||
		strings.Contains(e, "already exist") {
		return true
	}
	return false
}

func IsGroupTypeValid(t string) bool {
	// this function checks if the accounting entitiy's type matches
	// one of the 3 allowed types, which are assets, liabilities and equity
	if t != GroupTypes.Assets &&
		t != GroupTypes.Liabilities &&
		t != GroupTypes.Expenses &&
		t != GroupTypes.Income &&
		t != GroupTypes.Equity {
		return false
	}
	return true
}

func IsAccountTypeValid(t string) bool {
	if t != AccountTypes.CreditIncreased && t != AccountTypes.DebitIncreased {
		return false
	}
	return true
}

func IsMoneyAmountValid(a int) bool {
	if a > (MaxBigInt / 1000) { // we store money in mils, and the max amount we can store is bigint divided by 1000
		return false
	} else if a < (-MaxBigInt / 1000) { // negative case
		return false
	}
	return true
}

func IsInventoryMovementTypeValid(t MoveType) bool {
	if t == InventoryMovementTypes.Increase ||
		t == InventoryMovementTypes.Decrease {
		return true
	}
	return false
}

func IsStorehouseTypeValid(t StorehouseType) bool {
	if t == StorehouseTypes.Owned ||
		t == StorehouseTypes.NotOwned {
		return true
	}
	return false
}
