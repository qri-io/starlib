load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat", "meow", 123],
                            ["dog", "bark", 456]])
  print(df)
  print('')

  df = df.append([["eel", "zap", 789]])
  print(df)
  print('')

  other = dataframe.DataFrame([["frog", "ribbit", 321],
                               ["giraffe", "hum", 654],
                               ["hippo", "grunt", 987]])
  df = df.append(other)
  print(df)
  print('')


f()
