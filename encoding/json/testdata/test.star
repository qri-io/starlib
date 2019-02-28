
load('encoding/json.star', 'json')
load('assert.star', 'assert')

json_string_1 = """[["a","b","c"],
[1,2,3],
[4,5,6],
[7,8,9]
]
"""

json_1 = json.loads(json_string_1)
assert.eq(json_1, [["a","b","c"],[1,2,3],[4,5,6],[7,8,9]])
assert.eq(json.dumps(json_1), "[[\"a\",\"b\",\"c\"],[1,2,3],[4,5,6],[7,8,9]]")