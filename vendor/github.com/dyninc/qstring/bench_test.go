package qstring

import (
	"net/url"
	"testing"
)

// Straight benchmark literal.
func BenchmarkUnmarshall(b *testing.B) {
	query := url.Values{
		"limit":  []string{"10"},
		"page":   []string{"1"},
		"fields": []string{"a", "b", "c"},
	}
	type QueryStruct struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &QueryStruct{}
		err := Unmarshal(query, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Parallel benchmark literal.
func BenchmarkRawPLiteral(b *testing.B) {
	query := url.Values{
		"limit":  []string{"10"},
		"page":   []string{"1"},
		"fields": []string{"a", "b", "c"},
	}
	type QueryStruct struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := &QueryStruct{}
			err := Unmarshal(query, data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
