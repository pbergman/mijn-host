package mijn_host

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/libdns/libdns"
)

func TestProvider(t *testing.T) {

	provider := NewProvider()
	provider.SetApiKey(os.Getenv("API_KEY"))

	if _, ok := os.LookupEnv("DEBUG"); ok {
		provider.SetDebug(os.Stdout)
	}

	zones, err := provider.ListZones(context.Background())

	if err != nil {
		t.Fatalf("ListZones failed: %v", err)
	}

	for _, zone := range zones {

		var name = fmt.Sprintf("__test_%s", randomString(8))
		var value1 = randomString(32)
		var value2 = randomString(32)
		var records = []libdns.Record{
			libdns.TXT{
				Name: name,
				Text: value1,
			},
			libdns.TXT{
				Name: name,
				Text: value2,
			},
		}

		t.Logf("adding 2 TXT records records (%s) with values \"%s\" && \"%s\"", name, value1, value2)

		set, err := provider.SetRecords(context.Background(), zone.Name, records)

		if err != nil {
			t.Fatalf("SetRecords failed: %v", err)
		}

		if len(records) != len(set) || false == validateRecords(set, records) {
			t.Fatalf("SetRecords returned %#v records, expected %#v", set, records)
		}

		t.Log("successfully set records")

		t.Log("validating records by querying all available records")

		curr, err := provider.GetRecords(context.Background(), zone.Name)

		if err != nil {
			t.Fatalf("GetRecords failed: %v", err)
		}

		if false == validateRecords(records, curr) {
			t.Fatalf("GetRecords returned %#v records, expected %#v", set, records)
		}

		t.Log("successfully found new records in remote set")

		t.Log("removing created records")

		removed, err := provider.DeleteRecords(context.Background(), zone.Name, []libdns.Record{libdns.TXT{Name: name}})

		if err != nil {
			t.Fatalf("DeleteRecords failed: %v", err)
		}

		if len(records) != len(removed) || false == validateRecords(removed, records) {
			t.Fatalf("DeleteRecords returned %#v records, expected %#v", removed, records)
		}

		t.Logf("successfully removed %d records", len(removed))

		t.Log("validating removed records by querying and validating against available records")

		x, err := provider.GetRecords(context.Background(), zone.Name)

		if err != nil {
			t.Fatalf("GetRecords failed: %v", err)
		}

		for item := range RecordIterator(records).Iterate() {
			if validateRecords([]libdns.Record{item}, x) {
				t.Fatalf("not expected to find record %#v", item)
			}
		}

		t.Log("successfully validating removed records")
	}
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
