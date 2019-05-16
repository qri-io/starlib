load('xlsx.star', 'xlsx')
load('assert.star', 'assert')

file_1 = xlsx.get_url(test_server_url + "/eg_one.xlsx")

assert.eq(file_1.get_sheets(), { 1 : "Sheet1"})
assert.eq(file_1.get_rows("Sheet1"), [["foo", "bar", "baz"],["bat", "cat","mouse"]])
