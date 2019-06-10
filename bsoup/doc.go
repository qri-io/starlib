/*Package bsoup defines a beautiful-soup-like API for working with HTML documents
in starlark

 outline: bsoup
   bsoup defines a beautiful-soup-like API for working with HTML documents
   path: bsoup
   types:
     SoupNode
       methods:
         find(name, attrs, recursive, string, **kwargs)
           retrieve the first occurance of an element that matches arguments passed to find.
           works similarly to [node.find()](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#find)
         find_all(name, attrs, recursive, string, limit, **kwargs)
           retrieves all descendants that match arguments passed to find_all.
           works similarly to [node.find_all()](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#find-all)
         attrs()
           get a dictionary of element attributes
           works similarly to [node.attrs](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#attributes)
         contents()
           gets the list of children of an element
           works similarly to [soup.contents](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#contents-and-children)
         child()
           gets a single child element with the given tag name
           works like accessing a node [using its tag name](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#navigating-using-tag-names)
         parent()
           gets the parent node of an element
           works like [node.parent](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#parent)
         next_sibling()
           gets the next sibling of an element
           works like [node.next_sibling](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#next-sibling-and-previous-sibling)
         prev_sibling()
           gets the previous sibling of an element
           works like [node.prev_sibling](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#next-sibling-and-previous-sibling)
         get_text()
           all the text in a document or beneath a tag, as a single Unicode string:
           works like [soup.get_text](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#get-text)
   functions:
     parseHtml(html string) SoupNode
       parseHtml parses html from a string, returning the root SoupNode
*/
package bsoup
