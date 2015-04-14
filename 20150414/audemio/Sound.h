//
//  AudSource.h
//  audemio
//
//  Created by Jeff Younker on 4/14/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#ifndef audemio_AudSource_h
#define audemio_AudSource_h

#import <AudioToolbox/AudioToolbox.h>
#import <AVFoundation/AVFoundation.h>

#define NUM_FILES  2

@interface Sound : NSObject {
    NSString* filePath; // input file path
    AudioFileID file; // input file
    AudioStreamBasicDescription format; // data format
    AudioUnit au; // audio unit
    UInt32 busIndex; // busIndex
    float duration; // audio length in seconds
}

- (id) initWithAudioFile:(NSString*)path busIndex:(UInt32)bus;

- (OSStatus) open;

- (AudioUnit*) audioUnit;

- (void) cleanup;

@end


#endif
