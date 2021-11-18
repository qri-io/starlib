load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame(columns=['id','animal','sound'],
                           data=[[1,'cat','meow'],
                                 [2,'cat','purr'],
                                 [3,'dog','bark'],
                                 [4,'dog','woof'],
                                 [5,'eel','zap'],
                                 ])
  print(df)
  print("")

  series = df.groupby(['animal'])['id'].count()
  print(series)
  print("")

  series = series.reset_index()
  print(series)
  print("")


f()
