//
//  MixerUtils.cpp
//  audemio
//
//  Created by Jeff Younker on 4/14/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#include "MixerUtils.h"

OSStatus setMixerBusCount(AudioUnit mixer, UInt32 inBusCount) {
    UInt32 busCount = inBusCount;
    UInt32 size = sizeof(busCount);
    return (AudioUnitSetProperty (mixer,
                                  kAudioUnitProperty_BusCount,
                                  kAudioUnitScope_Input,
                                  0,
                                  &busCount,
                                  size));
    
}

OSStatus setReverb(AudioUnit mixer, UInt32 busIndex) {
    UInt32    renderFlags3d;
    UInt32    outSize = sizeof(renderFlags3d);
    
    // get the current render flags for this bus
    OSStatus res = AudioUnitGetProperty (mixer,
                                         kAudioUnitProperty_3DMixerRenderingFlags,
                                         kAudioUnitScope_Input,
                                         busIndex,
                                         &renderFlags3d,
                                         &outSize);
    if (res != noErr) {
        return res;
    }
    
    // turn on this render flag and then set the bus
    renderFlags3d |= k3DMixerRenderingFlags_DistanceDiffusion;
    return AudioUnitSetProperty(mixer,
                                kAudioUnitProperty_3DMixerRenderingFlags,
                                kAudioUnitScope_Input, busIndex,
                                &renderFlags3d,
                                sizeof(renderFlags3d));
    
}

OSStatus setObjectCoordinates(AudioUnit mixer, UInt32 mixerBus, Float32 inAzimuth,
                                     Float32 inDistance) {
    OSStatus res = AudioUnitSetParameter(mixer,
                                         k3DMixerParam_Azimuth,
                                         kAudioUnitScope_Input,
                                         mixerBus,
                                         inAzimuth,
                                         0);
    if (res != noErr) {
        return res;
    }
    
    return AudioUnitSetParameter(mixer,
                                 k3DMixerParam_Distance,
                                 kAudioUnitScope_Input,
                                 mixerBus,
                                 inDistance,
                                 0);
}
