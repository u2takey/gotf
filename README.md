# gotf
Use golang functions in template

## Aim
Map golang functions to template functions. 
Now you can use golang functions in template.

## Buildin Functions

### stringsToLower
```
{{ "Hello World" | stringsToLower }} => "hello world"
```

### stringsToUpper
```
{{ "Hello World" | stringsToUpper }} => "HELLO WORLD"
```

### stringsTrimSpace
```
{{ "    Hello World    " | stringsTrimSpace }} => "Hello World"
```

### stringsReplace
```
{{ "Hello World" | stringsReplace "Hello"  "HELLO" -1  }} => "HELLO World"
```

### urlQueryEscape
```
{{ \"a=b?c\" | urlQueryEscape }}" => "a%3Db%3Fc"
```

### timeParse
```
{{ "2013-Feb-03" | timeParse "2006-Jan-02" | structFormat "Jan-02 2006" }} => "Feb-03 2013"
```

### structFormat
```
{{ "2013-Feb-03" | timeParse "2006-Jan-02" | structFormat "Jan-02 2006" }} => "Feb-03 2013"
```

## Add functions From Std library

```
// gotf.Makefunc make handler used for template.FuncMap
// name is function name, if function is a class method, prefix it by "stuct"
// valIndex is input value index in gofunc
// gofunc is funtion you want to map into template
func Makefunc(name string, valIndex int, gofunc interface{}) handler 


// usage
import "strings"
GoFuncMap = template.FuncMap{"stringsReplace" : Makefunc("stringsReplace", 0, strings.Replace)}

template.New(name).Funcs(GoFuncMap)
```