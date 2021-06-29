load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series(['abc', None, '', 'def', None])
  print(series)
  print('')

  nn = series.notnull()
  print(nn)
  print('')

  out = series[series.notnull()]
  print(out)
  print('')


f()
