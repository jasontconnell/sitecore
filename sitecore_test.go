package sitecore

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jasontconnell/sitecore/api"
)

var connstr string = os.Getenv("SitecoreAPITestConnectionString")

func TestBlob(t *testing.T) {
	id := api.MustParseUUID("{8A28BA6E-5A13-47E8-8D25-983D3BD0E761}")
	b, err := api.LoadBlob(connstr, id)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	fmt.Println(b.GetData(), b.GetId())
}

func TestLoadItemMap(t *testing.T) {
	t.Log(connstr)
	start := time.Now()
	items, err := api.LoadItems(connstr)
	t.Log("Loaded items", time.Since(start))

	if err != nil {
		t.Fatal(err)
	}

	start = time.Now()
	root, itemMap := api.LoadItemMap(items)
	t.Log("Loaded item map", time.Since(start))

	if root == nil {
		t.Fatal("couldn't find root")
	}

	if len(itemMap) == 0 {
		t.Fatal("no items")
	}

	t.Log("Root path", root.GetPath())

	filtered := api.FilterItemMap(itemMap, []string{"/sitecore/templates"})
	t.Log("Filtered item map", len(filtered))

	t.Log(root.GetId(), root.GetPath(), len(itemMap), time.Since(start))

	start = time.Now()
	fields, ferr := api.LoadFields(connstr)
	if ferr != nil {
		t.Fatal(ferr)
	}

	t.Log("Fields loaded", time.Since(start))

	if len(fields) == 0 {
		t.Fatal("no fields received")
	}

	start = time.Now()
	npfields, nperr := api.LoadFieldsParallel(connstr, 12)
	if nperr != nil {
		t.Fatal(nperr)
	}

	t.Log("Loaded fields parallel", time.Since(start))

	if len(npfields) == 0 {
		t.Fatal("non parallel fields is empty")
	}

	api.AssignFieldValues(itemMap, npfields)

	if len(fields) != len(npfields) {
		t.Log("len of non parallel fields", len(npfields))
		t.Log("len of parallel fields", len(fields))

		t.Fatal("len received from non-parallel is not equal to parallel version")
	}

	testuid := api.MustParseUUID("9541e67d-ce8c-4225-803d-33f7f29f09ef")

	start = time.Now()
	fieldMap := api.LoadFieldMap(fields)

	t.Log("Loaded field map", time.Since(start))
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

func TestRenderings(t *testing.T) {
	testId := api.MustParseUUID("{768D4C5A-E8D8-4D63-BD2F-6C7D4F2BD080}")
	items, err := api.LoadItems(connstr)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	tlist, err := api.LoadTemplates(connstr)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	_, itemMap := api.LoadItemMap(items)

	npfields, nperr := api.LoadFieldsParallel(connstr, 12)
	if nperr != nil {
		t.Fatal(nperr)
	}

	api.AssignFieldValues(itemMap, npfields)

	tm := api.GetTemplateMap(tlist)

	t.Log("Mapping Layouts")
	mperr := api.MapAllLayouts(itemMap, tm, true)
	if mperr != nil {
		t.Fatal(mperr)
	}

	testItem, ok := itemMap[testId]
	if !ok {
		t.Log("Item not found")
		t.Fail()
	}

	t.Log("getting renderings for ", testItem.GetId())
	for _, r := range testItem.GetRenderings() {
		t.Log(r)
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
