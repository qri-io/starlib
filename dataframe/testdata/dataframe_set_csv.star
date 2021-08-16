load("dataframe.star", "dataframe")


text = """a,b,c
1,2,3
4,5,6
"""


def f():
  df = dataframe.DataFrame()
  df.set_csv(text)
  print(df)


f()
