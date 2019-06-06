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
					should work like [soup.find()](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#find)
			 find_all(name, attrs, recursive, string, limit, **kwargs)
					retrieves all descendants that match arguments passed to find_all.
					should work like [soup.find_all()](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#find-all)
			 attrs()
					get a dictionary of element attributes
					should work like [soup.attrs](https://www.crummy.com/software/BeautifulSoup/bs4/doc/#attributes)
			 contents()
			 child()
       parent()
       next_sibling()
       prev_sibling()
   functions:
     parseHtml(html string) SoupNode
       parseHTML parses html from a string, returning the root SoupNode
*/
package bsoup
