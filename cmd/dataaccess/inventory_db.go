package dataaccess

import (
	"database/sql"
	"errors"
	cmn "inventory_money_tracking_software/cmd/common"
	mdls "inventory_money_tracking_software/cmd/models"
)

func selectItem[T sqlConection](conn T, itemId int) (mdls.Item, error) {
	var item mdls.Item
	item.Id = itemId
	stmt := readSql.readItem
	var unit sql.NullString
	var count sql.NullFloat64
	err := conn.QueryRow(stmt, item.Id).Scan(&item.Name,
		&item.Description, &item.Imagelink, &count, &unit)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	if unit.Valid {
		item.InventoryUnit = unit.String
	}
	if count.Valid {
		item.Count = count.Float64
	}
	return item, err
}

// updates itemQuantity as part of the inventory journal entry
func updateItemCount(tx *sql.Tx, itemId int, newCount float64) error {
	_, err := tx.Exec(updateSql.updateItemCount, newCount, itemId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}

// selects storehosue by id, accepts either a transaction or a normal db coneection
func selectStorehouse[T sqlConection](conn T, id int) (mdls.Storehouse, error) {
	s := mdls.Storehouse{}
	s.Id = id
	err := conn.QueryRow(readSql.readStorehouse, id).Scan(&s.Name, &s.Description, &s.StorehouseType)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return s, err
}

// selects storehosue by name, accepts either a transaction or a normal db coneection
func selectStorehouseByName[T sqlConection](conn T, name string) (mdls.Storehouse, error) {
	s := mdls.Storehouse{Name: name}
	err := conn.QueryRow(readSql.readStorehouseByName, name).Scan(&s.Id,
		&s.Description, &s.StorehouseType)
	if err == sql.ErrNoRows {
		cmn.HandleError(err, cmn.ErrorLevels.Operation)
	} else {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
	}
	return s, err
}

func insertStorehouse[T sqlConection](conn T, s mdls.Storehouse) (int, error) {
	cmn.Log("insert new storehouse to database", cmn.LogLevels.Operation)
	var sId int
	stmt := insertSql.insertStorehouse
	err := conn.QueryRow(stmt, s.Name, s.Description, string(s.StorehouseType)).Scan(&sId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return sId, err
}

// inserts item without the unit, unit has to be done in a seperate
func insertItem(i mdls.Item) (int, error) {
	cmn.Log("insert new Item to database", cmn.LogLevels.Operation)
	var itemId int
	stmnt := insertSql.insertItem
	err := db.QueryRow(stmnt, i.Name, i.Description, i.Imagelink, i.InventoryUnit).Scan(&itemId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return itemId, err
}

// updates the item, but doesn't allow you to update the count
// the item count must only be updated within an inventory journal entry
func updateItem(i mdls.Item) error {
	cmn.Log("update item id %d ", cmn.LogLevels.Operation, i.Id)
	stmnt := updateSql.updateItem
	_, err := db.Exec(stmnt, i.Name, i.Description, i.Imagelink,
		i.InventoryUnit)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}

// deletion needes to consider deletion from the
// storehouse_item junction table too
func deleteItem(tx *sql.Tx, i mdls.Item) error {
	// cmn.Log("delete item id %d ", cmn.LogLevels.Operation, i.Id)
	// stmnt := deleteSql.deleteItem
	// var itemId string
	// err := db.QueryRow(stmnt, i.Name, i.Description, i.Imagelink,
	// 	i.InventoryUnitId).Scan(&itemId)
	// cmn.HandleError(err, cmn.ErrorLevels.Error)
	return nil
}

func insertInventoryJournalEntry(tx *sql.Tx, ij mdls.InventoryJournalEntry) (int, error) {
	cmn.Log("insertInventoryJournalEntry", cmn.LogLevels.Operation)
	s := insertSql.insertInventoryJournalEntry
	id := 0
	err := tx.QueryRow(s, ij.Description).Scan(&id)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return id, err
}

func insertInventoryMovement(tx *sql.Tx, m mdls.InventoryMovement) error {
	cmn.Log("insertInventoryMovement with movetype %v storehouseId %d, itemId %d and quantity %f",
		cmn.LogLevels.Operation,
		m.MoveType, m.StorehouseId, m.ItemId, m.Quantity)
	stmt := insertSql.insertInventoryMovement
	t := string(m.MoveType)
	_, err := tx.Exec(stmt, &m.StorehouseId, &m.ItemId, &m.Quantity, &t,
		&m.InventoryJournalEntryId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}

// selects storehouseItemQuantity, accepts either a transaction or a normal db coneection
func selectStorehouseItemQuantity[T sqlConection](conn T, itemId,
	storehouseId int) (mdls.StorehouseItemQuantity, error) {
	q := mdls.StorehouseItemQuantity{
		ItemId:       itemId,
		StorehouseId: storehouseId,
	}
	stmt := readSql.readStorehouseItemQuantity
	err := conn.QueryRow(stmt, storehouseId, itemId).Scan(&q.Quantity)
	if err != nil && err == sql.ErrNoRows {
		return q, err
	} else if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return q, err
	}
	return q, nil
}

func insertStorehouseItemQuantity(tx *sql.Tx, q mdls.StorehouseItemQuantity) (int, error) {
	cmn.Log("insertStorehouseItemQuantity with storehouseId %d, itemId %d and quantity %f",
		cmn.ErrorLevels.Operation,
		q.StorehouseId, q.ItemId, q.Quantity)
	if q.StorehouseId == 0 || q.ItemId == 0 {
		err := errors.New("Either StorehouseId or ItemId is zero")
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return 0, err
	}
	s := insertSql.insertStorehouseItemQuantity
	_, err := tx.Exec(s, q.StorehouseId, q.ItemId, q.Quantity)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return q.Id, err
}
func updateStorehouseItemQuantity(tx *sql.Tx, q mdls.StorehouseItemQuantity) error {
	if q.StorehouseId == 0 || q.ItemId == 0 {
		err := errors.New("Either StorehouseId or ItemId is zero")
		cmn.HandleError(err, cmn.ErrorLevels.Error)

		return err
	}
	s := updateSql.updateStorehouseItemQuantity
	_, err := tx.Exec(s, q.Quantity, q.StorehouseId, q.ItemId)
	cmn.HandleError(err, cmn.ErrorLevels.Error)
	return err
}

func selectInventoryMovement(id int) (mdls.InventoryMovement, error) {
	s := readSql.readInventoryMovement
	m := mdls.InventoryMovement{
		Id: id,
	}
	err := db.QueryRow(s, id).Scan(&m.StorehouseId, &m.ItemId, &m.Quantity, &m.MoveType,
		&m.InventoryJournalEntryId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return m, err
	}
	return m, nil
}

func selectInventoryMovementByInventoryJournalEntryId(inventoryJournalEntryId int) ([]mdls.InventoryMovement,
	error) {
	o := []mdls.InventoryMovement{}
	if inventoryJournalEntryId == 0 {
		err := errors.New("Inventory journal ID cant be zero...")
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return o, err
	}
	s := readSql.readInventoryMovementByInventoryJournalId
	rows, err := db.Query(s, inventoryJournalEntryId)
	if err != nil {
		cmn.HandleError(err, cmn.ErrorLevels.Error)
		return o, err
	}
	for rows.Next() {
		m := mdls.InventoryMovement{InventoryJournalEntryId: inventoryJournalEntryId}
		err := rows.Scan(&m.Id, &m.StorehouseId, &m.ItemId, &m.Quantity, &m.MoveType)
		if err != nil {
			cmn.HandleError(err, cmn.ErrorLevels.Error)
			return o, err
		}
		o = append(o, m)
	}
	return o, nil
}
