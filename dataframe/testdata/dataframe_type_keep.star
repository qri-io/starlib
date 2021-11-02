load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat", "meow", 1.2, 3]])
  df = df.append([["zebra", "neigh", 456, 78]])
  print(df)
  print('')


f()
