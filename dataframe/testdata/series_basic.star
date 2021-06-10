load("dataframe.star", "dataframe")

def f():
  print('case 0:')
  series = dataframe.Series([123,456,789])
  print(series)
  print('')

  print('case 1:')
  series = dataframe.Series(data=[123,456,789])
  print(series)
  print('')

  print('case 2:')
  series = dataframe.Series(data=['123','456','789'])
  print(series)
  print('')

  print('case 3:')
  series = dataframe.Series(data=[123.4,456.7,789.1])
  print(series)
  print('')

  print('case 4:')
  series = dataframe.Series(3)
  print(series)
  print('')

  print('case 5:')
  series = dataframe.Series(123.4)
  print(series)
  print('')

  print('case 6:')
  series = dataframe.Series('123')
  print(series)
  print('')

  print('case 7:')
  series = dataframe.Series(data=3)
  print(series)
  print('')

  print('case 8:')
  series = dataframe.Series(3, dtype='object')
  print(series)
  print('')

  print('case 9:')
  series = dataframe.Series(3, dtype='float64')
  print(series)
  print('')

  print('case 10:')
  series = dataframe.Series({'a': 123, 'b':456, 'c':789})
  print(series)
  print('')

  print('case 11:')
  series = dataframe.Series(data={'a': 123, 'b':456, 'c':789})
  print(series)
  print('')

  print('case 12:')
  series = dataframe.Series([123,456,789], index=['a', 'b', 'c'])
  print(series)
  print('')

  print('case 13:')
  series = dataframe.Series(data=[123,456,789], index=['a', 'b', 'c'])
  print(series)
  print('')

  print('case 14:')
  series = dataframe.Series({'a': 123, 'b':456, 'c':789}, index=['a', 'b', 'c'])
  print(series)
  print('')

  print('case 15:')
  series = dataframe.Series(data={'a': 123, 'b':456, 'c':789}, index=['a', 'b', 'c'])
  print(series)
  print('')

  # TODO:
  #print('case 16:')
  #series = dataframe.Series(data={'a': 123, 'b':456, 'c':789}, index=['b', 'a', 'd'])
  #print(series)
  #print('')

  # TODO:
  #print('case 17:')
  #series = dataframe.Series(data={'a': 123, 'b':456, 'c':789}, index=['x', 'y', 'z'])
  #print(series)
  #print('')

  print('case 18:')
  series = dataframe.Series(data=[123,456,789], name='numbers')
  print(series)
  print('')

  print('case 19:')
  series = dataframe.Series({'a': 123, 'b':456, 'c':789}, index=['a', 'b', 'c'], name='numbers')
  print(series)
  print('')

  print('case 20:')
  series = dataframe.Series(data=[123,456,789], name='numbers', dtype='int64')
  print(series)
  print('')

  print('case 21:')
  series = dataframe.Series(data=[123,456,789], name='numbers', dtype='float64')
  print(series)
  print('')

  print('case 22:')
  series = dataframe.Series(data=[123,456,789], name='numbers', dtype='string')
  print(series)
  print('')

  print('case 23:')
  series = dataframe.Series(data=[12,34.9,56])
  print(series)
  print('')

  print('case 24:')
  series = dataframe.Series(data=[12.9,34,56])
  print(series)
  print('')

  print('case 25:')
  series = dataframe.Series(data=[12,'34',56])
  print(series)
  print('')

  print('case 26:')
  series = dataframe.Series(data=[12.9,'34.8',56.7])
  print(series)
  print('')

f()
