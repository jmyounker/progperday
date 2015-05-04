#!/usr/bin/env python

import sys


def lookup():
    s = 'abcdefghijklmnopqrstuvwxyz'
    lkup = {}
    n = 2
    i = 1
    for c in s:
        lkup[c] = str(n) * i
        i += 1
        if i > 3 and not (c == 'r' or c == 'y'):
            i = 1
            n += 1
    lkup[' '] = "0"
    return lkup


lut = lookup()

def solve():
    line = readline()
    r = ''
    p = ''
    for x in line:
        v = lut[x]
        if len(p) > 0 and p[0] == v[0]:
            r = r + " "
        r = r + v
        p = v
    return r


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


