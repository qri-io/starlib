load("assert.star", "assert")
load("dataframe.star", "dataframe")


def myFunc(row):
  return '{}:{}'.format(row['lkey'], row['value'])


def f():
  df = dataframe.DataFrame({"lkey": ["foo", "bar", "baz", "foo"],
                            "value": [1, 2, 3, 5]})
  r = df.apply(myFunc, axis=1)
  print(r)
  print('')

  df['combined'] = r
  print(df)
  print('')


f()
