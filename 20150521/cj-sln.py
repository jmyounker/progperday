#!/usr/bin/env python

import sys


def solve():
    cr1 = readint()
    a1 = []
    for i in range(0, 4):
        a1.append(readints(4))
    cr2 = readint()
    a2 = []
    for i in range(0, 4):
        a2.append(readints(4))
    r1 = set(a1[cr1-1])
    r2 = set(a2[cr2-1])
    nc = 0
    fc = None
    for c in r1:
        if c in r2:
            nc += 1
            fc = c
    if nc == 0:
        return "Volunteer cheated!"
    elif nc == 1:
        return str(fc)
    else:
        return "Bad magician!"


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


