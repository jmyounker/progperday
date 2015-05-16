#!/usr/bin/env python

import sys


def solve():
    n = readint()
    ws = []
    for i in range(0, n):
        (s, r) = readints(2)
        ws.append((s, r))
    ws = sorted(ws, key=lambda x: -x[0])
    cross = 0
    for i in range(0, len(ws)-1):
      wte = ws[i][1]
      for j in range(i+1, len(ws)):
          wle = ws[j][1]
          if wle > wte:
              cross += 1
    return str(cross)
         
        
       


def main():
    ncl = readline()
    nc = int(ncl)
    for i in range(1, nc+1):
        res = solve()
        print "Case #%d: %s" % (i, res)


def readint():
    return readints(1)[0]


def readints(n):
    itms = [int(x) for x in readline().split()]
    assert len(itms) == n, "wrong input count: expected %d ints but got %d" % (n, len(itms))
    return itms


def readline():
    return sys.stdin.readline().replace('\n', '')


if __name__ == '__main__':
   main()


