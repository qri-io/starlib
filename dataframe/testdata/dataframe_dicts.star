load("dataframe.star", "dataframe")


def f():
  print("case 0: columns with dict will re-index the data")
  df = dataframe.DataFrame({"a": ["apple", "apricot"],
                            "b": ["banana", "blueberry"],
                            "c": ["cherry", "currant"],
                            "d": ["date", "durian"]},
                           columns=["a","b","o","d","e"])
  print(df)
  print(df.columns)
  print('')


  print("case 1: list of dicts will have correct column names")
  rows = [{"name": "cat", "sound": "meow"},
          {"name": "dog", "sound": "bark"},
          {"name": "eel", "sound": "zap"}]
  df = dataframe.DataFrame(rows)
  print(df)
  print(df.columns)
  print('')


  print("case 2: columns are merged by name, null values for missing cells")
  rows = [{"month": "June", "year": 2001, "day": 4},
          {"day": 10, "year": 1996, "weekday": "Monday"},
          {"month": "December", "day": 25, "weekday": "Wednesday"}]
  df = dataframe.DataFrame(rows)
  print(df)
  print(df.columns)
  print('')


f()

