// Copyright 2013 Patrice FERLET
// Ue of this source code is governed by MIT-style
// license that can be found in the LICENSE file
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
	if v.Regex().String() != "(?m)[a-z0-9]" {
		t.Errorf("%s is not (?m)[a-z0-9]", v.Regex())
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
		// if no panic... the test fails
		if r := recover(); r == nil {
			t.Errorf("Call must panic !")
		}
	}()

	New().Range("a", "z", 0, 9, 10)
}

func TestOneLine(t *testing.T) {
	s := "atlanta\narkansas\nalabama\narachnophobia"

	v := New().SearchOneLine(true).Find("a").EndOfLine().Regex()
	res := v.FindAllStringIndex(s, -1)
	if len(res) != 1 {
		t.Errorf("%v should be length 1, %d instead", res, len(res))
	}
	if len(res[0]) != 2 {
		t.Errorf("%v should be length 2, %d instead", res[0], len(res[0]))
	}

	v = New().SearchOneLine(false).Find("a").EndOfLine().Regex()
	res = v.FindAllStringIndex(s, -1)
	if len(res) != 3 {
		t.Errorf("%v should be length 3, %d instead", res, len(res))
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

func TestCaptures(t *testing.T) {

	s := "this is a foobarsystem to get bar"

	v := New().Anything().Find("foo").Find("bar").Word()
	res := v.Regex().FindAllStringSubmatch(s, -1)

	if len(res[0]) > 1 {
		t.Errorf("%v is not a slice of only one match (globale match)", res)
	}
	if res[0][0] != "this is a foobarsystem" {
		t.Errorf("global capture \"%s\" is not \"this is a foobarsystem\"", res[0][0])
	}

	v = New().Anything().Find("foo").BeginCapture().Find("bar").Word().EndCapture()
	res = v.Regex().FindAllStringSubmatch(s, -1)

	if len(res) != 1 {
		t.Errorf("%v is not slice length 1", res)
	}

	if res[0][0] != "this is a foobarsystem" {
		t.Errorf("global capture \"%s\" is not \"this is a foobarsystem\"", res[0][0])
	}
	if res[0][1] != "barsystem" {
		t.Errorf("capture %s is not barsystem", res[0][1])
	}

}

func TestSeveralCaptures(t *testing.T) {

	s := `
	this is a foobarsystem that matches my test
	And there, a new foobartest that should be ok
`

	v := New().Anything().Find("foo").
		BeginCapture().
		Find("bar").Word().
		EndCapture().
		SearchOneLine(false)
	res := v.Regex().FindAllStringSubmatch(s, -1)

	for i, r := range res {
		switch i {
		case 0:
			if r[1] != "barsystem" {
				t.Errorf("%s is not \"barsystem\"", r[1])
			}
		case 1:
			if r[1] != "bartest" {
				t.Errorf("%s is not \"bartest\"", r[1])
			}
		default:
			t.Errorf("%v is not allowed result", r)
		}
	}

}

func TestCapturingSeveralGroups(t *testing.T) {

	s := `

<b>test 1</b>
<b>foo 2</b>

`
	v := New().
		Find("<b>").
		BeginCapture().
		Word().
		EndCapture().
		Any(" ").
		BeginCapture().
		Range("0", "9").
		EndCapture().
		Find("</b>")

	res := v.Captures(s)
	if len(res) != 2 {
		t.Errorf("%v is not length 2", res)
	}
	for i, r := range res {
		switch i {
		case 0:
			if r[1] != "test" || r[2] != "1" {
				t.Errorf("%s,%s is not test,1", r[1], r[2])
			}
		case 1:
			if r[1] != "foo" || r[2] != "2" {
				t.Errorf("%s,%s is not test,1", r[1], r[2])
			}
		default:
			t.Errorf("%d is not a valid result index for %v", i, res)
		}
	}

}

func TestORMethod(t *testing.T) {

	s := "foobarbaz footestbaz foonobaz"
	expected := []string{"foobarbaz", "footestbaz"}

	v := New().Find("foobarbaz").Or().Find("footestbaz")
	if !v.Test(s) {
		t.Errorf("%s doesn't match %s", v.Regex(), s)
	}
	res := []string{}
	res = v.Regex().FindAllString(s, -1)

	if len(res) != 2 {
		t.Errorf("%v is not length 2", res)
	}

	for i, r := range res {
		if r != expected[i] {
			t.Errorf("%s is not expected value: %s", r, expected[i])
		}
	}

}
