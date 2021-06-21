load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series(['123', '456', '789'])
  print(series)
  print('')

  out = series.astype('int64')
  print(out)
  print('')

  series = dataframe.Series(['2010-03-22 23:20:50', '2012-04-07 11:07:43'], dtype='datetime64[ns]')
  print(series)
  print('')

  out = series.astype('int64')
  print(out)
  print('')


f()
