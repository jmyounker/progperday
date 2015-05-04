//
//  MoonView.swift
//  lander
//
//  Created by Jeff Younker on 5/4/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

import Cocoa
import Quartz

class MoonView: NSView {

    override func drawRect(dirtyRect: NSRect) {
        super.drawRect(dirtyRect)

        // Draw a mountain range
        let deltaX = 3
        let deltaXCrag = 2
        let deltaY = 5
        let deltaYCrag = 4
        let floor = 30
        let ceiling = 200
        
        let width = self.frame.size.width
        let height = self.frame.size.height
        var h = Int(random(floor, ceiling))
        var dir = h > 70 ? -1 : 1

        var bp = NSBezierPath()
        bp.moveToPoint(CGPointMake(0, CGFloat(h)))
        for var i = 0; i < Int(width); i += deltaX {
            // Every so often make a little crag
            if oneIn(7) {
                i += (deltaXCrag - deltaX)
                h += deltaYCrag * -dir
            } else {
                h += deltaY * dir
            }
            // Less often start a new peak or valley
            if oneIn(20) {
                dir *= -1
            }
            // If below the minimum height then start going up
            if h <= floor && dir == -1 {
                dir = 1
            }
            // If above the maximum height then start going down
            if h >= ceiling && dir == 1 {
                dir = -1
            }
            let p = CGPointMake(CGFloat(min(i, Int(width))), CGFloat(h))
            bp.lineToPoint(p)
        }
        bp.stroke()
    }
}

func random(lower: Int, upper: Int) -> Int {
    let r = arc4random_uniform(UInt32(Int64(upper) - Int64(lower)))
    return Int(Int64(r) + Int64(lower))
}

func oneIn(prob: Int) -> Bool {
    return 1 == random(0, prob)
}