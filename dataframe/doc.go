/*Package dataframe ...

  outline: dataframe
    dataframe is a 2d columnar data structure that provides many powerful analysis and manipulation tools, similar to a spreadsheet or SQL engine
    path: dataframe
    functions:
      DataFrame(data, index, columns, dtype) DataFrame
        constructs a DataFrame containing the given data
        params:
          data any
            data for the content of the DataFrame. Can be a list, dict, Series, or another DataFrame
          index Index
            an Index that describes the rows
          columns Index
            an Index that describes the columns
          dtype string
            data type to force. If not provided, it will be inferred for each column
      parse_csv(text) DataFrame
        constructs a DataFrame by parsing the text as csv data. Assumes the first row is a header row
        params:
          text string
            the string to parse as csv data
      Index(data, name) Index
        constructs an Index, which describes a single axis of a dataframe
        params:
          data list(string)
            a list of strings for the index
          name string
            the name of the Index
      Series(data, index, dtype, name) Series
        constructs an Series, a homogeneously typed dataframe column
        params:
          data list
            a list of data values. They will be coerced to use a single data type
          index Index
            the index that describes the elements in the Series
          dtype string
            data type of the values in the Series
          name string
            name of the Series
    types:
      DataFrame
        a dataframe
        methods:
          append(other) DataFrame
            appends data to the rows of this DataFrame, returned as a new DataFrame
            params:
              other list
                data to append
          apply(function, axis) Series
            travel the given axis and apply the function to each slice. The result values of that function are collected into a Series, which is returned
            params:
              function function
                the function to apply to each slice
              axis int
                which to travel, either 0 for columns, or 1 for rows
          drop(labels, axis, index, columns)
            drop columns or rows from the DataFrame
            params:
              labels list(string)
                what to drop from the DataFrame, axis is required to specify what the labels mean. axis=0 if the labels are for the index, axis=1 if the labels are for the columns
              axis int
                which axis to drop from. axis=0 for index, axis=1 for columns
              index list(string)
                values to drop from the index
              columns Index
                values to drop from the columns
          drop_duplicates(subset)
            drop duplicate rows of the DataFrame
            params:
              subset list(string)
                which subset of each row to consider for uniqueness
          group_by(by) GroupByResult
            group a set of row according to some given column value
            params:
              by list(string)
                a list of column names to use for grouping the rows together
          head(n?) DataFrame
            return the first n row of the DataFrame
            params:
              n int
                number of rows to include, defaulting to 5
          merge(right, left_on, right_on, how, suffixes) DataFrame
            merge this with the right DataFrame, returned as a new DataFrame
            params:
              right DataFrame
                the DataFrame to merge with this one
              left_on string
                which column of the left DataFrame to merge on
              right_on string
                which column of the right DataFrame to merge on
              how string
                how to merge the columns, only "inner" is supported, and is the default
              suffixes list(string)
                suffixes to use for merged column names, defaulting to ["_x", "_y"]
          reset_index()
            resets the index to be an empty index, turning the previous index into its own column
        fields:
          at AtIndexer
            returns an AtIndexer, which can be used to retrieve an arbitrary cell from the DataFrame
          columns Index
            returns the columns of the DataFrame as an Index
          index Index
            returns the Index of the DataFrame, if it exists
          shape tuple(int,int)
            returns a tuple with the size of the DataFrame, as (number rows, number columns)

*/
package dataframe
