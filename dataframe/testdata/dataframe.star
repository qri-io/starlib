load("assert.star", "assert")
load("dataframe.star", "dataframe")

df = dataframe.DataFrame(data={
  "id":      ("1", "2", "3", "4", "5", "6", "7"),
  "name":    ("apple", "orange", "banana", "banana", "lemon", "apple", "banana"),
  "weight":  (1.1, 2.0, 3.3, 3.5, 2, 1.4, 1.4),
  "spoiled": (False, False, True, False, False, True, False),
})

assert.eq(df.weight, (1.1, 2.0, 3.3, 3.5, 2, 1.4, 1.4))
