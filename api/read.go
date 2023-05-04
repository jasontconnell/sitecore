package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"github.com/jasontconnell/sitecore/scprotobuf"
	"github.com/jasontconnell/sqlhelp"
	"google.golang.org/protobuf/proto"
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

	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	rserr := sqlhelp.GetResultsChannel(conn, sqlstr, rchan)

	if rserr != nil {
		return nil, rserr
	}

	if c <= 0 {
		c = 12
	}
	fvchan := make(chan data.FieldValueNode, 20000000)

	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func(id int, records chan map[string]interface{}, fv chan data.FieldValueNode) {
			count := 0

			for row := range records {
				fieldValue := data.NewFieldValue(
					getUUID(row["FieldID"]),
					getUUID(row["ItemID"]),
					row["Name"].(string),
					row["Value"].(string),
					data.GetLanguage(row["Language"].(string)),
					row["Version"].(int64),
					data.GetSource(row["Source"].(string)),
				)
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

func LoadFieldValuesMetadata(connstr string, c int) ([]data.FieldValueNode, error) {
	sqlstr := `
	with FieldValues (ValueID, ItemID, FieldID, Version, Language, Source)
	as
	(
		select
			ID, ItemId, FieldId, 1, 'en', 'SharedFields'
		from SharedFields
		union
		select
			ID, ItemId, FieldId, Version, Language, 'VersionedFields'
		from VersionedFields
		union
		select
			ID, ItemId, FieldId, 1, Language, 'UnversionedFields'
		from UnversionedFields
	)
	select 
		cast(fv.ValueID as char(36)) as ValueID, 
		cast(fv.ItemID as char(36)) as ItemID, 
		f.Name, 
		cast(fv.FieldID as char(36)) as FieldID, 
		fv.Version, 
		fv.Language, 
		fv.Source
	from
		FieldValues fv
			join Items f
				on fv.FieldID = f.ID
	order by fv.Source, f.Name, fv.Language, fv.Version;
`
	rchan := make(chan map[string]interface{}, 500000)

	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	rserr := sqlhelp.GetResultsChannel(conn, sqlstr, rchan)

	if rserr != nil {
		return nil, rserr
	}

	if c <= 0 {
		c = 12
	}
	fvchan := make(chan data.FieldValueNode, 20000000)

	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func(id int, records chan map[string]interface{}, fv chan data.FieldValueNode) {

			for row := range records {
				fieldValue := data.NewFieldValue(
					getUUID(row["FieldID"]),
					getUUID(row["ItemID"]),
					row["Name"].(string),
					"",
					data.GetLanguage(row["Language"].(string)),
					row["Version"].(int64),
					data.GetSource(row["Source"].(string)),
				)

				fv <- fieldValue
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

func LoadFilteredFieldValues(connstr string, fieldIds []uuid.UUID, c int) ([]data.FieldValueNode, error) {
	if len(fieldIds) == 0 {
		return LoadFieldsParallel(connstr, c)
	}

	filters := []string{}
	for _, fieldId := range fieldIds {
		filters = append(filters, "'"+fieldId.String()+"'")
	}
	filter := strings.Join(filters, ",")
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
		where fv.FieldID in (%s)
        order by fv.Source, f.Name, fv.Language, fv.Version;
	`

	query := fmt.Sprintf(sqlstr, filter)

	rchan := make(chan map[string]interface{}, 500000)

	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	rserr := sqlhelp.GetResultsChannel(conn, query, rchan)

	if rserr != nil {
		return nil, rserr
	}

	if c <= 0 {
		c = 12
	}
	fvchan := make(chan data.FieldValueNode, 20000000)

	var wg sync.WaitGroup
	for i := 0; i < c; i++ {
		wg.Add(1)
		go func(id int, records chan map[string]interface{}, fv chan data.FieldValueNode) {

			for row := range records {
				fieldValue := data.NewFieldValue(
					getUUID(row["FieldID"]),
					getUUID(row["ItemID"]),
					row["Name"].(string),
					row["Value"].(string),
					data.GetLanguage(row["Language"].(string)),
					row["Version"].(int64),
					data.GetSource(row["Source"].(string)),
				)

				fv <- fieldValue
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

func loadTemplatesFromDb(connstr string) ([]*data.TemplateQueryRow, error) {
	query := fmt.Sprintf(itemSelect, `isnull(sf.Value, '') as Type, isnull(Replace(Replace(UPPER(b.Value), '}',''), '{', ''), '') as BaseTemplates, isnull(Replace(Replace(UPPER(sv.Value), '}',''), '{', ''), '') as StandardValuesId, isnull(sh.Value, '0') as Shared, isnull(unv.Value, '0') as Unversioned`,
		`left join SharedFields sf
                    on i.ID = sf.ItemId
                        and sf.FieldId = 'AB162CC0-DC80-4ABF-8871-998EE5D7BA32'
                left join SharedFields b
                    on i.ID = b.ItemID
						and b.FieldId = '12C33F3F-86C5-43A5-AEB4-5598CEC45116'
				left join SharedFields sv
						on i.ID = sv.ItemID
							and sv.FieldId = 'F7D48A55-2158-4F02-9356-756654404F73'
				left join SharedFields sh
					on i.ID = sh.ItemID
						and sh.FieldId = 'BE351A73-FCB0-4213-93FA-C302D8AB4F51'
				left join SharedFields unv
					on i.ID = unv.ItemID
						and unv.FieldId = '39847666-389D-409B-95BD-F2016F11EED5'
						`)

	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	records, rerr := sqlhelp.GetResultSet(conn, query)

	if rerr != nil {
		return nil, rerr
	}

	var rows []*data.TemplateQueryRow
	for _, row := range records {
		tmp := &data.TemplateQueryRow{
			ID:               getUUID(row["ID"]),
			Name:             row["Name"].(string),
			TemplateID:       getUUID(row["TemplateID"]),
			ParentID:         getUUID(row["ParentID"]),
			BaseTemplateIds:  getUUIDs(row["BaseTemplates"], "|"),
			StandardValuesId: getUUID(row["StandardValuesId"]),
			Type:             row["Type"].(string),
			Shared:           row["Shared"].(string),
			Unversioned:      row["Unversioned"].(string),
		}
		// inner := data.NewItemNode(getUUID(row["ID"]), row["Name"].(string), getUUID(row["TemplateID"]), getUUID(row["ParentID"]), getUUID(row["MasterID"]))
		// tmp := data.NewTemplateNode(inner, row["Type"].(string), getUUIDs(row["BaseTemplates"], "|"))

		rows = append(rows, tmp)
	}

	return rows, nil
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

func ReadProtobuf(filename string) ([]data.ItemNode, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't read items from protobuf file %s. %w", filename, err)
	}

	list := []data.ItemNode{}

	var items scprotobuf.ItemsData
	err = proto.Unmarshal(b, &items)
	if err != nil {
		return nil, fmt.Errorf("couldn't deserialize items from protobuf file %s. %w", filename, err)
	}

	nmap := make(map[uuid.UUID]string)
	for _, pitem := range items.ItemDefinitions {
		n := data.NewItemNode(
			getUUIDFromProtoGuid(pitem.ID),
			pitem.Item.Name,
			getUUIDFromProtoGuid(pitem.Item.TemplateID),
			getUUIDFromProtoGuid(pitem.Item.ParentID),
			getUUIDFromProtoGuid(pitem.Item.MasterID),
		)

		nmap[n.GetId()] = n.GetName()

		list = append(list, n)
	}

	m := make(map[uuid.UUID][]data.FieldValueNode)
	for _, pfld := range items.LanguageData {
		id := getUUIDFromProtoGuid(pfld.ID)

		var flist []data.FieldValueNode
		for _, ld := range pfld.LanguageData {
			for _, v := range ld.VersionsData {
				for _, f := range v.Fields {
					if f.Value == "" {
						continue
					}

					fieldID := getUUIDFromProtoGuid(f.ID)
					name := nmap[fieldID]

					fv := data.NewFieldValue(fieldID, id, name, f.Value, data.Language(ld.Language), int64(v.Version), data.VersionedFields)
					flist = append(flist, fv)
				}
			}
		}

		m[id] = flist
	}

	for _, sfld := range items.SharedData {
		id := getUUIDFromProtoGuid(sfld.ID)

		for _, fld := range sfld.SharedDataItems {
			fieldID := getUUIDFromProtoGuid(fld.ID)
			name := nmap[fieldID]
			fv := data.NewFieldValue(fieldID, id, name, fld.Value, data.English, 1, data.SharedFields)

			m[id] = append(m[id], fv)
		}
	}

	for i := 0; i < len(list); i++ {
		flds := m[list[i].GetId()]

		for _, f := range flds {
			list[i].AddFieldValue(f)
		}
	}

	return list, nil
}

func getUUIDFromProtoGuid(g *scprotobuf.Guid) uuid.UUID {
	var plo, phi uint64
	if g != nil {
		plo = *g.Lo
		phi = *g.Hi
	}

	return MustParseUUIDProto(plo, phi)
}
