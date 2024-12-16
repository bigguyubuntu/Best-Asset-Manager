package mdls

type Tag struct {
	dateTracking
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TagFinancialJournalEntry struct {
	dateTracking
	Id                      int `json:"id"`
	TagId                   int `json:"tag_id"`
	FinancialJournalEntryId int `json:"financial_journal_entry_id"`
}

func (t Tag) AreTagsEqual(t1, t2 Tag) bool {
	if t1.Id == t2.Id &&
		t1.Name == t2.Name &&
		t1.Description == t2.Description &&
		t1.CreatedAt == t2.CreatedAt &&
		t1.UpdatedAt == t2.UpdatedAt {
		return true
	}
	return false
}
