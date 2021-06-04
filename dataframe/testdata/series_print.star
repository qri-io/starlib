load("dataframe.star", "dataframe")

def f():
  series = dataframe.Series([chr(c) for c in range(ord('a'), ord('m'))])
  print(series)
  print('')

  series = dataframe.Series(data=[123,456,789], index=['abc', 'd', 'fghi'])
  print(series)
  print('')

f()
