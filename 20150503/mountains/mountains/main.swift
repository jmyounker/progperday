//
//  main.swift
//  mountains
//
//  Created by Jeff Younker on 5/4/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

import Foundation
import Cocoa

func main() {
    var app = NSApplication()
    
    let frame = NSMakeRect(0, 0, 400, 200)
    var window = NSWindow(
        contentRect: frame,
        styleMask: NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask,
        backing: NSBackingStoreType.Buffered,
        defer: false)
    var view = MoonView()
    window.title = "The Moon"
    window.contentView.addSubview(view)
    window.delegate = view
    window.makeKeyAndOrderFront(nil)
    app.run()
}


class MoonView : NSView, NSWindowDelegate {
    
}
main()
