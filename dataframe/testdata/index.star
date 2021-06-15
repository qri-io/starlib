load("assert.star", "assert")
load("dataframe.star", "dataframe")

fruits = dataframe.Index(["apple", "orange", "banana", "banana", "lemon", "apple", "banana"])

assert.eq(
  fruits.eq("banana"),
  dataframe.Index(data=[False, False, True, True, False, False, True])
)