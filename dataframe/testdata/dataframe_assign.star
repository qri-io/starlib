load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([[1,"cat","meow"],
                            [2,"dog","bark"],
                            [3,"eel","zap"]],
                           columns=["id","animal","sound"])
  print(df)
  print('')

  ans = df.assign(rating=[8.2, 7.6, 3.4])
  print(ans)
  print('')

  # Original is not modified
  print(df)
  print('')


f()
