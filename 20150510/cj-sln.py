#!/usr/bin/env python

import sys


def solve(t):
    w = readline()
    pp = parsepossible(w)
    return findwords(t, pp)


def parsepossible(w):
    st = []
    START = 1
    ONECHAR = 2
    STATE = START
    cw = set()
    for c in w:
        if c == '(' and STATE == START:
            STATE = ONECHAR
        elif c == ')' and STATE == ONECHAR:
            st.append(cw)
            cw = set()
            STATE = START
        elif STATE == ONECHAR:
            cw.add(c)
        elif STATE == START:
            cw.add(c)
            st.append(cw)
            cw = set()
        else:
            raise Exception("BADPARSEIMPLEMENTATION")
    return st


def findwords(dtree, w):
    if len(w) == 0:
        return 1
    wf = 0
    for c in w[0]:
        if c in dtree:
            wf += findwords(dtree[c], w[1:])
    return wf


def trie(words):
    t = {}
    for w in words:
        ct = t
        for c in w:
            if c not in ct:
                ct[c] = {}
            ct = ct[c]
        ct = ['WORD']
    return t


def main():
    (N, D, L) = readints(3)

    words = []
    for _ in range(1, D+1):
        words.append(readline())
    t = trie(words)

    for i in range(1, L+1):
        res = solve(t)
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


