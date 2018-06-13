package qstring

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

type TestStruct struct {
	Name string `qstring:"name"`
	Do   bool

	// int fields
	Page  int `qstring:"page"`
	ID    int8
	Small int16
	Med   int32
	Big   int64

	// uint fields
	UPage  uint
	UID    uint8
	USmall uint16
	UMed   uint32
	UBig   uint64

	// Floats
	Float32 float32
	Float64 float64

	// slice fields
	Fields   []string `qstring:"fields"`
	DoFields []bool   `qstring:"dofields"`
	Counts   []int
	IDs      []int8
	Smalls   []int16
	Meds     []int32
	Bigs     []int64

	// uint fields
	UPages  []uint
	UIDs    []uint8
	USmalls []uint16
	UMeds   []uint32
	UBigs   []uint64

	// Floats
	Float32s []float32
	Float64s []float64
	hidden   int
	Hidden   int `qstring:"-"`
}

func TestUnmarshall(t *testing.T) {
	var ts TestStruct
	query := url.Values{
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
		"uids":     []string{"5", "6", "7"},
		"usmalls":  []string{"8", "9", "10"},
		"umeds":    []string{"9", "10", "11"},
		"ubigs":    []string{"12", "13", "14"},
		"float32s": []string{"6000", "6001", "6002"},
		"float64s": []string{"7000", "7001", "7002"},
	}

	err := Unmarshal(query, &ts)
	if err != nil {
		t.Fatal(err.Error())
	}

	if ts.Page != 1 {
		t.Errorf("Expected page to be 1, got %d", ts.Page)
	}

	if len(ts.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(ts.Fields))
	}
}

func TestUnmarshalNested(t *testing.T) {
	type Paging struct {
		Page  int
		Limit int
	}

	type Params struct {
		Paging Paging
		Name   string
	}

	query := url.Values{
		"name":  []string{"SomeName"},
		"page":  []string{"1"},
		"limit": []string{"50"},
	}

	params := &Params{}

	err := Unmarshal(query, params)
	if err != nil {
		t.Fatal(err.Error())
	}

	if params.Paging.Page != 1 {
		t.Errorf("Nested Struct Failed to Unmarshal. Expected 1, got %d", params.Paging.Page)
	}
}

func TestUnmarshalTime(t *testing.T) {
	type Query struct {
		Created     time.Time
		LastUpdated time.Time
	}

	createdTS := "2006-01-02T15:04:05Z"
	updatedTS := "2016-01-02T15:04:05-07:00"

	query := url.Values{
		"created":     []string{createdTS},
		"lastupdated": []string{updatedTS},
	}

	params := &Query{}
	err := Unmarshal(query, params)
	if err != nil {
		t.Fatal(err.Error())
	}

	if params.Created.Format(time.RFC3339) != createdTS {
		t.Errorf("Expected created ts of %s, got %s instead.", createdTS, params.Created.Format(time.RFC3339))
	}

	if params.LastUpdated.Format(time.RFC3339) != updatedTS {
		t.Errorf("Expected update ts of %s, got %s instead.", updatedTS, params.LastUpdated.Format(time.RFC3339))
	}
}

func TestUnmarshalInvalidTypes(t *testing.T) {
	var err error
	var ts *TestStruct
	testio := []struct {
		inp       interface{}
		errString string
	}{
		{inp: nil, errString: "qstring: Unmarshal(nil)"},
		{inp: TestStruct{}, errString: "qstring: Unmarshal(non-pointer qstring.TestStruct)"},
		{inp: ts, errString: "qstring: Unmarshal(nil *qstring.TestStruct)"},
	}

	for _, test := range testio {
		err = Unmarshal(url.Values{}, test.inp)
		if err == nil {
			t.Errorf("Expected invalid type error, got success instead")
		}

		if err.Error() != test.errString {
			t.Errorf("Got %q error, expected %q", err.Error(), test.errString)
		}
	}
}

var errNoNames = errors.New("No Names Provided")

type MarshalInterfaceTest struct {
	Names []string
}

func (u *MarshalInterfaceTest) UnmarshalQuery(v url.Values) error {
	var ok bool
	if u.Names, ok = v["names"]; ok {
		return nil
	}
	return errNoNames
}

func TestUnmarshaller(t *testing.T) {
	testIO := []struct {
		inp      url.Values
		expected interface{}
	}{
		{url.Values{"names": []string{"foo", "bar"}}, nil},
		{make(url.Values), errNoNames},
	}

	s := &MarshalInterfaceTest{Names: []string{}}
	for _, test := range testIO {
		err := Unmarshal(test.inp, s)
		if err != test.expected {
			t.Errorf("Expected Unmarshaller to return %s, but got %s instead", test.expected, err)
		}
	}
}
