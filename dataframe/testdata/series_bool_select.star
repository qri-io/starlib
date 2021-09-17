load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series(data=['123','cat','456', '789', 'cat', 'cat', '321'])
  print(series)
  print('')

  bools = series.equals('cat')
  print(bools)
  print('')

  res = series[bools]
  print(res)
  print('')


f()
