load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat","meow","purr"],
                            ["dog","bark","woof"]])
  print(df.columns)
  print(df.index)


f()
