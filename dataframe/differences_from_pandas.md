# Difference from Pandas

This implementation strives to immitate the python version of Pandas as closely as possible, but there inevitably will be some differences due to how starlark behaves.

## bool(Series)

In python, calling bool(Series) raises this exception: "ValueError: The truth
value of a Series is ambiguous. Use a.empty, a.bool(), a.item(), a.any() or a.all()."
Since starlark does not have exceptions, just always return true.

