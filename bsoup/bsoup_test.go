package bsoup

import (
	"strings"
	"testing"

	"go.starlark.net/starlark"
)

const htmlPage = `<html>
  <head>
    <title>Sample page</title>
  </head>
  <body>
    <div id="header">header</div>
    <div id="content">
      <ul>
        <li>one</li>
        <li>two</li>
        <li>
          <b>three</b>
        </li>
      </ul>
    </div>
    <div class="footer">
      <span id="abc" style="color:red;">foot</span>
    </div>
  </body>
</html>`

func runTestScript(htmlContent, scriptContent string) string {
	// The environment for tests adds a global function GetHtml to get the html string.
	environment := make(map[string]starlark.Value)
	environment["GetHtml"] = starlark.NewBuiltin("GetHtml", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		// Convert htmlContent string to a starlark value, pass it as an argument to ParseHTML.
		htmlValue := starlark.Value(starlark.String(htmlContent))
		return ParseHTML(thread, b, starlark.Tuple([]starlark.Value{htmlValue}), nil)
	})

	// Execute the script.
	thread := &starlark.Thread{}
	outEnv, err := starlark.ExecFile(thread, "", scriptContent, environment)
	if err != nil {
		panic(err)
	}

	// Remove whitespace so that tests are easier to check.
	result := outEnv["result"].String()
	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "\n", "", -1)
	return result
}

func TestFindByTagName(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  return doc.find("ul")
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<ul><li>one</li><li>two</li><li><b>three</b></li></ul>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestChildByTagName(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  elem = doc.find("ul")
  return elem.child("li")
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<li>one</li>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestChildren(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  elem = doc.find("ul")
  return elem.contents()
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `[,<li>one</li>,,<li>two</li>,,<li><b>three</b></li>,]`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestParent(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  elem = doc.find("b")
  return elem.parent()
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<li><b>three</b></li>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestFindByTagNameAndDict(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  return doc.find("div", {"class": "footer"})
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divclass="footer"><spanid="abc"style="color:red;">foot</span></div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestFindOnlyDict(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  return doc.find("", {"class": "footer"})
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divclass="footer"><spanid="abc"style="color:red;">foot</span></div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestFindOnlyDictUsingId(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  return doc.find("", {"id": "content"})
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divid="content"><ul><li>one</li><li>two</li><li><b>three</b></li></ul></div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestFindKeywordForAttribute(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  return doc.find("", id="content")
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divid="content"><ul><li>one</li><li>two</li><li><b>three</b></li></ul></div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestFindKeywordWithoutTagname(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  return doc.find(id="content")
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divid="content"><ul><li>one</li><li>two</li><li><b>three</b></li></ul></div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestNextSibling(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  elem = doc.find(id="content")
  return elem.next_sibling().next_sibling()
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divclass="footer"><spanid="abc"style="color:red;">foot</span></div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestPrevSibling(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  elem = doc.find(id="content")
  return elem.prev_sibling().prev_sibling()
result = run()
`
	actual := runTestScript(htmlPage, script)
	expect := `<divid="header">header</div>`
	if actual != expect {
		t.Errorf("error, expected: %s, got %s", expect, actual)
	}
}

func TestAttrs(t *testing.T) {
	script := `
def run():
  doc = GetHtml()
  elem = doc.find("span")
  return elem.attrs()
result = run()
`
	actual := runTestScript(htmlPage, script)
	// TODO(dlong): Investigate why this is non-deterministic, shouldn't starlark order the keys?
	expectOne := `{"id":"abc","style":"color:red;"}`
	expectTwo := `{"style":"color:red;","id":"abc"}`
	if actual != expectOne && actual != expectTwo {
		t.Errorf("error, expected: %s, got %s", expectOne, actual)
	}
}
