load("dataframe.star", "dataframe")


def f():
  series = dataframe.Series([12.3, 45.6, 78.9]).shift(1)
  series.astype('int64')


f()
