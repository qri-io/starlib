load("dataframe.star", "dataframe")


def f():
  df = dataframe.read_csv("""id,animal,sound
1,cat,meow
2,dog,bark
3,eel,zap
""")
  print(df)

  df = dataframe.DataFrame({"a": ["apple", "apricot"],
                            "b": ["banana", "blueberry"],
                            "c": ["cherry", "currant"]})
  print(df)

  df = dataframe.DataFrame({"city": ["New York", "Chicago", "San Jose"],
                            "pop": [8419000, 2710000, 1028000]})
  print(df)


f()

