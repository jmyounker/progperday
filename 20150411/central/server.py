#!/usr/bin/env python

__author__ = 'jeff'


import argparse
from functools import wraps
import json
import os


# Using flask because basically all we need is URL routing to
# interact via JSON.
from flask import Flask, request, Response, send_file


base_dir = os.path.dirname(os.path.dirname(os.path.realpath(__file__)))
static_dir = os.path.join(base_dir, 'static')
app = Flask(__name__, static_folder=static_dir, static_url_path='/static')


def returns_json(f):
    @wraps(f)
    def decorator(*args, **kwargs):
        r = f(*args, **kwargs)
        return Response(json.dumps(r), content_type='text/json; charset=utf-8')
    return decorator


def returns_plain(f):
    @wraps(f)
    def decorator(*args, **kwargs):
        r = f(*args, **kwargs)
        return Response(r, content_type='text/plain; charset=utf-8')
    return decorator


@app.route('/data', methods=['GET'])
@returns_json
def content():
    return {
        'status': 'ok',
        'data': [
        dict(id=1, content='item 1', start='2014-04-20'),
        dict(id=2, content='item 2', start='2014-04-14'),
        dict(id=3, content='item 3', start='2014-04-19'),
        dict(id=4, content='item 4', start='2014-04-22'),
        dict(id=5, content='item 5', start='2014-04-15'),
        dict(id=6, content='item 6', start='2014-04-12'),
        dict(id=7, content='item 7', start='2014-04-01'),
    ]
    }


@app.route('/', methods=[u'GET'])
def index():
    """Returns all test runs"""
    return send_file(os.path.join(static_dir, 'frontend.html'))


def main():
    parser = argparse.ArgumentParser(description='show a timeline using vis.js')
    app.run(debug=True)


if __name__ == '__main__':
    main()

