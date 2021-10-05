load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([[1,"cat","meow"],
                            [2,"dog","bark"],
                            [3,"eel","zap"]],
                           columns=["id","animal","sound"])
  print(df)
  print('')

  df = dataframe.DataFrame({"a": ["apple", "apricot"],
                            "b": ["banana", "blueberry"],
                            "c": ["cherry", "currant"]})
  print(df)
  print('')

  df = dataframe.DataFrame({"id": [123, 456, 789],
                            "city": ["New York", "Chicago", "San Jose"],
                            "pop": [8419000, 2710000, 1028000]})
  print(df)
  print('')

  series = df["city"]
  print(series)
  print(type(series))
  print('')

  df = dataframe.DataFrame([[1,"cat","meow"],
                            [2,"dog","bark"],
                            [3,"eel","zap"]],
                           index=["feline","canine","anguillid"])
  print(df)
  print('')

  df = dataframe.DataFrame([["a", "b", "c", True, 2],
                            ["d", "e", False, False, 3]])
  print(df)
  print('')

  one_column = ["cat", "dog", "eel"]
  df = dataframe.DataFrame(one_column)
  print(df)
  print('')

f()

