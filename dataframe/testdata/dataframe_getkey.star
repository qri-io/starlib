load("assert.star", "assert")
load("dataframe.star", "dataframe")

df = dataframe.read_csv("""id,animal,sound
1,cat,meow
2,dog,bark
3,eel,zap
""")

series = df['animal']
print(series)
