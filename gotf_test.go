package gotf

import (
	"bytes"
	"testing"
	"time"
)

func AssertEqual(t *testing.T, buffer *bytes.Buffer, testString string) {
	if buffer.String() != testString {
		t.Errorf("Expected %s, got %s", testString, buffer.String())
	}
	buffer.Reset()
}

func ParseTest(buffer *bytes.Buffer, body string, data interface{}) {
	tpl := New("test")
	tpl.Parse(body)
	tpl.Execute(buffer, data)
}

func TestGotfFuncMap(t *testing.T) {
	var buffer bytes.Buffer
	testtime, _ := time.Parse("2006-Jan-02", "2013-Feb-03")
	testtimestuct := struct {
		T time.Time
	}{
		T: testtime,
	}
	testcase := []struct {
		body   string
		expect string
		data   interface{}
	}{
		{"{{ \"Hello World\" | stringsToLower }}", "hello world", ""},
		{"{{ \"Hello World\" | stringsToUpper }}", "HELLO WORLD", ""},
		{"{{ \"    Hello World    \" | stringsTrimSpace }}", "Hello World", ""},
		{"{{ \"Hello World\" | stringsReplace \"Hello\"  \"HELLO\" -1  }}", "HELLO World", ""},
		{"{{ \"a=b?c\" | urlQueryEscape }}", "a%3Db%3Fc", ""},
		{"{{ \"2013-Feb-03\" | timeParse \"2006-Jan-02\" | structFormat \"Jan-02 2006\" }}", "Feb-03 2013", ""},
		{"{{ .T | structFormat \"Jan-02 2006\" }}", "Feb-03 2013", testtimestuct},
	}

	for _, acase := range testcase {
		ParseTest(&buffer, acase.body, acase.data)
		AssertEqual(t, &buffer, acase.expect)
	}
}
