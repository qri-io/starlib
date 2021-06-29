load("dataframe.star", "dataframe")


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

  no_dups = df.drop_duplicates()
  print(no_dups)

  no_dup_animals = df.drop_duplicates(subset=['animal'])
  print(no_dup_animals)


f()
