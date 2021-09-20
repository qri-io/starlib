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
  print('Full DataFrame:')
  print(df)
  print('')

  print('name ends with "g":')
  res = df[df['name'].str.endswith('g')]
  print(res)
  print('')

  print('sound starts with "b":')
  res = df[df['sound'].str.startswith('b')]
  print(res)
  print('')

  print('sound does not start with "b":')
  res = df[~df['sound'].str.startswith('b')]
  print(res)
  print('')

  print('sound\'s 1-th char is "a":')
  res = df[df['sound'].str[1].equals('a')]
  print(res)
  print('')

  print('name\'s 1-th char is not "i":')
  res = df[df['name'].str[1].notequals('i')]
  print(res)
  print('')

  print('name contains an "o":')
  res = df[df['name'].str.contains('o')]
  print(res)
  print('')

  print('sound starts with "b", using list comprehension:')
  res = df[[e.startswith('b') for e in df.sound]]
  print(res)
  print('')

  print('name\'s 2-th char is not "g", using list comprehension:')
  res = df[[e[2] != 'g' for e in df.name]]
  print(res)
  print('')


f()
