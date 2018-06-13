package qstring

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestMarshallString(t *testing.T) {
	ts := TestStruct{
		Name:     "SomeName",
		Do:       true,
		Page:     1,
		ID:       12,
		Small:    13,
		Med:      14,
		Big:      15,
		UPage:    2,
		UID:      16,
		USmall:   17,
		UMed:     17,
		UBig:     17,
		Float32:  6000,
		Float64:  7000,
		Fields:   []string{"foo", "bar"},
		DoFields: []bool{true, false},
		Counts:   []int{1, 2},
		IDs:      []int8{3, 4},
		Smalls:   []int16{6, 7},
		Meds:     []int32{9, 10},
		Bigs:     []int64{12, 13},
		UPages:   []uint{2, 3},
		UIDs:     []uint8{5, 6},
		USmalls:  []uint16{8, 9},
		UMeds:    []uint32{9, 10},
		UBigs:    []uint64{12, 13},
		Float32s: []float32{6000, 6001},
		Float64s: []float64{7000, 7001},
	}

	expected := []string{"name=SomeName", "do=true", "page=1", "id=12", "small=13",
		"med=14", "big=15", "upage=2", "uid=16", "usmall=17", "umed=17", "ubig=17",
		"float32=6000", "float64=7000", "fields=foo", "fields=bar", "dofields=true",
		"dofields=false", "counts=1", "counts=2", "ids=3", "ids=4", "smalls=6",
		"smalls=7", "meds=9", "meds=10", "bigs=12", "bigs=13", "upages=2",
		"upages=3", "uids=5", "uids=6", "usmalls=8", "usmalls=9", "umeds=9",
		"umeds=10", "ubigs=12", "ubigs=13", "float32s=6000", "float32s=6001",
		"float64s=7000", "float64s=7001"}
	query, err := MarshalString(&ts)
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, param := range expected {
		if !strings.Contains(query, param) {
			t.Errorf("Expected %s to contain %s", query, param)
		}
	}
}

func TestMarshallValues(t *testing.T) {
	ts := TestStruct{
		Name:     "SomeName",
		Do:       true,
		Page:     1,
		ID:       12,
		Small:    13,
		Med:      14,
		Big:      15,
		UPage:    2,
		UID:      16,
		USmall:   17,
		UMed:     17,
		UBig:     17,
		Float32:  6000,
		Float64:  7000,
		Fields:   []string{"foo", "bar"},
		DoFields: []bool{true, false},
		Counts:   []int{1, 2},
		IDs:      []int8{3, 4},
		Smalls:   []int16{6, 7},
		Meds:     []int32{9, 10},
		Bigs:     []int64{12, 13},
		UPages:   []uint{2, 3},
		UIDs:     []uint8{5, 6},
		USmalls:  []uint16{8, 9},
		UMeds:    []uint32{9, 10},
		UBigs:    []uint64{12, 13},
		Float32s: []float32{6000, 6001},
		Float64s: []float64{7000, 7001},
	}

	expected := url.Values{
		"name":     []string{"SomeName"},
		"do":       []string{"true"},
		"page":     []string{"1"},
		"id":       []string{"12"},
		"small":    []string{"13"},
		"med":      []string{"14"},
		"big":      []string{"15"},
		"upage":    []string{"2"},
		"uid":      []string{"16"},
		"usmall":   []string{"17"},
		"umed":     []string{"18"},
		"ubig":     []string{"19"},
		"float32":  []string{"6000"},
		"float64":  []string{"7000"},
		"fields":   []string{"foo", "bar"},
		"dofields": []string{"true", "false"},
		"counts":   []string{"1", "2"},
		"ids":      []string{"3", "4", "5"},
		"smalls":   []string{"6", "7", "8"},
		"meds":     []string{"9", "10", "11"},
		"bigs":     []string{"12", "13", "14"},
		"upages":   []string{"2", "3", "4"},
		"uids":     []string{"5, 6, 7"},
		"usmalls":  []string{"8", "9", "10"},
		"umeds":    []string{"9", "10", "11"},
		"ubigs":    []string{"12", "13", "14"},
		"float32s": []string{"6000", "6001", "6002"},
		"float64s": []string{"7000", "7001", "7002"},
	}
	values, err := Marshal(&ts)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(values) != len(expected) {
		t.Errorf("Expected %d fields, got %d. Hidden is %q",
			len(expected), len(values), values["hidden"])
	}
}

