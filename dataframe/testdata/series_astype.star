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

  series = dataframe.Series([123.1, 564.5, 978.9])
  print(series)
  print('')

  series = series.shift(1)
  print(series)
  print('')

  out = series.astype('Int64')
  print(out)
  print('')

  out = out[out.notnull()]
  print(out)
  print('')


f()
