
load('encoding/csv.star', 'csv')
load('assert.star', 'assert')

csv_string_1 = """a,b,c
1,2,3
4,5,6
7,8,9
"""

assert.eq(csv.read_all(csv_string_1), [["a","b","c"],["1","2","3"],["4","5","6"],["7","8","9"]])