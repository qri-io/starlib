load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([['cat', 'tabby', 123],
                            ['cat', 'black', 456],
                            ['cat', 'calico', 789],
                            ['dog', 'doberman', 321],
                            ['dog', 'pug', 654]],
                           columns=['species', 'breed', 'id'])
  print(df)
  print('')

  mod = df.shift(1)
  print('case 0: shift 1')
  print(mod)
  print('')

  mod = df.shift(4)
  print('case 1: shift 4')
  print(mod)
  print('')

  mod = df.shift(1, axis='columns')
  print('case 2: shift 1 by columns')
  print(mod)
  print('')

  series = df['species']
  print('case 3: series')
  print(series)
  print('')

  mod = series.shift(1)
  print('case 4: series shift 1')
  print(mod)
  print('')

  mod = series.shift(3)
  print('case 5: series shift 3')
  print(mod)
  print('')


f()
