package verbalexpressions

import "testing"
import "strings"

func TestChaining(t *testing.T) {

	exp := "http://www.google.com"
	v := New().StartOfLine().
		Then("http").
		Maybe("s").
		Then("://").
		Maybe("www.").
		Word().
		Then(".").
		Word().
		Maybe("/").
		EndOfLine()
	if !v.Test(exp) {
		t.Errorf("%v regexp doesn't match %s", v.Regex(), exp)
	}
}

func TestRange(t *testing.T) {

	exp := "abcdef 123"

	v := New().Range("a", "z", 0, 9)
	if v.Regex().String() != "[a-z0-9]" {
		t.Errorf("%s is not [a-z0-9 ]", v.Regex())
	}
	if !v.Test(exp) {
		t.Errorf("%v regexp doesn't match %s", v.Regex(), exp)
	}
	exp = "ABCDEF"
	if v.Test(exp) {
		t.Errorf("%v regexp should not match %s", v.Regex(), exp)
	}

}

func TestPanicOnRangeOddParams(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Logf("panic accepted: %s", r)
		}
	}()

	New().Range("a", "z", 0, 9, 10)
	t.Errorf("Call must panic !")
}

func TestOneLine(t *testing.T) {
	s := "atlanta\narkansas\nalabama\narachnophobia"

	v := New().SearchOneLine(false).Find("a").EndOfLine().Regex()
	res := v.FindAllStringIndex(s, -1)
	if len(res) != 1 {
		t.Errorf("%v should be length 1, %d instead", res, len(res))
	}
	if len(res[0]) != 2 {
		t.Errorf("%v should be length 2, %d instead", res[0], len(res[0]))
	}

	v = New().SearchOneLine(true).Find("a").EndOfLine().Regex()
	res = v.FindAllStringIndex(s, -1)
	if len(res) != 3 {
		t.Errorf("%v should be length 1, %d instead", res, len(res))
	}
	for _, r := range res {
		if len(r) != 2 {
			t.Errorf("%v should be length 2, %d instead", r, len(r))
		}
	}
}

func TestExcepting(t *testing.T) {

	s := "This is a simple test"
	v := New().AnythingBut("im").Regex().FindAllString(s, -1)
	for _, st := range v {
		if strings.Contains(st, "i") {
			t.Errorf("%s should not find \"i\"", st)
		}
		if strings.Contains(st, "m") {
			t.Errorf("%s should not find \"m\"", st)
		}
	}
}

func TestAny(t *testing.T) {

	s := "foo1 foo5 foobar"
	v := New().Find("foo").Any("1234567890")
	res := v.Regex().FindAllString(s, -1)
	if len(res) != 2 {
		t.Errorf("len(res) : %d isn't 2", len(res))
	}

}

func TestReplace(t *testing.T) {

	s := "foomode barmode themodebaz"
	expect := "foochanged barchanged thechangedbaz"

	v := New().Find("mode")
	res := v.Replace(s, "changed")

	if res != expect {
		t.Errorf("Replacement hasn't worked as expected %s != %s", res, expect)
	}

}
