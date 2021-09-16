load("dataframe.star", "dataframe")


def f():
  df = dataframe.parse_csv("""id,animal,sound
1,cat,meow
2,dog,bark
3,eel,zap
1,cat,meow
20,dog,barks
3,eel,zap
4,frog,ribbit""")
  print(df)
  print('')

  print(df.head())
  print('')

  print(df.head(3))
  print('')


f()
