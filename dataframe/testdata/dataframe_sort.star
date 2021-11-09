load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame(columns=['id','animal','sound'],
                           data=[[1,'cat','meow'],
                                 [2,'dog','bark'],
                                 [3,'eel','zap'],
                                 [4,'frog','ribbit']])
  print("rows with no index")
  print(df)
  print("")

  sorted = df.sort_values(by=['sound'])
  print("case 0: sort ascending")
  print(sorted)
  print("")

  sorted = df.sort_values(by=['sound'], ascending=False)
  print("case 1: sort descending")
  print(sorted)
  print("")

  df = dataframe.DataFrame(columns=['id','animal','sound'],
                           data=[[1,'cat','meow'],
                                 [2,'dog','bark'],
                                 [3,'eel','zap'],
                                 [4,'frog','ribbit']],
                           index=['nyan', 'inu', 'unagi', 'kaeru'])
  print("rows with index")
  print(df)
  print("")

  sorted = df.sort_values(by=['sound'])
  print("case 2: sort with index")
  print(sorted)
  print("")


f()
