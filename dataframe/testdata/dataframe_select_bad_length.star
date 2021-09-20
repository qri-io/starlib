load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([["cat","meow"],
                            ["dog","bark"],
                            ["eel","zap"],
                            ["frog","ribbit"],
                            ["giraffe","hum"],
                            ["hippo","grunt"],
                            ["ibex","bleat"],
                            ["jaguar","roar"]],
                           columns=['name', 'sound'])
  select = [True, False, True, True, False]
  res = df[select]
  print(res)


f()
