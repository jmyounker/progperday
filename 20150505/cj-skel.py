#!/usr/bin/env python

import sys


def main(f):
  ncl = readline(f)
  nc = int(ncl)
  for i in range(1, nc+1):
    res = solve(f)
    print "Case #%d: %s" % (i, res)


def solve(f):
    return "foo"


def readline(f):
  return f.readline().replace('\n', '')


if __name__ == '__main__':
   main(sys.stdin)


