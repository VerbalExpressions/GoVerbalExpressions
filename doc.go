/*
This is a Go implementation of VerbalExpressions for other languages. Check https://github.com/VerbalExpressions to know the other implementations.

VerbalExperssions is a way to build complex regular expressions with a verbal language.

Important: Because the package is hosted on github organization, the name used to host repository is a bit complicated.
That's why we named package "verbalexpressions" instead of "GoVerbalExpressions". So, to import package you *MUST* do:

	import verbalexpressions "github.com/VerbalExpressions/GoVerbalExpressions"

This is very important !

Using verbalexpressions:

Use "New()" factory then you can chain calls. Go syntax allows you to set new line after seperators:

	v := verbalexpressions.New().
		StartOfLine().
		Find("foo").
		Word().
		Anything().
		EndOfLine()

Then, you can use "Test()" method to check if your string matches expression.

You may get the regexp.Regexp structure using "Regex()" method, then use common methods to split, replace, find submatches and so on... as usual


*/
package verbalexpressions
