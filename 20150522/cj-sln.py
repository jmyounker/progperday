#!/usr/bin/env python

import sys


def solve():
    num, src, dest = readfields(3)

    src_by = {}
    for i, c in enumerate(src):
        src_by[c] = i
    src_base = len(src)

    n = 0
    for c in num:
        n = n * src_base + src_by[c]

    out = []
    if n == 0:
        return dest[0]
    dst_base = len(dest)
    while n != 0:
        r = n % dst_base
        out.append(dest[r])
        n = int(n / dst_base)

    out.reverse()
    return "".join(out)


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
    return itm


def readfields(n):
    itms = readline().split()
    assert len(itms) == n, "wrong input count: expected %d ints but got %d" % (n, len(itms))
    return itms


def readline():
    return sys.stdin.readline().replace('\n', '')


if __name__ == '__main__':
   main()


