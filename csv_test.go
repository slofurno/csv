package csv

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

type unquotedTest struct {
	name    string
	src     string
	records [][]string
}

func TestUnquoted(t *testing.T) {

	tests := []unquotedTest{
		{"ok", "a,b,c\nd,e,f\n", [][]string{
			{"a", "b", "c"},
			{"d", "e", "f"}},
		},
		{"unquoted quotes", "a,\"b,c\n\"d,e,f\"\n", [][]string{
			{"a", "\"b", "c"},
			{"\"d", "e", "f\""}},
		},
		{"missing final newline", "a,b\nc,d", [][]string{
			{"a", "b"},
			{"c", "d"}},
		},
		{"carriage returns", "a,b\r\nc,d\r\n", [][]string{
			{"a", "b"},
			{"c", "d"}},
		},
	}

	for _, test := range tests {
		reader := NewUnquoted(strings.NewReader(test.src))
		j := 0

		for {
			fields, err := reader.Read()
			if err == io.EOF {
				if j != len(test.records) {
					t.Fatalf("test %s: expected records: %d, actual: %d\n",
						test.name, len(test.records), j)
				}

				break
			}
			if err != nil {
				t.Fatalf("test %s error with: %s\n", test.name, err)
			}

			fmt.Println(fields)

			if len(fields) != len(test.records[j]) {
				t.Fatalf("test %s: expected fields: %d, actual: %d\n",
					test.name, len(test.records[j]), len(fields))
			}
			for i := range fields {
				if fields[i] != test.records[j][i] {
					t.Fatalf("test %s: expected field: %s, actual: %s\n",
						test.name, test.records[j][i], fields[i])
				}
			}
			j++
		}
	}

}
