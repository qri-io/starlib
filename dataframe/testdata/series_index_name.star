load("dataframe.star", "dataframe")

def f():
  index = dataframe.Index(data=['a', 'b', 'c'], name='IDs')
  series = dataframe.Series([123,456,789], index=index, name='Cool Data')
  print(series)
  print('')

f()
