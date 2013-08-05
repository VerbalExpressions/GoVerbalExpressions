GoVerbalExpressions
===================

Go VerbalExpressions implementation

VerbalExpression is a concept to help building difficult regular expressions.

## Other Implementations
You can see an up to date list of all ports in our [organization](https://github.com/VerbalExpressions).
- [Javascript](https://github.com/jehna/VerbalExpressions)
- [Ruby](https://github.com/VerbalExpressions/RubyVerbalExpressions)
- [C#](https://github.com/VerbalExpressions/CSharpVerbalExpressions)
- [Python](https://github.com/VerbalExpressions/PythonVerbalExpressions)
- [Java](https://github.com/VerbalExpressions/JavaVerbalExpressions)
- [PHP](https://github.com/VerbalExpressions/PHPVerbalExpressions)
- [C++](https://github.com/VerbalExpressions/CppVerbalExpressions)
- [Haskell](https://github.com/VerbalExpressions/HaskellVerbalExpressions)


## Installation

Use this command line:
    
    go get github.com/VerbalExpressions/GoVerbalExpressions

This will install package in your $GOPATH and you will be ready to import it.

## Examples

```go

// import with a nice name
import (
    verbalexpressions "github.com/VerbalExpressions/GoVerbalExpressions"
    "fmt"
)

func main () {
    v := verbalexpressions.New().
            StartOfLine()
            Then("http").
            Maybe("s").
            Then( "://").
            Maybe("www.").
            AnythingBut(" ").
            EndOfLine()

    testMe := "https://www.google.com"
    
    if v.Test(testMe) {
       fmt.Println("You have a valid URL") 
    } else {
       fmt.Println("URL is incorrect") 
    }
}

```


