load("dataframe.star", "dataframe")
load("time.star", "time")


def f():
  a = time.time(year=2021, month=3, day=21)
  b = time.time(year=2021, month=5, day=4)
  when = dataframe.Series([a, b])
  print(when)
  print('')

  names = dataframe.Series(['apple', 'banana'])
  print(names)
  print('')

  df = dataframe.DataFrame({'when': when, 'names': names})
  print(df)
  print('')

  counts = when.astype('int')
  print(counts)
  print('')

  dts = counts.astype('datetime64[ns]')
  print(dts)
  print('')


f()
