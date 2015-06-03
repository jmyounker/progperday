#!/usr/bin/env python

"""Skeleton for Dawn's art project driver."""

import time


def empty_img():
    return "empty"


def main():
    image_delay = 1 # seconds
    cam = Camera()
    images = [empty_img() for x in range(0, 6)]
    screens = [s for s in get_screens()[1:]]
    windows = []
    for s in screens:
        w = Window(s.size())
        w.set_screen(s)
        windows.append(w)

    transforms = [None for s in range(0, 6)]
    transforms[0] = ColorBlend(01, 01, 01)
    transforms[1] = ColorBlend(02, 02, 02)
    transforms[2] = ColorBlend(03, 03, 03)
    transforms[3] = ColorBlend(04, 04, 04)
    transforms[4] = ColorBlend(05, 05, 05)
    transforms[5] = ColorBlend(06, 06, 06)

    while True:
        for i in range(5, 0, -1):
            images[i] = images[i-1]
        images[0] = cam.capture()
        for i in range(0, 6):
            windows[i].render(transforms[i](images[i], windows[i]))
        for i in range(0, 6):
            windows[i].show()
        time.sleep(image_delay)


def get_screens():
    return [Screen(i) for i in range(0, 8)]


class Screen(object):
    def __init__(self, number):
        self.number = number

    def size(self):
        return (680, 400) 


def ColorBlend(hue, sat, light):
    def transform(img, window):
        f_pre = "img.raw"
        f_post = "img.adjusted"
        # write file
        # build command to gimp
        batch = '(transform-image "%(fpre)s" "%(fpost)s" %(width)d %(height)d %(hue)d %(sat)s %(light)s)'
        subst = {
            'fpre': f_pre,
            'fpost': f_post,
            'width': window.size[0],
            'height': window.size[1],
            'hue': hue,
            'sat': sat,
            'light': light,
        }
        # run gimp command
        # read file into window object
        return "transformed:" + img
    return transform


class Window(object):
    def __init__(self, size):
        self.size = size

    def set_screen(self, screen):
        pass

    def render(self, img):
        self.img = img

    def show(self):
        print self.img
       

class Camera(object):
    def __init__(self):
        self.cnt = 0

    def capture(self):
        self.cnt += 1
        return "img%d" % self.cnt


if __name__ == '__main__':
     main()
