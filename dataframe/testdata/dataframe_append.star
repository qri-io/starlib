load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat", "meow", 123],
                            ["dog", "bark", 456]])
  print(df)
  print(df.index)
  print('')

  # Append a list of lists (only 1 row)
  df = df.append([["eel", "zap", 789]])
  print(df)
  print(df.index)
  print('')

  # Append a dataframe
  other = dataframe.DataFrame([["frog", "ribbit", 321],
                               ["giraffe", "hum", 654],
                               ["hippo", "grunt", 987]])
  df = df.append(other)
  print(df)
  print('')

  # Append with not enough columns
  other = dataframe.DataFrame([["iguana", "wheeze"]])
  df = df.append(other)
  print(df)
  print('')

  # Append with too many columns
  other = dataframe.DataFrame([["jaguar", "growl", 444, 555]])
  df = df.append(other)
  print(df)
  print('')

  # Append to an empty dataframe
  df = dataframe.DataFrame()
  df = df.append(other)
  print(df)
  print('')

  df = dataframe.DataFrame([["cat", "meow", 123],
                            ["dog", "bark", 456]],
                           index=["c", "d"])
  print(df)
  print(df.index)
  print('')

  df = df.append([["eel", "zap", 789]])
  print(df)
  print(df.index)
  print('')


f()
