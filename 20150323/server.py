#!/usr/bin/env python

__author__ = 'jeff'

# Using flask because basically all we need is URL routing to
# interact via JSON.
from flask import Flask, redirect, url_for

# Endpoints
# GET - /run/ - all test IDS
# GET - /run/:id: - get data for a single test
# POST - /run/:id: - create test ID
app = Flask(__name__)


@app.route(u'/run', methods=[u'GET'])
@app.route(u'/run/', methods=[u'GET'])
def all_runs():
    return "GOT RUN REQUEST"


@app.route(u'/run/<int:run_id>', methods=[u'GET', u'POST'])
@app.route(u'/run/<int:run_id>/', methods=[u'GET', u'POST'])
def one_run(run_id):
    return "GOT RUN_ID %d" % run_id


if __name__ == "__main__":
    app.run(debug=True)


