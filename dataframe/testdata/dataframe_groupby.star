load("dataframe.star", "dataframe")


def example_animals(series):
  examples = ['{} {}'.format(elem, series.name) for elem in series]
  return ', '.join(examples)


def make_full_name_series(series):
  return series + ' ' + series.name


def f():
  df = dataframe.DataFrame({"IDs": ["cat", "dog", "eel", "dog", "cat", "frog", "cat", "eel"],
                            "count": [1, 2, 3, 4, 5, 6, 7, 8]})
  print(df)
  print('')

  print('case 0: group and sum a grouped column')
  unit_sums = df.groupby(['IDs'])['count'].sum()
  print(unit_sums)
  print(unit_sums.index)
  print('')

  df = dataframe.DataFrame([["cat", "tabby"],
                            ["cat", "black"],
                            ["cat", "calico"],
                            ["dog", "doberman"],
                            ["dog", "pug"]],
                           columns=["species", "breed"],
                           index=['A','B','C','D','E'])

  print('case 1: group and count a grouped column')
  num_breeds = df.groupby(['species'])['breed'].count()
  print(num_breeds)
  print(num_breeds.index)
  print('')

  # If `apply` returns a scalar, the result is a Series whose
  # length is equal to the number of groupings. The index is
  # the keys of the grouping.
  list_of_examples = df.groupby(['species'])['breed'].apply(example_animals)
  print(list_of_examples)
  print('')
  print('case 2: apply a function that returns a scalar')
  print('type = %s' % type(list_of_examples))
  print('name = %s' % list_of_examples.name)
  print('index.name = %s' % list_of_examples.index.name)
  print('index = %s' % list_of_examples.index)
  print('')

  # If `apply` returns a Series, the result is a Series whose
  # length is equal to the original non-grouped Series. The
  # index matches the original input DataFrame.
  list_of_examples = df.groupby(['species'])['breed'].apply(make_full_name_series)
  print(list_of_examples)
  print('')
  print('case 3: apply a function that returns a series')
  print('type = %s' % type(list_of_examples))
  print('name = %s' % list_of_examples.name)
  print('index.name = %s' % list_of_examples.index.name)
  print('index = %s' % list_of_examples.index)

f()
