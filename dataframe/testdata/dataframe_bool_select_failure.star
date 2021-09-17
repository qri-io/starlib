load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame(columns=["id","animal","sound"],
                           data=[[1,"cat","meow"],
                                 [2,"dog","bark"],
                                 [3,"eel","zap"],
                                 [1,"cat","meow"],
                                 [20,"dog","barks"],
                                 [3,"eel","zap"],
                                 [4,"frog","ribbit"]])
  print(df)
  print("")

  # Invalid way to create a bool series, binary equal operator doesn't work
  is_id_three = df["animal"] == "dog"
  print(is_id_three)
  print("")

  selection = df[is_id_three]
  print(selection)
  print("")


f()
