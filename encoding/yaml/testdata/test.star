
load('encoding/yaml.star', 'yaml')
load('assert.star', 'assert')

yaml_list = """- Apple
- Orange
- Strawberry
- Mango
"""

native_list = ["Apple", "Orange", "Strawberry", "Mango"]
assert.eq(yaml.loads(yaml_list), native_list) 
assert.eq(yaml.dumps(native_list), yaml_list)

yaml_dict = """martin:
  job: Developer
  name: Martin D'vloper
  skill: Elite
"""

native_dict = {"martin": {"name": "Martin D'vloper", "job": "Developer", "skill": "Elite"}}
assert.eq(yaml.loads(yaml_dict), native_dict)
assert.eq(yaml.dumps(native_dict), yaml_dict)

