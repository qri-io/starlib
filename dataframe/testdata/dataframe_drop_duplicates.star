load("dataframe.star", "dataframe")
load("time.star", "time")


def f():
  df = dataframe.DataFrame([[1,"cat","meow"],
                            [2,"dog","bark"],
                            [3,"eel","zap"],
                            [1,"cat","meow"],
                            [20,"dog","barks"],
                            [3,"eel","zap"],
                            [4,"frog","ribbit"]],
                           columns=["id","animal","sound"])
  print(df)
  print('')

  no_dups = df.drop_duplicates()
  print(no_dups)
  print('')

  no_dup_animals = df.drop_duplicates(subset=["animal"])
  print(no_dup_animals)
  print('')

  ts0 = time.time(year=2021, month=1, day=23)
  ts1 = time.time(year=2021, month=5, day=19)
  ts2 = time.time(year=2021, month=7, day=4)
  df = dataframe.DataFrame([["cat", ts0],
                            ["cat", ts1],
                            ["dog", ts2]],
                           columns=["animal","when"])
  print(df)
  print('')

  no_dup_animals = df.drop_duplicates(subset=["animal"])
  print(no_dup_animals)
  print('')


f()
