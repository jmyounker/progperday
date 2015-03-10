#!/usr/bin/env python

"""Insnely stupid and inefficient program to sort an array.

Example:

$ ./randomsort.py
It only took 4797251 exchanges to sort 10 elements.

"""

import sys
import random


ARR_SIZE = 10


def main():
    data = [random.randint(0, 10000) for x in range(0, ARR_SIZE)]
    i = 0
    while not is_sorted(data):
        exc(data)
        i = i + 1
    print "It only took %d exchanges to sort %s elements" % (i, len(data))


def exc(data):
    i = random.randint(0, ARR_SIZE-1)
    j = random.randint(0, ARR_SIZE-1)
    data[i], data[j] = data[j], data[i]


def is_sorted(data):
    c = data[0]
    for i in range(1, len(data)):
        p = c
        c = data[i]
        if p > c:
            return False
    return True


if __name__ == '__main__':
    sys.exit(main())
