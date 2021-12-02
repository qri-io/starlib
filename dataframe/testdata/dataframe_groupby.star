load("dataframe.star", "dataframe")


def example_animals(series):
  examples = ['{} {}'.format(series[i], series.name) for i in range(series.size)]
  return ', '.join(examples)


def f():
  df = dataframe.DataFrame({"IDs": ["cat", "dog", "eel", "dog", "cat", "frog", "cat", "eel"],
                            "count": [1, 2, 3, 4, 5, 6, 7, 8]})
  unit_sums = df.groupby(['IDs'])['count'].sum()
  print(unit_sums)
  print(unit_sums.index)
  print('')

  df = dataframe.DataFrame([["cat", "tabby"],
                            ["cat", "black"],
                            ["cat", "calico"],
                            ["dog", "doberman"],
                            ["dog", "pug"]],
                           columns=["species", "breed"])

  num_breeds = df.groupby(['species'])['breed'].count()
  print(num_breeds)
  print(num_breeds.index)
  print('')

  list_of_examples = df.groupby(['species'])['breed'].apply(example_animals)
  print(list_of_examples)
  print('')
  print(type(list_of_examples))
  print(list_of_examples.name)
  print(list_of_examples.index.name)


f()
