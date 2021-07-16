package transformer

import (
	"encoding/json"
	"testing"
)

type Address struct {
	Place string `bson:"ort" keep:"update,omitempty"`
}
type TestStruct struct {
	Name    string       `bson:"name" keep:"update,create,omitempty"`
	Age     int          `bson:"ag" keep:"create,omitempty"`
	Childs  []TestStruct `keep:"update,omitempty"`
	Address Address      `bson:"addr" keep:"update,create,omitempty"`
}

func TestTransformer(t *testing.T) {

	testCases := []struct {
		desc string
		in   interface{}
		out  interface{}
	}{
		{
			desc: "Simple Struct",
			in:   TestStruct{Name: "Timo", Address: Address{""}},
		},
		{
			desc: "Simple Struct with Struct",
			in:   TestStruct{Name: "Timo", Age: 21, Address: Address{"Mannheim"}},
		},
		{
			desc: "Simple Struct with Array",
			in: TestStruct{Name: "Timo", Age: 21, Childs: []TestStruct{
				{Name: "Nico", Age: 23, Address: Address{"Lindau"}}, {Name: "Jeannine", Age: 22},
			}},
		},
	}
	t.Log("Running")
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			json, _ := json.MarshalIndent(Clean(tC.in, "update"), "", "  ")
			t.Log("\n\n" + string(json) + "\n")
		})
	}
	t.Log("Done")
}
