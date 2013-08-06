// Copyright 2013 Patrice FERLET
// Ue of this source code is governed by MIT-style
// license that can be found in the LICENSE file
package verbalexpressions

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// VerbalExpression structure to create expression
type VerbalExpression struct {
	expression string
	suffixes   string
	prefixes   string
	modifiers  string
}

// quote is an alias to regexp.QuoteMeta
func quote(s string) string {
	return regexp.QuoteMeta(s)
}

// utility function to return only strings
func tostring(i interface{}) string {
	var r string
	switch x := i.(type) {
	case string:
		r = x
	case int64:
		r = strconv.FormatInt(x, 64)
	case uint:
		r = strconv.FormatUint(uint64(x), 64)
	case int:
		r = strconv.FormatInt(int64(x), 32)
	default:
		log.Panicf("Could not convert %v %t", x, x)
	}
	return r
}

// Instanciate a new VerbalExpression. You should use this method to
// initalize some internal var.
//
// Example:
//		v := verbalexpression.New().Find("foo")
func New() *VerbalExpression {
	r := new(VerbalExpression)
	r.modifiers = "m"
	return r
}

// append a modifier
func (v *VerbalExpression) addmodifier(m string) *VerbalExpression {
	if !strings.Contains(v.modifiers, m) {
		v.modifiers += "m"
	}
	return v
}

// remove a modifier
func (v *VerbalExpression) removemodifier(m string) *VerbalExpression {
	v.modifiers = strings.Replace(v.modifiers, m, "", -1)
	return v
}

// add method, append expresions to the internal string that will be parsed
func (v *VerbalExpression) add(s string) *VerbalExpression {
	v.expression += s
	return v
}

// Start to capture something, stop with EndCapture()
func (v *VerbalExpression) BeginCapture() *VerbalExpression {
	v.suffixes += ")"
	return v.add("(")
}

// Stop capturing expresions parts
func (v *VerbalExpression) EndCapture() *VerbalExpression {
	v.suffixes = strings.Replace(v.suffixes, ")", "", 1)
	return v.add(")")
}

// Anything will match any char
func (v *VerbalExpression) Anything() *VerbalExpression {
	return v.add(`(?:.*)`)
}

// AnythingBut will match anything excpeting the given string.
// Example:
//		s := "This is a simple test"
//		v := verbalexpressions.New().AnythingBut("ie").RegExp().FindAllString(s, -1)
//		[Th s  s a s mple t st]
func (v *VerbalExpression) AnythingBut(s string) *VerbalExpression {
	return v.add(`(?:[^` + quote(s) + `]*)`)
}

// EndOfLine tells verbalexpressions to match a end of line.
// Warning, to check multiple line, you must use SearchOneLine(true)
func (v *VerbalExpression) EndOfLine() *VerbalExpression {
	v.suffixes += "$"
	return v
}

// Maybe will search string zero on more times
func (v *VerbalExpression) Maybe(s string) *VerbalExpression {
	return v.add(`(?:` + quote(s) + `)?`)
}

// StartOfLine seeks the begining of a line. As EndOfLine you should use
// SearchOneLine(true) to test multiple lines
func (v *VerbalExpression) StartOfLine() *VerbalExpression {
	v.prefixes += `^`
	return v
}

// Find seeks string. The string MUST be there (unlike Maybe() method)
func (v *VerbalExpression) Find(s string) *VerbalExpression {
	return v.add(`(?:` + quote(s) + `)`)
}

// Alias to Find()
func (v *VerbalExpression) Then(s string) *VerbalExpression {
	return v.Find(s)
}

// Any accepts caracters to be matched
//
// Example:
//		s := "foo1 foo5 foobar"
//		v := New().Find("foo").Any("1234567890").Regex().FindAllString(s, -1)
//		[foo1 foo5]
func (v *VerbalExpression) Any(s string) *VerbalExpression {
	return v.add(`(?:[` + quote(s) + `])`)
}

//AnyOf is an alias to Any
func (v *VerbalExpression) AnyOf(s string) *VerbalExpression {
	return v.Any(s)
}

