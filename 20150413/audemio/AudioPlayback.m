//
//  AudioPlayback.m
//  audemio
//
//  Created by Jeff Younker on 4/13/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#include <unistd.h>

#import "AudioPlayback.h"


#define kNumberPlaybackBuffers	3
#define kBuffTime 0.5
#define kTotalBuffDrainTime kNumberPlaybackBuffers * kBuffTime

typedef struct AudPlayer {
    AudioFileID					audioFile; // reference to your output file
    AudioStreamBasicDescription format; // data format
    AUGraph                     graph; // audio graph that plays events
    AudioUnit                   fileAu; // audio unit that reads from file

    SInt64						packetPosition; // current packet index in output file
    UInt32						numPacketsToRead; // number of packets to read from file
    AudioStreamPacketDescription *packetDescs; // array of packet descriptions for read buffer
    Boolean						isDone; // playback has completed
} AudPlayer;

@implementation AudioPlayback


- (OSStatus) playSound: (NSString*) file {
    // Construct URL to sound file
    NSURL *soundUrl = [NSURL fileURLWithPath:file];
    
    // Determine how long the sound is
    AVURLAsset* audioAsset = [AVURLAsset URLAssetWithURL:soundUrl options:nil];
    CMTime duration = audioAsset.duration;
    float durationSeconds = CMTimeGetSeconds(duration);                                ;
    
    NSNumberFormatter *formatter = [[NSNumberFormatter alloc] init];
    formatter.roundingIncrement = [NSNumber numberWithDouble:0.01];
    
    NSLog(@"playing sound for %@ seconds",
          [formatter stringFromNumber:[NSNumber numberWithDouble: durationSeconds]]);

    AudPlayer player = {0};

    OSStatus res = AudioFileOpenURL(CFBridgingRetain(soundUrl), kAudioFileReadPermission, 0, &player.audioFile);
    if (res != noErr) {
        return res;
    }

    createAuGraph(&player);
    res = prepareInput(&player);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphStart(player.graph);
    if (res != noErr) {
        return res;
    }
    NSLog(@"playing sound...");
    
    sleep(durationSeconds);

    NSLog(@"playback complete...");
    
cleanup:
    AUGraphStop(player.graph);
    AUGraphUninitialize(player.graph);
    AUGraphClose(player.graph);
    AudioFileClose(player.audioFile);

    return noErr;
}


static OSStatus createAuGraph(AudPlayer* player) {
    OSStatus res = NewAUGraph(&player->graph);
    if (res != noErr) {
        return res;
    }

    AudioComponentDescription outputCd = {0};
    outputCd.componentType = kAudioUnitType_Output;
    outputCd.componentSubType = kAudioUnitSubType_DefaultOutput;
    outputCd.componentManufacturer = kAudioUnitManufacturer_Apple;
    
    AUNode outputNode;
    res = AUGraphAddNode(player->graph, &outputCd, &outputNode);

    if (res != noErr) {
        return res;
    }

    AudioComponentDescription fileCd = {0};
    fileCd.componentType = kAudioUnitType_Generator;
    fileCd.componentSubType = kAudioUnitSubType_AudioFilePlayer;
    fileCd.componentManufacturer = kAudioUnitManufacturer_Apple;

    AUNode fileNode;
    res = AUGraphAddNode(player->graph, &fileCd, &fileNode);
    if (res != noErr) {
        return res;
    }

    res = AUGraphOpen(player->graph);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphNodeInfo(player->graph,
                          fileNode,
                          NULL,
                          &player->fileAu);
    if (res != noErr) {
        return res;
    }

    res = AUGraphConnectNodeInput(player->graph,
                                  fileNode,
                                  0,
                                  outputNode,
                                  0);

    if (res != noErr) {
        return res;
    }

    res = AUGraphInitialize(player->graph);
    if (res != noErr) {
        return res;
    }
    
    return noErr;
}

static OSStatus prepareInput(AudPlayer* player) {
    UInt32 formatPropSize = sizeof(player->format);
    OSStatus res = AudioFileGetProperty(player->audioFile,
                                        kAudioFilePropertyDataFormat,
                                        &formatPropSize,
                                        &player->format);
    if (res != noErr) {
        return res;
    }

    
    res = AudioUnitSetProperty(player->fileAu,
                               kAudioUnitProperty_ScheduledFileIDs,
                               kAudioUnitScope_Global,
                               0,
                               &player->audioFile,
                               sizeof(player->audioFile));
    if (res != noErr) {
        return res;
    }

    
    UInt64 numPackets;
    UInt32 numPacketsPropSize = sizeof(numPackets);
    res = AudioFileGetProperty(player->audioFile,
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
    region.mAudioFile = player->audioFile;
    region.mLoopCount = 1;
    region.mStartFrame = 0;
    region.mFramesToPlay = (UInt32)numPackets * player->format.mFramesPerPacket;

    res = AudioUnitSetProperty(player->fileAu,
                               kAudioUnitProperty_ScheduledFileRegion,
                               kAudioUnitScope_Global,
                               0,
                               &region,
                               sizeof(region));
    if (res != noErr) {
        return res;
    }
    
    AudioTimeStamp startTime;
    memset (&startTime, 0, sizeof(startTime));
    startTime.mFlags = kAudioTimeStampSampleTimeValid;
    startTime.mSampleTime = -1;

    res = AudioUnitSetProperty(player->fileAu,
                               kAudioUnitProperty_ScheduleStartTimeStamp,
                               kAudioUnitScope_Global,
                               0,
                               &startTime,
                               sizeof(startTime));
    if (res != noErr) {
        return res;
    }
    
    return noErr;
}

@end