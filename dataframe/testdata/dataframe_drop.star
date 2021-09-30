load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat", "meow", 123],
                            ["dog", "bark", 456],
                            ["eel", "zap",  789]],
                           columns=["name", "sound", "num"])
  print(df)
  print('')

  df = df.drop(['sound'], axis=1)
  print(df)
  print('')

  df = df.drop(columns=['num'])
  print(df)
  print('')

  df = df.drop(index=[1])
  print(df)
  print('')


f()
