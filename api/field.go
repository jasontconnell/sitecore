package api

import (
	// "fmt"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"regexp"
	"strings"
)

func getRepFunc(itemMap data.ItemMap, repMap map[uuid.UUID]uuid.UUID, notFound string) func(string) string {
	uidMap := make(map[uuid.UUID]uuid.UUID)
	for _, item := range itemMap {
		if repMap == nil {
			uidMap[item.GetId()] = emptyUuid
		} else if id, ok := repMap[item.GetId()]; ok {
			uidMap[item.GetId()] = id
		} else {
			uidMap[item.GetId()] = emptyUuid
		}
	}
	return func(s string) string {
		orig := s
		curlies := s[0] == '{'
		nohyphens := strings.IndexAny(s, "-") == -1
		upper := strings.IndexAny(s, "ABCDEF") != -1
		if curlies {
			s = string(s[1 : len(s)-1])
		}
		if len(s) == 32 {
			s = strings.ToLower(string(s[:8]) + "-" + string(s[8:12]) + "-" + string(s[12:16]) + "-" + string(s[16:20]) + "-" + string(s[20:]))
		}
		u2 := uuid.Must(uuid.Parse(s))
		var repId uuid.UUID
		outstring := orig
		if id, ok := uidMap[u2]; ok {
			repId = id
		} else {
			return orig
		}

		if repId == emptyUuid { // no replacement, don't do anything
			return orig
		}

		outstring = repId.String()
		if nohyphens {
			outstring = strings.Replace(outstring, "-", "", -1)
		}
		if upper {
			outstring = strings.ToUpper(outstring)
		}
		if curlies {
			outstring = "{" + outstring + "}"
		}
		return outstring
	}
}

var uuidReg *regexp.Regexp = regexp.MustCompile(`(?i){?([a-f0-9]{8})-?([a-f0-9]{4}-?){3}([a-f0-9]{12})}?`)

func FilterFields(itemMap, filterBy data.ItemMap) {
	f := getRepFunc(filterBy, nil, "")
	for _, item := range itemMap {
		for _, fv := range item.GetFieldValues() {
			newval := uuidReg.ReplaceAllStringFunc(fv.GetValue(), f)
			fv.SetValue(newval)
		}
	}
}
