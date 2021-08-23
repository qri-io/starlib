load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat","meow"],
                            ["dog","bark"],
                            ["eel","zap"],
                            ["frog","ribbit"]])
  print(df)
  print('')

  idx = df.index
  print(idx)
  print('')

  print(len(idx))
  print(df.shape)
  print('')

  for v in idx:
    print(v)
  print('')

  df = dataframe.DataFrame([["cat","meow"],
                            ["dog","bark"],
                            ["eel","zap"],
                            ["frog","ribbit"]],
                           index=["Garfield", "Clifford", "Abaia", "Kermit"])
  print(df)
  print('')

  idx = df.index
  print(idx)
  print('')

  print(len(idx))
  print(df.shape)
  print('')

  for v in idx:
    print(v)
  print('')


f()

