#!/usr/bin/env python

import os
import sys


def solve():
    ne, nn = readints(2)
    e = []
    for i in range(0, ne):
        e.append(readline())
    n = []
    for i in range(0, nn):
        n.append(readline())
    de = set()
    dn = set()
    for d in e:
      existing_dirs(de, d)
    for d in n:
        created_dirs(de, dn, d)
    return len(dn)


def existing_dirs(de, d):
    if d == '/':
        return
    if d in de:
        return
    de.add(d)
    existing_dirs(de, os.path.dirname(d))


def created_dirs(de, dn, d):
    if d == "/":
        return
    if d in de:
        return
    if d in dn:
        return
    dn.add(d)
    created_dirs(de, dn, os.path.dirname(d))

    
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


