load("dataframe.star", "dataframe")


def f():
  left = dataframe.DataFrame([["cat","meow"],
                              ["dog","bark"]])
  rite = dataframe.DataFrame([["eel","zap"],
                              ["frog","ribbit"]])
  df = left + rite
  print(df)
  print('')


f()

