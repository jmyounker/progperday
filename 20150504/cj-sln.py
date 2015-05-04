#!/usr/bin/env python

import sys


def solve():
    m = readline()
    w = m.split()
    w.reverse()
    return " ".join(w)


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


