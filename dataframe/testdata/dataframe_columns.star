load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame({" First ": ["foo", "bar", "baz", "foo"],
                            "Second": [1, 2, 3, 5],
                            "Third Column": [300.1, 301.2, 303.4, 305.7]})

  col = df.columns
  print(df.columns)

  col = df.columns.str.lower()
  print(col)

  col = df.columns.str.replace(' ', '_')
  print(col)

  col = df.columns.str.strip()
  print(col)

  col = df.columns.str.strip().str.lower().str.replace(' ', '_')
  print(col)
  print('')

  df.columns = col
  print(df)
  print('')

  new = dataframe.DataFrame()
  new['my_column'] = [123, 456, 789]
  print(new)
  print('')

  new['another_col'] = ['cat', 'dog', 'eel']
  print(new)
  print('')


f()