func TestInvalidMarshalString(t *testing.T) {
	var err error
	var ts *TestStruct
	testio := []struct {
		inp       interface{}
		errString string
	}{
		{inp: nil, errString: "qstring: MarshalString(nil)"},
		{inp: TestStruct{}, errString: "qstring: MarshalString(non-pointer qstring.TestStruct)"},
		{inp: ts, errString: "qstring: MarshalString(nil *qstring.TestStruct)"},
	}

	for _, test := range testio {
		_, err = MarshalString(test.inp)
		if err == nil {
			t.Errorf("Expected invalid type error, got success instead")
		}

		if err.Error() != test.errString {
			t.Errorf("Got %q error, expected %q", err.Error(), test.errString)
		}
	}
}

func TestMarshalTime(t *testing.T) {
	type Query struct {
		Created     time.Time
		LastUpdated time.Time
	}

	createdTS := "2006-01-02T15:04:05Z"
	createdTime, _ := time.Parse(time.RFC3339, createdTS)
	updatedTS := "2016-01-02T15:04:05-07:00"
	updatedTime, _ := time.Parse(time.RFC3339, updatedTS)

	q := &Query{Created: createdTime, LastUpdated: updatedTime}
	result, err := MarshalString(q)
	if err != nil {
		t.Fatalf("Unable to marshal timestamp: %s", err.Error())
	}

	var unescaped string
	unescaped, err = url.QueryUnescape(result)
	if err != nil {
		t.Fatalf("Unable to unescape query string %q: %q", result, err.Error())
	}

	expected := []string{"created=2006-01-02T15:04:05Z",
		"lastupdated=2016-01-02T15:04:05-07:00"}
	for _, ts := range expected {
		if !strings.Contains(unescaped, ts) {
			t.Errorf("Expected query string %s to contain %s", unescaped, ts)
		}
	}
}

func TestMarshalNested(t *testing.T) {
	type Paging struct {
		Page  int
		Limit int
	}

	type Params struct {
		Paging Paging
		Name   string
	}

	params := &Params{Name: "SomeName",
		Paging: Paging{Page: 1, Limit: 50},
	}

	result, err := MarshalString(params)
	if err != nil {
		t.Fatalf("Unable to marshal nested struct: %s", err.Error())
	}

	var unescaped string
	unescaped, err = url.QueryUnescape(result)
	if err != nil {
		t.Fatalf("Unable to unescape query string %q: %q", result, err.Error())
	}

	// ensure the nested struct isn't iteself included in the query string
	if strings.Contains(unescaped, "paging=") {
		t.Errorf("Nested struct was included in %q", unescaped)
	}

	// ensure fields we expect to be present are
	expected := []string{"name=SomeName", "page=1", "limit=50"}
	for _, q := range expected {
		if !strings.Contains(unescaped, q) {
			t.Errorf("Expected query string %s to contain %s", unescaped, q)
		}
	}
}

func TestMarshalNestedPtrs(t *testing.T) {
	type Paging struct {
		Page  int
		Limit int
	}

	type Params struct {
		Paging *Paging
		Name   string
	}

	params := &Params{Name: "SomeName",
		Paging: &Paging{Page: 1, Limit: 50},
	}

	result, err := MarshalString(params)
	if err != nil {
		t.Fatalf("Unable to marshal nested struct: %s", err.Error())
	}

	var unescaped string
	unescaped, err = url.QueryUnescape(result)
	if err != nil {
		t.Fatalf("Unable to unescape query string %q: %q", result, err.Error())
	}

	// ensure the nested struct isn't iteself included in the query string
	if strings.Contains(unescaped, "paging=") {
		t.Errorf("Nested struct was included in %q", unescaped)
	}

	// ensure fields we expect to be present are
	expected := []string{"name=SomeName", "page=1", "limit=50"}
	for _, q := range expected {
		if !strings.Contains(unescaped, q) {
			t.Errorf("Expected query string %s to contain %s", unescaped, q)
		}
	}
}

func (u *MarshalInterfaceTest) MarshalQuery() (url.Values, error) {
	return url.Values{"names": u.Names}, nil
}

func TestMarshaller(t *testing.T) {
	s := &MarshalInterfaceTest{Names: []string{"foo", "bar"}}

	testIO := []struct {
		inp  *MarshalInterfaceTest
		vals url.Values
		err  error
	}{
		{s, url.Values{"names": s.Names}, nil},
	}

	for _, test := range testIO {
		v, err := Marshal(test.inp)
		if err != test.err {
			t.Errorf("Expected Marshaller to return %s, but got %s instead", test.err, err)
		}

		if len(v) != len(test.vals) {
			t.Errorf("Expected %q, got %q instead", test.vals, v)
		}
	}
}
