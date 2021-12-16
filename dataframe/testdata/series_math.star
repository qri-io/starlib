load("dataframe.star", "dataframe")

def f():
  left = dataframe.Series([123,456,789], index=['a','b','c'])
  print(left)
  print('')

  rite = dataframe.Series([101,505,909], index=['a','b','c'])
  print(rite)
  print('')

  answer = left + rite
  print(answer)
  print('')

  answer = left - rite
  print(answer)
  print('')

  answer = left + rite.shift(1)
  print(answer)
  print('')

  answer = dataframe.abs(left - rite)
  print(answer)
  print('')

  answer = dataframe.abs(left - rite.shift(1))
  print(answer)
  print('')

  answer = left + 1000
  print(answer)
  print('')

  answer = left + 0.5
  print(answer)
  print('')

  left = dataframe.Series(['a', 'b', 'c'])
  print(left)
  print('')

  rite = dataframe.Series(['pple', 'anana', 'herry'])
  print(rite)
  print('')

  answer = left + rite
  print(answer)
  print('')

  answer = left + ' ' + rite
  print(answer)
  print('')


f()
