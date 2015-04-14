//
//  AudSource.m
//  audemio
//
//  Created by Jeff Younker on 4/14/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#import <Foundation/Foundation.h>

#import "Sound.h"

@implementation Sound

- (id) initWithAudioFile:(NSString*)path busIndex:(UInt32)bus  {
    self = [super init];
    if (self) {
        self.filePath = path;
        self.busIndex = bus;
    }
    return self;
}

- (void) cleanup {
    AudioFileClose(self->file);
}

- (AudioUnit*) audioUnit {
    return &self->au;
}

- (OSStatus) open {
    NSURL* soundUrl = [NSURL fileURLWithPath:self.filePath];
    AVURLAsset* audioAsset = [AVURLAsset URLAssetWithURL:soundUrl options:nil];
    self.duration = CMTimeGetSeconds(audioAsset.duration);
    
    OSStatus res = AudioFileOpenURL(CFBridgingRetain(soundUrl), kAudioFileReadPermission, 0, &self->file);
    if (res != noErr) {
        NSLog(@"error: could not set open file");
        return res;
    }

    UInt32 formatPropSize = sizeof(self->format);
    res = AudioFileGetProperty(self->file,
                               kAudioFilePropertyDataFormat,
                               &formatPropSize,
                               &self->format);
    if (res != noErr) {
        NSLog(@"error: could not set get file properties");
        return res;
    }
    
    
    res = AudioUnitSetProperty(self->au,
                               kAudioUnitProperty_ScheduledFileIDs,
                               kAudioUnitScope_Global,
                               0,
                               &self->file,
                               sizeof(self->file));
    if (res != noErr) {
        NSLog(@"error: could not set file id");
        return res;
    }
    
    
    UInt64 numPackets;
    UInt32 numPacketsPropSize = sizeof(numPackets);
    res = AudioFileGetProperty(self->file,
                               kAudioFilePropertyAudioDataPacketCount,
                               &numPacketsPropSize,
                               &numPackets);
    
    // Tell the file player AU to play the entire file
    ScheduledAudioFileRegion region;
    memset (&region.mTimeStamp, 0, sizeof(region.mTimeStamp));
    region.mTimeStamp.mFlags = kAudioTimeStampSampleTimeValid;
    region.mTimeStamp.mSampleTime = 0;
    region.mCompletionProc = NULL;
    region.mCompletionProcUserData = NULL;
    region.mAudioFile = self->file;
    region.mLoopCount = 1;
    region.mStartFrame = 0;
    region.mFramesToPlay = (UInt32)numPackets * self->format.mFramesPerPacket;
    
    res = AudioUnitSetProperty(self->au,
                               kAudioUnitProperty_ScheduledFileRegion,
                               kAudioUnitScope_Global,
                               0,
                               &region,
                               sizeof(region));
    if (res != noErr) {
        NSLog(@"error: could not set playback region");
        return res;
    }
    
    AudioTimeStamp startTime;
    memset (&startTime, 0, sizeof(startTime));
    startTime.mFlags = kAudioTimeStampSampleTimeValid;
    startTime.mSampleTime = -1;
    
    res = AudioUnitSetProperty(self->au,
                               kAudioUnitProperty_ScheduleStartTimeStamp,
                               kAudioUnitScope_Global,
                               0,
                               &startTime,
                               sizeof(startTime));
    if (res != noErr) {
        NSLog(@"error: could not set start time");
        return res;
    }
    
    return noErr;
}

@end