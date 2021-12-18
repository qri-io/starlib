load("dataframe.star", "dataframe")
load("time.star", "time")


def f():
  a = time.time(year=2021, month= 3, day=21, hour=18, minute=50, second=44)
  b = time.time(year=2021, month= 5, day= 4, hour= 7, minute=32, second=19)
  c = time.time(year=2021, month=12, day= 7, hour=15, minute=26, second= 2)
  d = time.time(year=2021, month=12, day= 9, hour= 9, minute=10, second=39)
  later = dataframe.Series([a, b, c, d])
  print('later:')
  print(later)
  print('')

  s = time.time(year=2021, month= 3, day= 9, hour=14, minute=52, second= 7)
  t = time.time(year=2021, month= 4, day=14, hour=21, minute=44, second=26)
  u = time.time(year=2020, month= 2, day=11, hour= 6, minute=21, second=50)
  v = time.time(year=2009, month=11, day= 7, hour=13, minute=37, second=14)
  earlier = dataframe.Series([s, t, u, v])
  print('earlier:')
  print(earlier)
  print('')

  diff = later - earlier
  print('diff:')
  print(diff)
  print('')

  days = diff.astype('timedelta64[D]')
  print('days:')
  print(days)
  print('')

  months = diff.astype('timedelta64[M]')
  print('months:')
  print(months)
  print('')

  years = diff.astype('timedelta64[Y]')
  print('years:')
  print(years)
  print('')

  hours = diff.astype('timedelta64[h]')
  print('hours:')
  print(hours)
  print('')

  minutes = diff.astype('timedelta64[m]')
  print('minutes:')
  print(minutes)
  print('')

  seconds = diff.astype('timedelta64[s]')
  print('seconds:')
  print(seconds)
  print('')


f()
