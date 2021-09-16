load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame()
  df.ffill()


f()
