package bsoup

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestFindByTagName(t *testing.T) {
	script := `
def run():
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
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
  doc = OpenHtml("doc.html")
  elem = doc.find("span")
  return elem.attrs()
result = run()
`
	actual := runTestScript(htmlPage, script)
	// TODO: Investigate why this is non-deterministic, shouldn't starlark order the keys?
	expectOne := `{"id":"abc","style":"color:red;"}`
	expectTwo := `{"style":"color:red;","id":"abc"}`
	if actual != expectOne && actual != expectTwo {
		t.Errorf("error, expected: %s, got %s", expectOne, actual)
	}
}

func runTestScript(htmlContent, scriptContent string) string {
	tmpDir, err := ioutil.TempDir("", "run_script_test")
	if err != nil {
		panic(err)
	}

	htmlFile := filepath.Join(tmpDir, "doc.html")
	err = ioutil.WriteFile(htmlFile, []byte(htmlContent), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.Chdir(tmpDir)
	if err != nil {
		panic(err)
	}

	result, err := runScriptContent(scriptContent, "result")
	if err != nil {
		panic(err)
	}

	result = strings.Replace(result, " ", "", -1)
	result = strings.Replace(result, "\n", "", -1)

	return result
}
