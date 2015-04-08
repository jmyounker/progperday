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


@app.route('/content/<int:cid>', methods=['GET'])
@returns_json
def content(cid):
    b = 10 * cid
    return {'status': 'ok', 'data': {'elts': range(b, b+10)}}


@app.route('/content_html/<int:cid>', methods=['GET'])
def content_html(cid):
    b = 10 * cid
    return "\n".join("<p>%d</p>" % x for x in range(b, b+10))


@app.route(u'/', methods=[u'GET'])
def index():
    """Returns all test runs"""
    return send_file(os.path.join(static_dir, 'frontend.html'))


def main():
    parser = argparse.ArgumentParser(description='test environment controller')
    app.run(debug=True)


if __name__ == '__main__':
    main()

