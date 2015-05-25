#!/usr/bin/env python

"""Demonstrate inference algorithm."""

import sys


class TFun(object):
    def __init__(self, name, args, retv):
        self.name = name
        self.args = args
        self.retv = retv


class TCon(object):
    def __init__(self, name):
        self.name = name


class TVar(object):
    def __init__(self, name):
        self.name = name


class TApp(object):
    def __init__(self, name, args):
        self.name = name
        self.args = args


def main():
    pass


if __name__ == '__main__':
    sys.exit(main())
