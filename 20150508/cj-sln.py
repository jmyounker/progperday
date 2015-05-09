#!/usr/bin/env python

import sys


def solve():
    n = readint()
    v1 = readints(n)
    v2 = readints(n)
    vs1 = sorted(v1)
    vs2 = sorted(v2)
    vs2.reverse()
    return dp(vs1, vs2)


def dp(v1, v2):
    s = 0
    for i in range(0, len(v1)):
        s += v1[i]*v2[i]
    return s


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


