package mijnhost

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/libdns/libdns"
)

func TestProvider_Unmarshall(t *testing.T) {
	var provider *Provider
	var buf = `{
"api_key": "testkey",
"base_uri": "http://127.0.0.1:8080",
"debug": true
}`

	if err := json.Unmarshal([]byte(buf), &provider); err != nil {
		t.Fatal(err)
	}

	if provider.GetApiKey() != "testkey" {
		t.Fatalf("api key = %s; want %s", provider.GetApiKey(), "testkey")
	}

	if provider.GetBaseUri().String() != "http://127.0.0.1:8080" {
		t.Fatalf("base uri = %s; want %s", provider.GetBaseUri().String(), "http://127.0.0.1:8080")
	}

	if provider.GetDebug() != os.Stdout {
		t.Fatal("expected debug to be os.Stdout")
	}
}

func TestProvider(t *testing.T) {

	provider := &Provider{ApiKey: os.Getenv("API_KEY")}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		provider.Debug = true
	}

	zones, err := provider.ListZones(context.Background())

	if err != nil {
		t.Fatalf("ListZones failed: %v", err)
	}

	for _, zone := range zones {

		t.Logf("zone: %s", zone.Name)

		var name = fmt.Sprintf("__test_%s", randomString(8))
		var records = []libdns.Record{
			libdns.TXT{
				Name: name,
				Text: randomString(32),
			},
			libdns.TXT{
				Name: name,
				Text: randomString(32),
			},
		}

		var original = validateRemote(records, -1, zone.Name, provider, t)

		t.Logf("found %d records", len(original))
		output(original, "data:", t)

		setRecords(records, zone.Name, provider, t)
		validateRemote(records, 0, zone.Name, provider, t)
		removeRecords(records, name, zone.Name, provider, t)
		validateRemote(records, 1, zone.Name, provider, t)
		validateRemote(original, 2, zone.Name, provider, t)
	}
}

func output(list []libdns.Record, prefix string, t *testing.T) {
	if out, e := json.MarshalIndent(list, "  ", " "); e == nil {
		t.Logf("%s\n%s", prefix, string(out))
	}
}

func setRecords(list []libdns.Record, zone string, provider *Provider, t *testing.T) {
	output(list, "adding records:", t)

	set, err := provider.SetRecords(context.Background(), zone, list)

	if err != nil {
		t.Fatalf("SetRecords failed: %v", err)
	}

	if len(list) != len(set) || false == validateRecords(set, list) {
		t.Fatalf("SetRecords returned %#v records, expected %#v", set, list)
	}

	t.Log("successfully set records")
}

func removeRecords(list []libdns.Record, name, zone string, provider *Provider, t *testing.T) {

	t.Log("removing created records")

	removed, err := provider.DeleteRecords(context.Background(), zone, []libdns.Record{libdns.TXT{Name: name}})

	if err != nil {
		t.Fatalf("DeleteRecords failed: %v", err)
	}

	if len(list) != len(removed) || false == validateRecords(removed, list) {
		t.Fatalf("DeleteRecords returned %#v records, expected %#v", removed, list)
	}

	t.Logf("successfully removed %d records", len(removed))
}

func validateRemote(list []libdns.Record, mode int, zone string, provider *Provider, t *testing.T) []libdns.Record {
	items, err := provider.GetRecords(context.Background(), zone)

	if err != nil {
		t.Fatalf("GetRecords failed: %v", err)
	}

	switch mode {
	case 0:
		t.Log("validating records by querying all available records")

		if false == validateRecords(list, items) {
			t.Fatalf("not all set records found in %+v looking for %+v", items, list)
		}

		t.Log("successfully found new records in remote set")
	case 1:
		t.Log("validating removed records by querying and validating against available records")

		for item := range RecordIterator(list).Iterate() {
			if validateRecords([]libdns.Record{item}, items) {
				t.Fatalf("not expected to find record %#v", item)
			}
		}

		t.Log("successfully removed records")
	case 2:
		t.Log("validating original records by querying against current available records")

		if len(list) != len(items) || false == validateRecords(list, items) {
			t.Fatalf("records not same a starting:\n%+v\n%+v", items, list)
		}

		t.Log("successfully matched original list against current available records")

	}

	return items

}

func validateRecords(x []libdns.Record, set []libdns.Record) bool {
	for a := range RecordIterator(x).Iterate() {
		for b := range RecordIterator(set).Iterate() {
			if strings.EqualFold(a.Name, b.Name) && a.Data == b.Data && a.Type == b.Type {
				return true
			}
		}
		return false
	}
	return false
}

func randomString(n int) string {
	var x = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var s = len(x)

	b := make([]byte, n)

	for i := range b {
		b[i] = x[rand.Intn(s)]
	}

	return string(b)
}
