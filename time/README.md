# time
time defines time primitives for starlark

## Functions

#### `duration(string) duration`
parse a duration

#### `location(string) location`
parse a location

#### `now() time`
implementations would be able to make this a constant

#### `time(string, format=..., location=...) time`
parse a time

#### `zero() time`
a constant


## Types

### `duration`


**Fields**

| name | type | description |
|------|------|-------------|
| hours | float |  |
| minutes | float |  |
| nanoseconds | int |  |
| seconds | float |  |
### `time`


**Fields**

| name | type | description |
|------|------|-------------|
| year | int |  |
| month | int |  |
| day | int |  |
| hour | int |  |
| minute | int |  |
| second | int |  |
| nanosecond | int |  |
