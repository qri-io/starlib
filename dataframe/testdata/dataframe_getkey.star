load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame({"id": [1, 2, 3],
                            "animal": ["cat", "dog", "eel"],
                            "sound": ["meow", "bark", "zap"]})
  print(df)
  print('')

  series = df['animal']
  print(series)
  print('')

  series = df['id']
  print(series)
  print('')


f()
