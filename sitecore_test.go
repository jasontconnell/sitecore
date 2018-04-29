package sitecore

import (
	"github.com/google/uuid"
	"sitecore/api"
	"testing"
)

var connstr string = `user id=sa;password=S4M3amg;server=localhost\MSSQL_2014;Database=JGWentworth_Dev_Master`

func TestLoadItemMap(t *testing.T) {
	items, err := api.LoadItems(connstr)

	if err != nil {
		t.Fatal(err)
	}

	root, itemMap := api.LoadItemMap(items)

	if root == nil {
		t.Fatal("couldn't find root")
	}

	if len(itemMap) == 0 {
		t.Fatal("no items")
	}

	t.Log(root.GetId(), root.GetPath(), len(itemMap))

	fields, ferr := api.LoadFields(connstr)
	if ferr != nil {
		t.Fatal(ferr)
	}

	if len(fields) == 0 {
		t.Fatal("no fields received")
	}

	npfields, nperr := api.LoadFields(connstr)
	if nperr != nil {
		t.Fatal(nperr)
	}

	if len(npfields) == 0 {
		t.Fatal("non parallel fields is empty")
	}

	if len(fields) != len(npfields) {
		t.Log("len of non parallel fields", len(npfields))
		t.Log("len of parallel fields", len(fields))

		t.Fatal("len received from non-parallel is not equal to parallel version")
	}

	testuid, _ := uuid.Parse("9541e67d-ce8c-4225-803d-33f7f29f09ef")

	fieldMap := api.LoadFieldMap(fields)

	fl, ok := fieldMap[testuid]
	if !ok {
		t.Fatal("expected item not found")
	}

	if len(fl) == 0 {
		t.Fatal("expected item found but no field values")
	}

	if len(fieldMap) != len(itemMap) {
		t.Fatal("not the same amount of items in field map vs item map", len(fieldMap), len(itemMap))
	}
}

func TestLoadTemplates(t *testing.T) {
	tmps, err := api.LoadTemplates(connstr)

	if err != nil {
		t.Fatal("Error loading templates", err)
	}

	root, tmap := api.LoadItemMap(tmps)

	api.LoadTemplateData(tmap)

	t.Log(root.GetId(), len(tmap))

	list := api.GetTemplates(tmap)
	if len(list) == 0 {
		t.Fatal("No templates")
	}
}

func BenchmarkFieldLoad(b *testing.B) {
	fields, err := api.LoadFields(connstr)
	b.ReportAllocs()

	if err != nil {
		b.Fatal(err)
	}

	if len(fields) == 0 {
		b.Fatal("No fields")
	}
}

func BenchmarkFieldLoadParallel(b *testing.B) {
	fields, err := api.LoadFieldsParallel(connstr, 20)

	b.ReportAllocs()

	if err != nil {
		b.Fatal(err)
	}

	if len(fields) == 0 {
		b.Fatal("No fields")
	}
}
