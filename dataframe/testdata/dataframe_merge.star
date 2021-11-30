load("dataframe.star", "dataframe")


def f():
  df1 = dataframe.DataFrame({"lkey": ["foo", "bar", "baz", "foo"],
                             "value": [1, 2, 3, 5]})
  df2 = dataframe.DataFrame({"rkey": ["foo", "bar", "baz", "foo"],
                             "value": [5, 6, 7, 8]})
  print(df1)
  print('')
  print(df2)
  print('')

  df3 = df1.merge(df2, left_on="lkey", right_on="rkey")
  print(df3)
  print('')

  df3 = df1.merge(df2, left_on="lkey", right_on="rkey", how="left")
  print(df3)
  print('')

  df3 = df1.merge(df2, left_on="lkey", right_on="rkey",
                  suffixes=("_left", "_right"))
  print(df3)
  print('')

  df1 = dataframe.DataFrame({"num": [1, 2, 3, 4],
                             "animal": ["cat", "dog", "eel", "frog"],
                             "score": [15, 21, 9, 12]})
  df2 = dataframe.DataFrame({"num": [3, 4, 2, 1],
                             "animal": ["frog", "eel", "dog", "cat"],
                             "score": [17, 11, 23, 8]})

  df3 = df1.merge(df2, left_on="num", right_on="num")
  print(df3)
  print('')

  df3 = df1.merge(df2, left_on="animal", right_on="animal")
  print(df3)
  print('')


f()
