load("dataframe.star", "dataframe")


def f():
  rows = [["test", 31.2, 11.4, "ok", 107, 6.91],
          ["more",  7.8, 44.1, "hi",  94, 3.1],
          ["last", 90.2, 26.8, "yo", 272, 4.3]]
  df = dataframe.DataFrame(rows)
  print(df)
  print('')

  a = df.at
  print(a)
  print('')

  v = df.at[1,3]
  print(v)
  print('')

  df.at[0,4] = 567
  print(df)
  print('')

  df.at[1,0] = ["dog", "eel"]
  df.at[2,0] = {"cat": "meow"}
  print(df)
  print('')


f()
