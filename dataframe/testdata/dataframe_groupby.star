load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame({"IDs": ["cat", "dog", "eel", "dog", "cat", "frog", "cat", "eel"],
                            "count": [1, 2, 3, 4, 5, 6, 7, 8]})
  unit_sums = df.groupby(['IDs'])['count'].sum()
  print(unit_sums)


f()
