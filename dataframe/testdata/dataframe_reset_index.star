load("dataframe.star", "dataframe")


def f():
  df = dataframe.DataFrame([("bird", 389.0),
                            ("bird", 24.0),
                            ("mammal", 80.5),
                            ("mammal", 9.2)],
                           index=["falcon", "parrot", "lion", "monkey"],
                           columns=("class", "max_speed"))
  print(df)
  print('')

  df = df.reset_index()
  print(df)
  print('')


f()
