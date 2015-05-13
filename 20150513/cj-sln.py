#!/usr/bin/env python

import sys


def solve():
    ms, audraw = readline().split()
    ms = int(ms)
    aud = [int(x) for x in audraw]
    stand = 0
    add = 0
    for i, aud in enumerate(aud):
        if stand < i:
            add += i - stand
            stand = i
        stand += aud
    return str(add)


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


