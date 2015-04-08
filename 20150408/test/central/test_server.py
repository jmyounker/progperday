__author__ = 'jeff'


import json

import pytest

import central.server


@pytest.fixture
def clt(app):
    return app.test_client()


@pytest.fixture
@pytest.fixture
def app():
    return central.server.app


def test_statics(clt):
    assert clt.get('/', follow_redirects=True).status == '200 OK'
    assert clt.get('/static/frontend.html', follow_redirects=True).status == '200 OK'
    assert clt.get('/static/frontend.js', follow_redirects=True).status == '200 OK'
    assert clt.get('/static/jquery-1.11.2.js', follow_redirects=True).status == '200 OK'


def test_content(clt):
    res = clt.get('/content/2')
    assert status_ok(res)
    jres = json.loads(res.data)
    assert jres['data']['elts'] == range(20, 30)


def assert_error_msg(res, msg):
    jres = json.loads(res.data)
    assert 'status' in jres
    assert 'error' in jres
    assert jres['status'] == 'fail'
    assert msg in jres['error']


def status_ok(res):
    return res.status == '200 OK'

