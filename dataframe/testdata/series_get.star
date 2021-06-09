load("dataframe.star", "dataframe")

def f():
  print('case 0:')
  series = dataframe.Series(data=[123,456,789], index=['a', 'b', 'c'])
  print(series['a'])
  print(series[1])
  print(type(series[1]))
  print('')

  print('case 2:')
  series = dataframe.Series(data=['123','456','789'], index=['a', 'b', 'c'])
  print(series['a'])
  print(series[1])
  print(type(series[1]))
  print('')

  print('case 3:')
  series = dataframe.Series(data=[123.4,456.7,789.1], index=['a', 'b', 'c'])
  print(series['a'])
  print(series[1])
  print(type(series[1]))
  print('')

  print('case 4:')
  series = dataframe.Series(data=[123,456,789], index=['a', 'b', 'c'])
  print(series.get(0))
  print(series.get('a'))
  print('')

f()
