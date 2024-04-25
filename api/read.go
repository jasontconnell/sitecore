package api

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/api/queries"
	"github.com/jasontconnell/sitecore/data"
	"github.com/jasontconnell/sitecore/scprotobuf"
	"github.com/jasontconnell/sqlhelp"
	_ "github.com/microsoft/go-mssqldb"
	"google.golang.org/protobuf/proto"
)

var emptyUuid uuid.UUID = MustParseUUID("00000000-0000-0000-0000-000000000000")

func LoadItems(connstr string) ([]data.ItemNode, error) {
	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}

	defer conn.Close()

	records, rerr := sqlhelp.GetResultSet(conn, queries.Items)

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

func LoadItemsByTemplates(connstr string, templateIds []uuid.UUID) ([]data.ItemNode, error) {
	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}
	defer conn.Close()

	idstr := []string{}
	for _, id := range templateIds {
		idstr = append(idstr, "'"+id.String()+"'")
	}

	query := fmt.Sprintf(queries.ItemsByTemplate, strings.Join(idstr, ","))
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
func processFieldValueQuery(connstr string, query string, c int) ([]data.FieldValueNode, error) {
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
			count := 0

			for row := range records {
				fieldValue := data.NewFieldValue(
					getUUID(row["FieldID"]),
					getUUID(row["ItemID"]),
					row["Name"].(string),
					row["Value"].(string),
					data.GetLanguage(row["Language"].(string)),
					row["Version"].(int64),
					row["Created"].(time.Time),
					row["Updated"].(time.Time),
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

func LoadFieldsParallel(connstr string, c int) ([]data.FieldValueNode, error) {
	return processFieldValueQuery(connstr, queries.FieldValues, c)
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

	query := fmt.Sprintf(queries.FieldValuesByField, filter)
	return processFieldValueQuery(connstr, query, c)
}

func LoadFieldValuesTemplates(connstr string, fieldIds, templateIds []uuid.UUID, c int) ([]data.FieldValueNode, error) {
	if len(fieldIds) == 0 {
		return LoadFieldsParallel(connstr, c)
	}

	if len(templateIds) == 0 {
		return LoadFilteredFieldValues(connstr, fieldIds, c)
	}

	fields := []string{}
	for _, fieldId := range fieldIds {
		fields = append(fields, "'"+fieldId.String()+"'")
	}
	fieldFilter := strings.Join(fields, ",")

	templates := []string{}
	for _, tmpId := range templateIds {
		templates = append(templates, "'"+tmpId.String()+"'")
	}
	templateFilter := strings.Join(templates, ",")

	query := fmt.Sprintf(queries.FieldValuesByFieldAndTemplate, fieldFilter, templateFilter)
	return processFieldValueQuery(connstr, query, c)
}

func loadTemplatesFromDb(connstr string) ([]*data.TemplateQueryRow, error) {
	rootIdStr := data.TemplatesRootID.String()
	query := fmt.Sprintf(queries.TemplatesByRoot, rootIdStr, rootIdStr)

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
			StandardValuesId: getUUID(row["StandardValuesField"]),
			Type:             row["Type"].(string),
			Shared:           row["Shared"].(string),
			Unversioned:      row["Unversioned"].(string),
		}

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
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't stat file %s. %w", filename, err)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't read items from protobuf file %s. %w", filename, err)
	}

	var items scprotobuf.ItemsData
	err = proto.Unmarshal(b, &items)
	if err != nil {
		return nil, fmt.Errorf("couldn't deserialize items from protobuf file %s. %w", filename, err)
	}

	var list []data.ItemNode
	nmap := make(map[uuid.UUID]data.ItemNode)
	for _, pitem := range items.ItemDefinitions {
		n := data.NewItemNode(
			getUUIDFromProtoGuid(pitem.ID),
			pitem.Item.Name,
			getUUIDFromProtoGuid(pitem.Item.TemplateID),
			getUUIDFromProtoGuid(pitem.Item.ParentID),
			getUUIDFromProtoGuid(pitem.Item.MasterID),
		)

		nmap[n.GetId()] = n
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
					fldItem, ok := nmap[fieldID]
					var name string
					if ok {
						name = fldItem.GetName()
					}

					fv := data.NewFieldValue(fieldID, id, name, f.Value, data.Language(ld.Language), int64(v.Version), stat.ModTime(), stat.ModTime(), data.VersionedFields)
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
			fldItem, ok := nmap[fieldID]
			var name string
			if ok {
				name = fldItem.GetName()
			}
			fv := data.NewFieldValue(fieldID, id, name, fld.Value, data.English, 1, stat.ModTime(), stat.ModTime(), data.SharedFields)
			m[id] = append(m[id], fv)
		}
	}

	for i := 0; i < len(list); i++ {
		flds := m[list[i].GetId()]

		for _, f := range flds {
			list[i].AddFieldValue(f)
		}
	}

	for _, item := range list {
		parent, ok := nmap[item.GetParentId()]
		if ok {
			item.SetParent(parent)
			item.GetParent().AddChild(item)
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
