package gotf

import (
	"html/template"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type functomap struct {
	name       string
	gofunc     interface{}
	valueindex int
}

var (
	t                = time.Now()
	StructFuncPrefix = "struct"
	GotfFuncMap      = template.FuncMap{}
	golangFuncs      = []functomap{
		{"stringsToLower", strings.ToLower, 0},
		{"stringsToUpper", strings.ToUpper, 0},
		{"stringsTrimSpace", strings.TrimSpace, 0},
		{"stringsReplace", strings.Replace, 0},
		{"urlQueryEscape", url.QueryEscape, 0},
		{"timeParse", time.Parse, 1},
		{"structFormat", t.Format, -1},
	}
)

type handler func(args ...interface{}) reflect.Value

// gotf.Makefunc make handler used for template.FuncMap
// name is function name, if function is a class method, prefix it by "stuct"
// valIndex is input value index in gofunc
// gofunc is funtion you want to map into template
func Makefunc(name string, valIndex int, gofunc interface{}) handler {
	return func(args ...interface{}) (ret reflect.Value) {
		var retv []reflect.Value

		inputs := make([]reflect.Value, len(args))
		leninput := len(inputs)
		for i, _ := range args {
			inputs[i] = reflect.ValueOf(args[i])
		}

		if valIndex >= leninput {
			return
		} else if valIndex != leninput-1 && valIndex != -1 {
			s := make([]reflect.Value, 0)
			s = append(s, inputs[:valIndex]...)
			s = append(s, inputs[leninput-1])
			s = append(s, inputs[valIndex:leninput-1]...)
			inputs = s
		}

		var fv reflect.Value
		ft := reflect.TypeOf(gofunc)

		if strings.HasPrefix(name, StructFuncPrefix) && leninput > 0 {
			method := strings.TrimPrefix(name, StructFuncPrefix)
			fv = inputs[leninput-1].MethodByName(method)
			inputs = inputs[:leninput-1]
		} else {
			fv = reflect.ValueOf(gofunc)
		}
		for i := 0; i != ft.NumIn(); i++ {
			it := ft.In(i)
			if it != inputs[i].Type() {
				ret = reflect.ValueOf("param not valid")
				return
			}
		}
		retv = fv.Call(inputs)

		switch len(retv) {
		case 0:
		case 1:
			ret = retv[0]
		case 2:
			ret = retv[0]
			if err, ok := retv[1].Interface().(error); ok {
				if err != nil {
					ret = reflect.ValueOf(err)
				}
			}
		}
		return
	}
}

func init() {
	for _, tomap := range golangFuncs {
		GotfFuncMap[tomap.name] = Makefunc(tomap.name, tomap.valueindex, tomap.gofunc)
	}
}

// gotf.New returns template.Template with GotfFuncMap added.
func New(name string) *template.Template {
	return template.New(name).Funcs(GotfFuncMap)
}

// gotf.Inject injects gotf functions into the passed FuncMap.
// Set force with true if you want to replace same name function.
// Set prefix to avoid function conflict
func Inject(funcs map[string]interface{}, force bool, prefix string) {
	for k, v := range GotfFuncMap {
		if _, ok := funcs[k]; !ok || force {
			funcs[k+prefix] = v
		}
	}
}
