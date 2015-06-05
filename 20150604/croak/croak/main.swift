//
//  main.swift
//  croak
//
//  Created by Jeff Younker on 6/4/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

import Foundation

func soundServer() {
    var server:UDPServer = UDPServer(addr: "127.0.0.1", port: 9975);
    println("Waiting for messages")

    let frog1 = "/Users/jeff/repos/progperday/20150414/audemio/Rana_clamitans.mp3"
    let frog2 = "/Users/jeff/repos/progperday/20150414/audemio/rana_sylvatica_chorus.aiff"
    var ap = AudioPlayback(soundFiles: [frog1, frog2])
    ap.begin()
    while true {
        let (packet, addr, port) = server.recv(1024*10)
        if packet == nil {
            break
        }
        let frognum = packet!.count > 3 ? 1 : 0
        println("starting frog \(frognum)")
        ap.stopPlayback(frognum)
        ap.startPlayback(frognum)
    }
    ap.end()
}

soundServer()