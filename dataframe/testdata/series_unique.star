load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series(['123', '456', '123', '789', '456', '042', '456', '555'])
  print(series)
  print('')

  out = series.unique()
  print(out)
  print('')


f()
