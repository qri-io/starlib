load('html.star', 'html')
load('assert.star', 'assert')

htmlstr = """
<!DOCTYPE html>
<html>
<head>
  <title>example page</title>
</head>
<body id="A" class="le_body">
  <header id="B">Header Tag</header>
  <section id="C">
    <div id="D" class="le_div">
      <h1 id="E" data-name="foo">Heading One</h1>
      <p id="F" class="paragraph">Paragraph</p>
    </div>
  </section>
  <footer id="G">
    <a id="H" href="http://foo.com">link</a>
  </footer>
</body>
</html>
"""

doc = html(htmlstr)

assert.eq(doc.children().text(), "\n  example page\n\n\n  Header Tag\n  \n    \n      Heading One\n      Paragraph\n    \n  \n  \n    link\n  \n\n\n")
assert.eq(doc.find("#E").text(), "Heading One")
assert.eq(doc.find("body header").text(), "Header Tag")
assert.eq(doc.find("#F").text(), "Paragraph")
assert.eq(doc.find("#E").attr("data-name"), "foo")
assert.eq(doc.contents().text(), "\n  example page\n\n\n  Header Tag\n  \n    \n      Heading One\n      Paragraph\n    \n  \n  \n    link\n  \n\n\n")
assert.eq(doc.find("footer").children().attr("href"), "http://foo.com")
assert.eq(doc.find("#D").children_filtered("#E").attr("data-name"), "foo")
assert.eq(doc.find("#D").children().filter("h1").attr("data-name"), "foo")
assert.eq(doc.has("#B").find("header").text(), "Header Tag")
assert.eq(doc.find("#A").children().first().text(), "Header Tag")
assert.eq(doc.find("#A").children().last().text(), "\n    link\n  ")
assert.eq(doc.find("#A").children().eq(0).text(), "Header Tag")
assert.eq(doc.find("#A").children().len(), 3)


p = doc.find("p")
assert.eq(p.parent().attr("class"), "le_div")
assert.eq(p.parents_until("body").attr("class"), "le_div")
assert.eq(p.siblings().text(), "Heading One")
assert.eq(p.get(), ("p",))
assert.eq(p.get(0), "p")

