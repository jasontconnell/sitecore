package item

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"github.com/jasontconnell/sqlhelp"
	"time"
    "sync"
)

var emptyUuid uuid.UUID = uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000"))

func LoadItems(connstr string) ([]*ItemNode, error) {
	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	sqlstr := `
        select 
            cast(ID as char(36)) as ID, 
            Name, 
            cast(TemplateID as char(36)) as TemplateID, 
            cast(ParentID as char(36)) as ParentID, 
            cast(MasterID as char(36)) as MasterID, 
            Created, 
            Updated 
        from Items 
        order by Name`
	records, rerr := sqlhelp.GetResultSet(conn, sqlstr)

	if rerr != nil {
		return nil, rerr
	}

	items := []*ItemNode{}
	for _, row := range records {
		item := Item{
			ID:         getUUID(row["ID"]),
			Name:       row["Name"].(string),
			TemplateID: getUUID(row["TemplateID"]),
			ParentID:   getUUID(row["ParentID"]),
			MasterID:   getUUID(row["MasterID"]),
			Created:    row["Created"].(time.Time),
			Updated:    row["Updated"].(time.Time),
		}

		items = append(items, &ItemNode{Item: item})
	}

	return items, nil
}

func LoadFields(connstr string) ([]*FieldValue, error) {
    conn, cerr := sql.Open("mssql", connstr)
    if cerr != nil {
        return nil, cerr
    }

    defer conn.Close()

    sqlstr := `
        with FieldValues (ValueID, ItemID, FieldID, Value, Version, Language, Source)
        as
        (
            select
                ID, ItemId, FieldId, Value, 1, 'en', 'SharedFields'
            from SharedFields
            union
            select
                ID, ItemId, FieldId, Value, Version, Language, 'VersionedFields'
            from VersionedFields
            union
            select
                ID, ItemId, FieldId, Value, 1, Language, 'UnversionedFields'
            from UnversionedFields
        )
        select 
            cast(fv.ValueID as char(36)) as ValueID, 
            cast(fv.ItemID as char(36)) as ItemID, 
            f.Name, 
            cast(fv.FieldID as char(36)) as FieldID, 
            fv.Value, fv.Version, 
            fv.Language, 
            fv.Source
        from
            FieldValues fv
                join Items f
                    on fv.FieldID = f.ID
        order by fv.Source, f.Name, fv.Language, fv.Version;
    `

    records, rserr := sqlhelp.GetResultSet(conn, sqlstr)

    if rserr != nil {
        return nil, rserr
    }

    fieldValues := []*FieldValue{}

    for _, row := range records {
        fieldValue := &FieldValue{
            FieldValueID: getUUID(row["ValueID"]),
            ItemID:       getUUID(row["ItemID"]),
            Name:         row["Name"].(string),
            FieldID:      getUUID(row["FieldID"]),
            Value:        row["Value"].(string),
            Language:     row["Language"].(string),
            Version:      row["Version"].(int64),
            Source:       row["Source"].(string),
        }
        fieldValues = append(fieldValues, fieldValue)
    }

    return fieldValues, nil
}

// Load Fields can return a ton of results. Pass in 'c' to specify how many goroutines should be spawned
func LoadFieldsParallel(connstr string, c int) ([]*FieldValue, error) {
    if c <= 0 { c = 24 }
	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	sqlstr := `
        with FieldValues (ValueID, ItemID, FieldID, Value, Version, Language, Source)
        as
        (
            select
                ID, ItemId, FieldId, Value, 1, 'en', 'SharedFields'
            from SharedFields
            union
            select
                ID, ItemId, FieldId, Value, Version, Language, 'VersionedFields'
            from VersionedFields
            union
            select
                ID, ItemId, FieldId, Value, 1, Language, 'UnversionedFields'
            from UnversionedFields
        )
        select 
            cast(fv.ValueID as char(36)) as ValueID, 
            cast(fv.ItemID as char(36)) as ItemID, 
            f.Name, 
            cast(fv.FieldID as char(36)) as FieldID, 
            fv.Value, fv.Version, 
            fv.Language, 
            fv.Source
        from
            FieldValues fv
                join Items f
                    on fv.FieldID = f.ID
        order by fv.Source, f.Name, fv.Language, fv.Version;
    `

    rchan := make(chan map[string]interface{}, 500000)
	rserr := sqlhelp.GetResultsChannel(conn, sqlstr, rchan)

    if rserr != nil {
        return nil, rserr
    }

    fvchan := make(chan *FieldValue, 5000000)

    var wg sync.WaitGroup
    for i := 0; i < c; i++ {
        wg.Add(1)
        go func(id int, records chan map[string]interface{}, fv chan *FieldValue){
            count := 0

            for row := range records {
                fieldValue := &FieldValue{
                    FieldValueID: getUUID(row["ValueID"]),
                    ItemID:       getUUID(row["ItemID"]),
                    Name:         row["Name"].(string),
                    FieldID:      getUUID(row["FieldID"]),
                    Value:        row["Value"].(string),
                    Language:     row["Language"].(string),
                    Version:      row["Version"].(int64),
                    Source:       row["Source"].(string),
                }
                //fieldValues = append(fieldValues, fieldValue)
                fv <- fieldValue
                count++
            }
            wg.Done()
        }(i, rchan, fvchan)
    }

    wg.Wait()
    close(fvchan)

    wg.Add(1)
    fieldValues := []*FieldValue{}
    go func(fv chan *FieldValue){
        for fieldValue := range fvchan {
            fieldValues = append(fieldValues, fieldValue)
        }
        wg.Done()
    }(fvchan)

    wg.Wait()
    
	return fieldValues, nil
}

func getUUID(val interface{}) uuid.UUID {
    if val == nil { return emptyUuid }
	id, iderr := uuid.Parse(val.(string))
	if iderr != nil {
		id = emptyUuid
	}

	return id
}
