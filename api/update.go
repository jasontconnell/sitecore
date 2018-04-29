package api

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
    "sitecore/data"
    "sync"
    "strings"
)

func Update(connstr string, items []data.UpdateItem, fields []data.UpdateField) int64 {
	itemGroups := 2
	itemGroupSize := len(items)/itemGroups + 1

	if len(items) < 100 {
		itemGroups = 1
		itemGroupSize = len(items)
	}

	fieldGroups := 4
	fieldGroupSize := len(fields)/fieldGroups + 1
	if len(fields) < 100 {
		fieldGroups = 1
		fieldGroupSize = len(fields)
	}

	var updated int64 = 0
	var wg sync.WaitGroup

	if len(items) > 0 {
		wg.Add(itemGroups)
		// items - 2 processes
		for i := 0; i < itemGroups; i++ {
			grplist := items[i*itemGroupSize : (i+1)*itemGroupSize]
			go func(grp []data.UpdateItem) {
				updated += updateItems(connstr, grp)
				wg.Done()
			}(grplist)
		}
	} else {
		wg.Done()
	}

	if len(fields) > 0 {
		wg.Add(fieldGroups)
		// fields - 4 processes
		for i := 0; i < fieldGroups; i++ {
			grplist := fields[i*fieldGroupSize : (i+1)*fieldGroupSize]
			go func(grp []data.UpdateField) {
				updated += updateFields(connstr, grp)
				wg.Done()
			}(grplist)
		}
	} else {
		wg.Done()
	}

	wg.Wait()

	return updated
}

func updateItems(connstr string, items []data.UpdateItem) int64 {
	var updated int64 = 0
	if db, err := sql.Open("mssql", connstr); err == nil {
		defer db.Close()

		for _, sql := range getSqlForItems(items) {
			if result, err := db.Exec(sql); err == nil {
				i, _ := result.RowsAffected()
				updated += i
			} else {
				fmt.Println(err)
				return -1
			}
		}
	}

	return updated
}

func updateFields(connstr string, fields []data.UpdateField) int64 {
	var updated int64 = 0
	if db, err := sql.Open("mssql", connstr); err == nil {
		defer db.Close()

		for _, sql := range getSqlForFields(fields) {
			if result, err := db.Exec(sql); err == nil {
				i, _ := result.RowsAffected()
				updated += i
			} else {
				fmt.Println(err)
				return -1
			}
		}
	}
	return updated
}

var updateitemfmt string = "update Items set Name = '%[1]v', TemplateID = '%[2]v', ParentID = '%[3]v', MasterID = '%[4]v' where ID = '%[5]v'"
var insertitemfmt string = "insert into Items (ID, Name, TemplateID, ParentID, MasterID, Created, Updated) values ('%[5]v', '%[1]v', '%[2]v', '%[3]v', '%[4]v', getdate(), getdate())"
var deleteitemfmt string = "delete from Items where ID = '%v'"

func getSqlForItems(items []data.UpdateItem) []string {
	sqllist := []string{}
	for _, item := range items {
		var sql string
		switch item.UpdateType {
		case data.Update:
			sql = fmt.Sprintf(updateitemfmt, item.Name, item.TemplateID, item.ParentID, item.MasterID, item.ID)
		case data.Insert:
			sql = fmt.Sprintf(insertitemfmt, item.Name, item.TemplateID, item.ParentID, item.MasterID, item.ID)
		case data.Delete:
			sql = fmt.Sprintf(deleteitemfmt, item.ID)
		case data.Ignore:
			sql = ""
		}

		if len(sql) > 0 {
			sqllist = append(sqllist, sql)
		}
	}
	return sqllist
}

func getSqlForFields(fields []data.UpdateField) []string {
	updatemap := make(map[string]string)
	insertmap := make(map[string]string)
	deletemap := make(map[string]string)

	updatemap["SharedFields"] = "update %[1]v set Value = '%[4]v', Updated = getdate() where ItemID = '%[2]v' and FieldID = '%[3]v'"
	updatemap["UnversionedFields"] = "update %[1]v set Value = '%[4]v', Updated = getdate() where ItemID = '%[2]v' and FieldID = '%[3]v' and Language = '%[5]v'"
	updatemap["VersionedFields"] = "update %[1]v set Value = '%[4]v', Updated = getdate() where ItemID = '%[2]v' and FieldID = '%[3]v' and Language = '%[5]v' and Version = %[6]v"

	insertmap["SharedFields"] = "insert into %[1]v (ID, ItemID, FieldID, Value, Created, Updated) values (newid(), '%[2]v', '%[3]v', '%[4]v', getdate(), getdate())"
	insertmap["UnversionedFields"] = "insert into %[1]v (ID, ItemID, FieldID, Value, Language, Created, Updated) values (newid(), '%[2]v', '%[3]v', '%[4]v', '%[5]v', getdate(), getdate())"
	insertmap["VersionedFields"] = "insert into %[1]v (ID, ItemID, FieldID, Value, Language, Version, Created, Updated) values (newid(), '%[2]v', '%[3]v', '%[4]v', '%[5]v', '%[6]v', getdate(), getdate())"

	deletemap["SharedFields"] = "delete from %[1]v where ItemID = '%[2]v' and FieldID = '%[3]v'"
	deletemap["UnversionedFields"] = "delete from %[1]v where ItemID = '%[2]v' and FieldID = '%[3]v' and Language = '%[5]v'"
	deletemap["VersionedFields"] = "delete from %[1]v where ItemID = '%[2]v' and FieldID = '%[3]v' and Language = '%[5]v' and Version = %[6]v"

	sqllist := []string{}
	for _, field := range fields {
		var sql string
		value := strings.Replace(field.Value, "'", "''", -1)

		switch field.UpdateType {
		case data.Update:
			sqlfmt, _ := updatemap[field.Source]
			sql = fmt.Sprintf(sqlfmt, field.Source, field.ItemID, field.FieldID, value, field.Language, field.Version)
		case data.Insert:
			sqlfmt, _ := insertmap[field.Source]
			sql = fmt.Sprintf(sqlfmt, field.Source, field.ItemID, field.FieldID, value, field.Language, field.Version)
		case data.Delete:
			sqlfmt, _ := deletemap[field.Source]
			sql = fmt.Sprintf(sqlfmt, field.Source, field.ItemID, field.FieldID, value, field.Language, field.Version)
		case data.Ignore:
			sql = ""
		}

		if len(sql) > 0 {
			sqllist = append(sqllist, sql)
		}
	}
	return sqllist
}

func cleanOrphanedItems(connstr string) (rows int64) {
	subq := "select ID from Items where ParentID not in (select ID from Items) and ParentID <> '00000000-0000-0000-0000-000000000000'"
	sqlfmt := `
        delete from SharedFields where ItemID in ( %[1]v )
        delete from VersionedFields where ItemID in ( %[1]v )
        delete from UnversionedFields where ItemID in ( %[1]v )
        delete from Items where ID in ( %[1]v )
    `

	sqlstr := fmt.Sprintf(sqlfmt, subq)

	if db, err := sql.Open("mssql", connstr); err == nil {
		defer db.Close()

		if result, err := db.Exec(sqlstr); err == nil {
			rows, _ = result.RowsAffected()
		} else {
			fmt.Println("cleaning orphaned items", err)
		}
	}

	return
}
