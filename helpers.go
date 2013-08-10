package verbalexpressions

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
