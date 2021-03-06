package basin

import (
	"regexp"
	"strings"
	"testing"
)

type StrMessage string

func (s StrMessage) Field(f string) (interface{}, bool) {
	return string(s), true
}

func (s StrMessage) Bytes() []byte {
	return []byte(s)
}

func TestPasses_NoFilter(t *testing.T) {
	msg := StrMessage("foo bar baz")

	filterOk := NewNoFilter()
	if !filterOk.Passes(msg) {
		t.Errorf("filter should always pass, but didn't for '%s'", msg)
	}
}

func TestPasses_ContainsFilter(t *testing.T) {
	msg := StrMessage("foo bar baz")

	filterOk := NewContainsFilter("", "bar")
	if !filterOk.Passes(msg) {
		t.Errorf("'bar' is contained in '%s'", msg)
	}

	filterBad := NewContainsFilter("", "qwijibo")
	if filterBad.Passes(msg) {
		t.Errorf("'qwijibo' is not contained in '%s'", msg)
	}
}

func TestPasses_RegexpFilter(t *testing.T) {
	msg := StrMessage("foo bar baz")

	filterOk := NewRegexpFilter("", regexp.MustCompile("[a-z]a[a-z]"))
	if !filterOk.Passes(msg) {
		t.Errorf("'[a-z]a[a-z]' should match in '%s'", msg)
	}

	filterBad := NewRegexpFilter("", regexp.MustCompile("[0-9]+"))
	if filterBad.Passes(msg) {
		t.Errorf("No numbers are contained within '%s'", msg)
	}
}

func TestPasses_FuncFilter(t *testing.T) {
	msg := StrMessage("foo bar baz")

	filterOk := NewFuncFilter(func(m Message) bool {
		value, _ := m.Field("")
		return strings.Contains(value.(string), "bar")
	})

	if !filterOk.Passes(msg) {
		t.Errorf("'bar' is contained in '%s'", msg)
	}

	filterBad := NewFuncFilter(func(m Message) bool {
		value, _ := m.Field("")
		return strings.Contains(value.(string), "qwijibo")
	})

	if filterBad.Passes(msg) {
		t.Errorf("'qwijibo' is not contained in '%s'", msg)
	}
}

func TestPasses_ComboFilter(t *testing.T) {
	msg := StrMessage("foo bar baz")
	filterOk := NewComboFilter(NewContainsFilter("", "bar"), NewContainsFilter("", "foo"))

	if !filterOk.Passes(msg) {
		t.Errorf("'foo' and 'bar' are contained within '%s'", msg)
	}

	filterBad := NewComboFilter(NewContainsFilter("", "qwijibo"), NewContainsFilter("", "monkey"))
	if filterBad.Passes(msg) {
		t.Errorf("neither 'qwijibo' nor 'monkey' are contained within '%s'", msg)
	}
}
