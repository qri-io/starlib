import os
import subprocess
import sys
import tempfile


"""comparator.py compares the output of python and starlark dataframe scripts

We want the python and starlark versions of dataframe to work as close as
possible to each other. By running our test scripts using python, then diffing
the output, we can compare how they behave.

In some cases, we expect differences, because the python version may print
internal implementation details. For example, type(...) returns something
like `<class 'pandas.core.series.Series'>` while our starlark version simply
uses `dataframe.Series`. For these sort of situations, allow.txt contains
a list of differences that are allowed to appear.

TODO(dustmop): Add expected ERROR cases to allow.txt

See status.txt for the current differences between the implementations

TODO(dustmop): The goal is to get status.txt to be all `ok`s. Some cases need
bug fixes, other need errors to be added to allow.txt

TODO(dustmop): Remove the `-wB` flag from our diff commands, so that even
differences in whitespace are handled.

Run this program like this to see the status of all scripts:

  >  python comparator.py all

Run this program like this to check a single script (the 5th, in this case):

  >  python comparator.py 5

"""


def mkdir_p(path):
  try:
    os.makedirs(path)
  except OSError as e:
    if e.errno == errno.EEXIST:
      pass
    else: raise


def read_file(filename, mode='r'):
  fp = open(filename, mode)
  contents = fp.read()
  fp.close()
  return contents


def write_file(filename, contents, mode='w'):
  fp = open(filename, mode)
  fp.write(contents)
  fp.close()


def is_error(obj):
  return isinstance(obj, dict)


def get_files(glob):
  # TODO: This is a total hack, assumes that the '*' always appears at the
  # start of the basename. Example: '/path/to/*.star' works correctly,
  # but '/path/to/my_*.star' does not.
  (dirpath, suffix) = glob.split('*')
  ents = os.listdir(dirpath)
  return sorted([os.path.join(dirpath, e) for e in ents if e.endswith(suffix)])


def parse_allow_file(filepath):
  content = read_file(filepath)
  result = {}
  index, accum, before, after = [None]*4
  state = 'start'
  for i, line in enumerate(content.split('\n')):
    numline = i + 1
    if state == 'start':
      if not line:
        continue
      if line.isdigit():
        index = int(line)
        accum, before, after = [[], [], []]
        state = 'header'
        continue
    elif state == 'header':
      if line.startswith('==='):
        accum = []
        state = 'before'
        continue
      raise Exception('%s:%d: expected "===" after test case number, got "%s"' % (filepath, numline, line))
    elif state == 'before':
      if line.startswith('---'):
        before += accum
        accum = []
        state = 'after'
        continue
      if line.startswith('< '):
        accum.append(line[2:])
        continue
    elif state == 'after':
      if not line:
        if accum:
          after += accum
        result[index] = {'before': before, 'after': after}
        state = 'start'
        continue
      if line.startswith('==='):
        after += accum
        accum = []
        state = 'before'
        continue
      if line.startswith('> '):
        accum.append(line[2:])
        continue
    # If we got here, there's an error
    raise Exception('%s:%d: unexpected, got "%s" [state=%s] TODO: improve this error' % (filepath, numline, line, state))
  return result


def diff_files(expectfile, actualfile):
  """Diff the files, and display detailed output if any"""

  cmd = ['diff', '-wB', expectfile, actualfile]
  proc = subprocess.Popen(cmd, stdin=subprocess.PIPE,
                          stdout=subprocess.PIPE, stderr=sys.stderr)
  (stdout, _unused) = proc.communicate()
  if isinstance(stdout, bytes):
    stdout = stdout.decode('utf-8')

  if not stdout:
    print('no difference')
    return

  print('Expect: %s' % expectfile)
  print('Actual: %s' % actualfile)
  print('command: diff -w %s %s' % (expectfile, actualfile))
  print(stdout)


def has_difference(expectfile, actualfile):
  """Diff the files, return whether there's differences"""

  cmd = ['diff', '-wB', expectfile, actualfile]
  proc = subprocess.Popen(cmd, stdin=subprocess.PIPE,
                          stdout=subprocess.PIPE, stderr=sys.stderr)
  (stdout, _unused) = proc.communicate()
  if isinstance(stdout, bytes):
    stdout = stdout.decode('utf-8')
  return stdout


def run_python_script(filepath, allow_diff):
  tmpdir = tempfile.mkdtemp()
  srcfile = os.path.join(tmpdir, 'run.py')
  content = read_file(filepath)
  content = content.replace('load("dataframe.star", "dataframe")',
                            'import pandas as dataframe')
  write_file(srcfile, content)

  cmd = ['python', srcfile]
  proc = subprocess.Popen(cmd, stdin=subprocess.PIPE,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE)
  (stdout, stderr) = proc.communicate()
  if stderr:
    return {'err': stderr.decode('utf-8')}
  if isinstance(stdout, bytes):
    stdout = stdout.decode('utf-8')

  if allow_diff:
    stdout = apply_allow_diff(stdout, allow_diff)

  outfile = os.path.join(tmpdir, 'actual.txt')
  write_file(outfile, stdout)
  return outfile


def apply_allow_diff(text, allow_diff):
  result = []
  c = 0
  needle = allow_diff['after']
  replace = allow_diff['before']
  for line in text.split('\n'):
    if c < len(needle) and line == needle[c]:
      line = replace[c]
      c += 1
    result.append(line)
  return '\n'.join(result)


def compare_single_scriptfile(scriptfile, allow_diff):
  """Run a single script with python, show detailed diff output"""

  actualfile = run_python_script(scriptfile, allow_diff)
  if is_error(actualfile):
    print(actualfile['err'])
    return
  expectfile = scriptfile.replace('.star', '.expect.txt')
  if not os.path.isfile(expectfile):
    print('no expectation file for %s' % scriptfile)
  else:
    diff_files(expectfile, actualfile)


def compare_all_scriptfiles(all_files, allow_conf):
  """Run all scripts with python, show short summary of diffs"""

  for i, scriptfile in enumerate(all_files):
    actualfile = run_python_script(scriptfile, allow_conf.get(i))
    if is_error(actualfile):
      print('%2d: ERROR %s' % (i, scriptfile))
      continue
    expectfile = scriptfile.replace('.star', '.expect.txt')
    if not os.path.isfile(expectfile):
      print('%2d: no expect file %s' % (i, scriptfile))
    elif has_difference(expectfile, actualfile):
      print('%2d: DIFFERENT %s' % (i, scriptfile))
    else:
      print('%2d: ok %s' % (i, scriptfile))


def show_usage():
  usage = """usage: python comparator.py ['all' | which]

    all     Run all scripts, compare output of python and starlark
    which   An integer, run the single script and show differences
"""
  # TODO: Option to ignore whitespaces differences (current: true)
  print(usage)


def main():
  which = sys.argv[1] if len(sys.argv) > 1 else None
  star_files = get_files('../testdata/*.star')
  allow_conf = parse_allow_file('allow.txt')

  if which is None:
    show_usage()
    return

  if which == 'all':
    compare_all_scriptfiles(star_files, allow_conf)
    return

  index = int(which)
  compare_single_scriptfile(star_files[index], allow_conf.get(index))


if __name__ == '__main__':
  main()
