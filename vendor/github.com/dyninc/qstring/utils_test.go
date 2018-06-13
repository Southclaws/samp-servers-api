package qstring

import (
	"reflect"
	"testing"
	"time"
)

func TestIsEmptyValue(t *testing.T) {
	var ts *TestStruct
	ts = nil
	testIO := []struct {
		inp      reflect.Value
		expected bool
	}{
		{inp: reflect.ValueOf([]int{0}), expected: false},
		{inp: reflect.ValueOf([]int{}), expected: true},
		{inp: reflect.ValueOf(map[string]int{"a": 0}), expected: false},
		{inp: reflect.ValueOf(map[string]int{}), expected: true},
		{inp: reflect.ValueOf(false), expected: true},
		{inp: reflect.ValueOf(true), expected: false},
		{inp: reflect.ValueOf(5), expected: false},
		{inp: reflect.ValueOf(0), expected: true},
		{inp: reflect.ValueOf(uint(5)), expected: false},
		{inp: reflect.ValueOf(uint(0)), expected: true},
		{inp: reflect.ValueOf(float32(5)), expected: false},
		{inp: reflect.ValueOf(float32(0)), expected: true},
		{inp: reflect.ValueOf(&TestStruct{}), expected: false},
		{inp: reflect.ValueOf(ts), expected: true},
		{inp: reflect.ValueOf(nil), expected: false},
		{inp: reflect.ValueOf(time.Time{}), expected: true},
		{inp: reflect.ValueOf(time.Now()), expected: false},
		{inp: reflect.ValueOf(*NewComparativeTime()), expected: true},
		{inp: reflect.ValueOf(ComparativeTime{Operator: "=", Time: time.Now()}),
			expected: false},
	}

	var result bool
	for _, test := range testIO {
		result = isEmptyValue(test.inp)
		if result != test.expected {
			t.Errorf("Expected %t for input %s, got %t", test.expected, test.inp, result)
		}
	}
}

func TestTagParsing(t *testing.T) {
	testio := []struct {
		inp    string
		output string
		omit   bool
	}{
		{inp: "name,omitempty", output: "name", omit: true},
		{inp: "name", output: "name", omit: false},
		{inp: "name,", output: "name", omit: false},
		{inp: "name", output: "name", omit: false},
		{inp: "", output: "", omit: false},
		{inp: ",omitempty", output: "", omit: true},
		{inp: "-", output: "-", omit: false},
	}

	var name string
	var omit bool
	for _, test := range testio {
		name, omit = parseTag(test.inp)
		if name != test.output {
			t.Errorf("Expected tag name to be %q, got %q instead", test.output, name)
		}

		if omit != test.omit {
			t.Errorf("Expected omitempty to be %t, got %t instead", test.omit, omit)
		}
	}
}

func TestOmitEmpty(t *testing.T) {
	vis := 5
	testio := []struct {
		Visible     int
		Conditional int `qstring:"conditional,omitempty"`
		omit        bool
	}{
		{Visible: vis, Conditional: 5, omit: false},
		{Visible: vis, Conditional: 0, omit: true},
	}

	for _, test := range testio {
		values, _ := Marshal(&test)

		_, found := values["conditional"]
		if found && test.omit {
			t.Errorf("%d was unexpectedly marshaled", test.Conditional)
		} else if !found && !test.omit {
			t.Errorf("%d was not marshaled, but should have been", test.Conditional)
		}
	}
}
