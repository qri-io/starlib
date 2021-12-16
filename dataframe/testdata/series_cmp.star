load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series([111, 123, 190, 212, 142, 100, 180])
  result = series.cmp('<', 150)
  print(result)


f()
