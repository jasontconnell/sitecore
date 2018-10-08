package api

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jasontconnell/sitecore/data"
	"strings"
	"sync"
)

func Update(connstr string, items []data.UpdateItem, fields []data.UpdateField) (int64, error) {
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
	var processError error

	if len(items) > 0 {
		wg.Add(itemGroups)
		// items - 2 processes
		for i := 0; i < itemGroups; i++ {
			last := (i + 1) * itemGroupSize
			if last > len(items) {
				last = len(items)
			}

			grplist := items[i*itemGroupSize : last]
			go func(grp []data.UpdateItem) {
				defer wg.Done()

				up, err := updateItems(connstr, grp)
				if err != nil {
					processError = err
					return
				}

				updated += up
			}(grplist)
		}
	}

	if len(fields) > 0 {
		wg.Add(fieldGroups)
		// fields - 4 processes
		for i := 0; i < fieldGroups; i++ {
			last := (i + 1) * fieldGroupSize
			if last > len(fields) {
				last = len(fields)
			}
			grplist := fields[i*fieldGroupSize : last]
			go func(grp []data.UpdateField) {
				defer wg.Done()
				up, err := updateFields(connstr, grp)
				if err != nil {
					processError = err
					return
				}
				updated += up
			}(grplist)
		}
	}

	wg.Wait()

	return updated, processError
}

func updateItems(connstr string, items []data.UpdateItem) (int64, error) {
	var updated int64 = 0
	if db, err := sql.Open("mssql", connstr); err == nil {
		defer db.Close()

		for _, sql := range getSqlForItems(items) {
			if result, err := db.Exec(sql); err == nil {
				i, _ := result.RowsAffected()
				updated += i
			} else {
				return -1, fmt.Errorf("Updating items, encountered an error, %v", err)
			}
		}
	}

	return updated, nil
}

func updateFields(connstr string, fields []data.UpdateField) (int64, error) {
	var updated int64 = 0
	if db, err := sql.Open("mssql", connstr); err == nil {
		defer db.Close()

		for _, sql := range getSqlForFields(fields) {
			if result, err := db.Exec(sql); err == nil {
				i, _ := result.RowsAffected()
				updated += i
			} else {
				return -1, fmt.Errorf("Updating fields, encountered an error, %v", err)
			}
		}
	}
	return updated, nil
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

		prms := []interface{}{field.Source, field.ItemID, field.FieldID, value}
		if field.Source == "UnversionedFields" {
			prms = append(prms, field.Language)
		}

		if field.Source == "VersionedFields" {
			prms = append(prms, field.Language)
			prms = append(prms, field.Version)
		}

		switch field.UpdateType {
		case data.Update:
			sqlfmt, _ := updatemap[field.Source]
			sql = fmt.Sprintf(sqlfmt, prms...)
		case data.Insert:
			sqlfmt, _ := insertmap[field.Source]
			sql = fmt.Sprintf(sqlfmt, prms...)
		case data.Delete:
			sqlfmt, _ := deletemap[field.Source]
			sql = fmt.Sprintf(sqlfmt, prms...)
		case data.Ignore:
			sql = ""
		}

		if len(sql) > 0 {
			sqllist = append(sqllist, sql)
		}
	}
	return sqllist
}

func CleanOrphanedItems(connstr string) (rows int64, err error) {
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
			return -1, fmt.Errorf("Cleaning up orphaned items, encountered an error: %v", err)
		}
	}

	return
}
