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

  values = dataframe.DataFrame([['2021-01-23', '11:47:22'],
                                ['2021-06-04', '19:03:59']],
                               columns=['day', 'clock'])
  print(values)
  print('')

  when = values.apply(lambda row: time.parse_time('{}T{}Z'.format(row['day'], row['clock'])), axis=1)
  print(when)
  print('')

  counts = when.astype('int')
  print(counts)
  print('')

  values['when'] = when
  print(values)
  print('')

f()