// LineBreak to find "\n" or "\r\n"
func (v *VerbalExpression) LineBreak() *VerbalExpression {
	return v.add(`(?:(?:\n)|(?:\r\n))`)
}

// Alias to LineBreak
func (v *VerbalExpression) Br() *VerbalExpression {
	return v.LineBreak()
}

// Range accepts an even number of arguments. Each pair of values defines start and end of range.
// Think like this: Range(from, to [, from, to ...])
//
// Example:
//		s := "This 1 is 55 a TEST"
//		v := verbalexpressions.New().Range("a","z",0,9)
//		res := v.Regex().FindAllString()
//		[his 1 is 55 a]
func (v *VerbalExpression) Range(args ...interface{}) *VerbalExpression {
	if len(args)%2 != 0 {
		log.Panicf("Range: not even args number")
	}

	parts := make([]string, 3)
	app := ""
	for i := 0; i < len(args); i++ {
		app += tostring(args[i])
		if i%2 != 0 {
			parts = append(parts, quote(app))
			app = ""
		} else {
			app += "-"
		}
	}
	return v.add("[" + strings.Join(parts, "") + "]")
}

// Tab fetch tabulation char (\t)
func (v *VerbalExpression) Tab() *VerbalExpression {
	return v.add(`\t+`)
}

// Word matches any word (containing alpha char)
func (v *VerbalExpression) Word() *VerbalExpression {
	return v.add(`\w+`)
}

// Or, chains a alternate expression
//
// Example:
//		v := Verbalexpression.New().
//				Find("foobarbaz").
//				Or().
//				Find("footestbaz")
func (v *VerbalExpression) Or() *VerbalExpression {
	v.prefixes += "(?:"
	v.suffixes = ")" + v.suffixes
	return v.add(")|(?:")
}

// WithAnyCase asks verbalexpressions to match with or without case sensitivity
func (v *VerbalExpression) WithAnyCase(sensitive bool) *VerbalExpression {
	if !sensitive {
		return v.removemodifier("i")
	}
	return v.addmodifier("i")
}

// SearchOneLine deactivates "multiline" mode if online argument is true
// Default is false
func (v *VerbalExpression) SearchOneLine(oneline bool) *VerbalExpression {
	if oneline {
		return v.removemodifier("m")
	}
	return v.addmodifier("m")
}

// Regex returns the regular expression to use to test on string.
func (v *VerbalExpression) Regex() *regexp.Regexp {
	modifier := ""
	if len(v.modifiers) > 0 {
		modifier = "(?" + v.modifiers + ")"
	}

	return regexp.MustCompile(modifier + v.prefixes + v.expression + v.suffixes)
}

/* proxy and helpers to regexp.Regexp functions */

// Test return true if verbalexpressions matches something in string "s"
func (v *VerbalExpression) Test(s string) bool {
	return v.Regex().Match([]byte(s))
}

// Replace alias to regexp.ReplaceAllString. It replace the found expression from
// string src by string dst
func (v *VerbalExpression) Replace(src string, dst string) string {
	return v.Regex().ReplaceAllString(src, dst)
}

// Returns a slice of results from captures. If you didn't apply BeginCapture() and EnCapture(), the slices
// will return slice of []string where []string is length 1, and 0 index is the global capture
//
// Example:
//		s:="This should get barsystem and whatever..."
//		// get "bar" followed by a word
//		v := verbalexpressions.New().Anything().
//				BeginCatpure().
//				Find("bar").Word().
//				EndCapture()
//
//		res := v.Captures(s)
//		fmt.Println(res)
//		[
//			["This should get barsystem", "barsystem"] // 0: global capture, 1: catpure 1
//		]
//
// So, to range results, you can do:
//		for _, captures := range res {
//			fmt.Println(captures[1])
//		}
// Actualy, 1 matches first group, you can use several captures.
func (v *VerbalExpression) Captures(s string) [][]string {
	return v.Regex().FindAllStringSubmatch(s, -1)
}
