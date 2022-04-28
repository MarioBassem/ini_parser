package ini_parser

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	t.Run("read a file", func(t *testing.T) {
		parser := Parser{}
		want := Parser{
			"default": {
				"key1": "val1",
				"key2": "val2",
			},
		}

		parser.readFile("test.ini")

		if !reflect.DeepEqual(want, parser) {
			t.Errorf("maps don't match")
		}
	})

	t.Run("write to a file", func(t *testing.T) {
		parser := Parser{}
		parser["section1"] = map[string]string{}
		parser["section1"]["key1"] = "val1"
		parser.writeToFile("written.ini")

		got := Parser{}
		got.readFile("written.ini")
		want := parser

		if !reflect.DeepEqual(got, want) {
			t.Errorf("maps don't match")
		}
	})

	t.Run("validate sections", func(t *testing.T) {
		sections := []string{
			"[]", "[s", "s]", "s", "", "[s s]", "[s;s]", "[s]",
		}

		got := []bool{}
		for _, s := range sections {
			got = append(got, validateSection(s))
		}
		want := []bool{
			false, false, false, false, false, false, false, true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, but wanted %v", got, want)
		}
	})

	t.Run("validate keys value pair line", func(t *testing.T) {
		keyval := []string{
			"", "k", "k v", "k=k = v", "[k=v]", "k = v",
		}

		got := []bool{}
		for _, s := range keyval {
			got = append(got, validateKeyValPair(s))
		}

		want := []bool{
			false, false, false, false, false, true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, but wanted %v", got, want)
		}
	})

	t.Run("validate keys", func(t *testing.T) {
		keys := []string{
			"", "k k", "k;k", "k",
		}

		got := []bool{}
		for _, s := range keys {
			got = append(got, validateKey(s))
		}

		want := []bool{
			false, false, false, true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, but wanted %v", got, want)
		}
	})

	t.Run("validate values", func(t *testing.T) {
		vals := []string{
			"", "v v", "v;v", "v",
		}

		got := []bool{}
		for _, s := range vals {
			got = append(got, validateValue(s))
		}

		want := []bool{
			false, false, false, true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, but wanted %v", got, want)
		}
	})

}

/*
	test validity of sections and keys
*/
