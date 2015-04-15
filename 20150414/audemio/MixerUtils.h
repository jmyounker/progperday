//
//  MixerUtils.h
//  audemio
//
//  Created by Jeff Younker on 4/14/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#ifndef audemio_MixerUtils_h
#define audemio_MixerUtils_h

#import <AVFoundation/AVFoundation.h>
#import <Foundation/Foundation.h>

//
//  MixerUtils.cpp
//  audemio
//
//  Created by Jeff Younker on 4/14/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

OSStatus setMixerBusCount(AudioUnit mixer, UInt32 inBusCount);
OSStatus setReverb(AudioUnit mixer, UInt32 busIndex);
OSStatus setObjectCoordinates(AudioUnit mixer, UInt32 mixerBus, Float32 inAzimuth,
                              Float32 inDistance);
OSStatus configureGraphForChannelLayout(AudioUnit mixer, AudioUnit output);

#endif
