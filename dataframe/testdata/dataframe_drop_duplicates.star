load("assert.star", "assert")
load("dataframe.star", "dataframe")

df = dataframe.read_csv("""id,animal,sound
1,cat,meow
2,dog,bark
3,eel,zap
1,cat,meow
20,dog,barks
3,eel,zap
4,frog,ribbit""")
print(df)

no_dups = df.drop_duplicates()
print(no_dups)

no_dup_animals = df.drop_duplicates(subset=['animal'])
print(no_dup_animals)
