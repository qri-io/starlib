load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series([123,456,789])
  print(series)
  print('')

  df = series.to_frame()
  print(df)
  print('')


f()
