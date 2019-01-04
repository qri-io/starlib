/*Package html defines a jquery-like html selection & iteration functions for HTML documents

  outline: html
    html defines a jquery-like html selection & iteration functions for HTML documents

    functions:
      html(markup) selection
        parse an html document returing a selection at the root of the document
        params:
          markup string
            html text to build a document from

    types:
      selection
        an HTML document for querying
        methods:
          attr(name) string
            gets the specified attribute's value for the first element in the Selection.
            To get the value for each element individually, use a looping construct such as each or map method
            params:
              name string
                attribute name to get the value of
          children() selection
            gets the child elements of each element in the Selection
          children_filtered(selector) selection
            gets the child elements of each element in the Selection, filtered by the specified selector
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          contents(selector) selection
            gets the children of each element in the Selection, including text and comment nodes
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          find(selector) selection
            gets the descendants of each element in the current set of matched elements, filtered by a selector
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          filter(selector) selection
            filter reduces the set of matched elements to those that match the selector string
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          get(i) selection
            retrieves the underlying node at the specified index. alias: eq
            params:
              i int
                numerical index of node to get
          has(selector) selection
            reduces the set of matched elements to those that have a descendant that matches the selector
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          parent(selector) selection
            gets the parent of each element in the Selection
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          parents_until(selector) selection
            gets the ancestors of each element in the Selection, up to but not including the element matched by the selector
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          siblings() selection
            gets the siblings of each element in the Selection
          text() string
            gets the combined text contents of each element in the set of matched elements, including descendants
          first(selector) selection
            gets the first element of the selection
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          last() selection
            gets the last element of the selection
            params:
              selector string
                a query selector string to filter the current selection, returning a new selection
          len() int
            returns the number of the nodes in the selection
          eq(i) selection
            gets the element at index i of the selection
            params:
              i int
                numerical index of node to get

*/
package html
