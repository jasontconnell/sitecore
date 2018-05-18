package api

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"github.com/jasontconnell/sqlhelp"
	"strings"
	"sync"
)

var emptyUuid uuid.UUID = MustParseUUID("00000000-0000-0000-0000-000000000000")

const itemSelect = `select 
            cast(i.ID as char(36)) ID, 
            i.Name, 
            cast(i.TemplateID as char(36)) as TemplateID, 
            cast(i.ParentID as char(36)) as ParentID, 
            cast(i.MasterID as char(36)) as MasterID, 
            i.Created, 
            i.Updated,
            %v
        from Items i %v
        order by i.Name`

func LoadItems(connstr string) ([]data.ItemNode, error) {
	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	query := fmt.Sprintf(itemSelect, "0", "")
	records, rerr := sqlhelp.GetResultSet(conn, query)

	if rerr != nil {
		return nil, rerr
	}

	var items []data.ItemNode
	for _, row := range records {
		item := data.NewItemNode(getUUID(row["ID"]), row["Name"].(string), getUUID(row["TemplateID"]), getUUID(row["ParentID"]), getUUID(row["MasterID"]))
		items = append(items, item)
	}

	return items, nil
}

func LoadFields(connstr string) ([]data.FieldValueNode, error) {
	return LoadFieldsParallel(connstr, 1)
}

// Load Fields can return a ton of results. Pass in 'c' to specify how many goroutines should be spawned
func LoadFieldsParallel(connstr string, c int) ([]data.FieldValueNode, error) {
	if c <= 0 {
		c = 12
	}
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

	fvchan := make(chan data.FieldValueNode, 20000000)

	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func(id int, records chan map[string]interface{}, fv chan data.FieldValueNode) {
			count := 0

			for row := range records {
				fieldValue := &data.FieldValue{
					FieldValueID: getUUID(row["ValueID"]),
					ItemID:       getUUID(row["ItemID"]),
					Name:         row["Name"].(string),
					FieldID:      getUUID(row["FieldID"]),
					Value:        row["Value"].(string),
					Language:     row["Language"].(string),
					Version:      row["Version"].(int64),
					Source:       row["Source"].(string),
				}
				fv <- fieldValue
				count++
			}
			wg.Done()
		}(i, rchan, fvchan)
	}

	wg.Wait()

	close(fvchan)

	wg.Add(1)
	fieldValues := []data.FieldValueNode{}
	go func(fv chan data.FieldValueNode) {
		for fieldValue := range fvchan {
			fieldValues = append(fieldValues, fieldValue)
		}
		wg.Done()
	}(fvchan)

	wg.Wait()

	return fieldValues, nil
}

func loadTemplatesFromDb(connstr string) ([]data.TemplateNode, error) {
	query := fmt.Sprintf(itemSelect, `isnull(sf.Value, '') as Type, isnull(Replace(Replace(UPPER(b.Value), '}',''), '{', ''), '') as BaseTemplates`,
		`left join SharedFields sf
                    on i.ID = sf.ItemId
                        and sf.FieldId = 'AB162CC0-DC80-4ABF-8871-998EE5D7BA32'
                left join SharedFields b
                    on i.ID = b.ItemID
                        and b.FieldId = '12C33F3F-86C5-43A5-AEB4-5598CEC45116'`)

	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	records, rerr := sqlhelp.GetResultSet(conn, query)

	if rerr != nil {
		return nil, rerr
	}

	var items []data.TemplateNode
	for _, row := range records {
		inner := data.NewItemNode(getUUID(row["ID"]), row["Name"].(string), getUUID(row["TemplateID"]), getUUID(row["ParentID"]), getUUID(row["MasterID"]))
		tmp := data.NewTemplateNode(inner, row["Type"].(string), getUUIDs(row["BaseTemplates"], "|"))

		items = append(items, tmp)
	}

	return items, nil
}

func getUUIDs(val interface{}, splitchar string) []uuid.UUID {
	if val == nil {
		return nil
	}

	s := val.(string)
	ss := strings.Split(s, splitchar)
	list := []uuid.UUID{}
	for _, id := range ss {
		u := getUUID(id)
		list = append(list, u)
	}
	return list
}

func getUUID(val interface{}) uuid.UUID {
	if val == nil {
		return emptyUuid
	}
	id, iderr := TryParseUUID(val.(string))
	if iderr != nil {
		id = emptyUuid
	}

	return id
}
