load("assert.star", "assert")
load("dataframe.star", "dataframe")

df = dataframe.DataFrame({"id": [1,2,3],
                          "animal": ["cat","dog","eel"],
                          "sound": ["meow","bark","zap"]})

print('case 0:')
df['num'] = 7
print(df)

print('case 1:')
df['num'] = 8
print(df)

print('case 2:')
df['num'] = 'abc'
print(df)

series = dataframe.Series([123,456,789])

print('case 3:')
df['num'] = series
print(df)

print('case 4:')
df['more'] = series
print(df)
