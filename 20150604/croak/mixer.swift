//
//  mixer.swift
//  audemio
//
//  Created by Jeff Younker on 4/16/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

import Foundation
import AudioToolbox

func setMixerBusCount(mixer: AudioUnit, busCount: UInt32) -> OSStatus {
    var bc = busCount
    let bcSize: UInt32 = 4
    return AudioUnitSetProperty(mixer,
        OSType(kAudioUnitProperty_BusCount),
        OSType(kAudioUnitScope_Input),
        0,
        &bc,
        bcSize)
}

func setReverb(mixer: AudioUnit, busIndex: UInt32) -> OSStatus {
    var renderFlags3d: UInt32 = 0
    var renderFlagsSize: UInt32 = 4
    if let err = checked(AudioUnitGetProperty(mixer,
        OSType(kAudioUnitProperty_3DMixerRenderingFlags),
        OSType(kAudioUnitScope_Input),
        busIndex,
        &renderFlags3d,
        &renderFlagsSize)) {
            return err
    }
    
    renderFlags3d |= UInt32(k3DMixerRenderingFlags_DistanceDiffusion)
    return AudioUnitSetProperty(mixer,
        OSType(kAudioUnitProperty_3DMixerRenderingFlags),
        OSType(kAudioUnitScope_Input),
        busIndex,
        &renderFlags3d,
        renderFlagsSize)
}

func setObjectCoordinates(mixer: AudioUnit, busIndex: UInt32,
    azimuth: Float32, distance: Float32) -> OSStatus {
        if let err = checked(AudioUnitSetParameter(mixer, OSType(k3DMixerParam_Azimuth), OSType(kAudioUnitScope_Input), busIndex, azimuth, 0)) {
            return err
        }
        return AudioUnitSetParameter(mixer, OSType(k3DMixerParam_Distance), OSType(kAudioUnitScope_Input), busIndex, distance, 0)
}
