//
//  AudioPlayback.h
//  audemio
//
//  Created by Jeff Younker on 4/2/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#ifndef audemio_AudioPlayback_h
#define audemio_AudioPlayback_h

#import <AudioToolbox/AudioToolbox.h>
#import <AVFoundation/AVFoundation.h>

#define NUM_FILES  2

@interface AudioPlayback : NSObject {
}

- (OSStatus) playSound: (NSString*) file;

@end

#endif
