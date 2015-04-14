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

typedef struct AudSource {
    AudioFileID	file; // input file
    AudioStreamBasicDescription format; // data format
    AudioUnit au; // audio unit
    UInt32 busIndex; // busIndex
    float duration; // audio length in seconds

} AudSource;

typedef struct AudPlayer {
    AudSource*					source; // reference to your output file
    AUGraph                     graph; // audio graph that plays events
    AudioUnit                   mixerAu; // audio unit for 3d mixer
} AudPlayer;

@implementation AudioPlayback



static OSStatus initAudioSource(NSString* file, UInt32 busIndex, AudSource* audSource)  {
    audSource->busIndex = busIndex;
    
    NSURL* soundUrl = [NSURL fileURLWithPath:file];
    AVURLAsset* audioAsset = [AVURLAsset URLAssetWithURL:soundUrl options:nil];
    audSource->duration = CMTimeGetSeconds(audioAsset.duration);

    OSStatus res = AudioFileOpenURL(CFBridgingRetain(soundUrl), kAudioFileReadPermission, 0, &audSource->file);
    if (res != noErr) {
        return res;
    }

    return noErr;
}

- (OSStatus) playSound: (NSString*) file {
    AudPlayer player = {0};
    AudSource source = {0};
    player.source = &source;
    
    OSStatus res = initAudioSource(file, 0, player.source);
    if (res != noErr) {
        return res;
    }

    createAuGraph(&player);
    res = prepareInput(player.source);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphStart(player.graph);
    if (res != noErr) {
        return res;
    }
    NSLog(@"playing sound...");
    
    for (Float32 i=0; i<21*4; i++) {
        setObjectCoordinates(&player, 0, 0, i*0.25);
        usleep(0.25 * 1000.0 * 1000.0);
    }

    NSLog(@"playback complete...");
    
cleanup:
    AUGraphStop(player.graph);
    AUGraphUninitialize(player.graph);
    AUGraphClose(player.graph);
    AudioFileClose(player.source->file);

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
    
    AudioComponentDescription mixerCd = {0};
    mixerCd.componentType = kAudioUnitType_Mixer;
    mixerCd.componentSubType = kAudioUnitSubType_3DMixer;
    mixerCd.componentManufacturer = kAudioUnitManufacturer_Apple;
    mixerCd.componentFlags = 0;
    mixerCd.componentFlagsMask = 0;

    AUNode mixerNode;
    res = AUGraphAddNode(player->graph, &mixerCd, &mixerNode);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphOpen(player->graph);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphNodeInfo(player->graph, fileNode, NULL, &player->source->au);
    if (res != noErr) {
        return res;
    }

    res = AUGraphNodeInfo(player->graph, mixerNode, NULL, &player->mixerAu);
    if (res != noErr) {
        return res;
    }

    res = setMixerBusCount(player, 64);
    if (res != noErr) {
        return res;
    }

    UInt32 reverbSetting = 1; // turn it on;
    res = AudioUnitSetProperty(player->mixerAu,
                                  kAudioUnitProperty_UsesInternalReverb,
                                  kAudioUnitScope_Global,
                                  0,
                                  &reverbSetting,
                                  sizeof(reverbSetting));
    
    res = AUGraphConnectNodeInput(player->graph, fileNode, 0, mixerNode, 0);
    if (res != noErr) {
        return res;
    }

    res = AUGraphConnectNodeInput(player->graph, mixerNode, 0, outputNode, 0);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphInitialize(player->graph);
    if (res != noErr) {
        return res;
    }

    setReverb(player, 0);
    return noErr;
}

static OSStatus setMixerBusCount(AudPlayer* player, UInt32 inBusCount) {
    UInt32 busCount = inBusCount;
    UInt32 size = sizeof(busCount);
    return (AudioUnitSetProperty (player->mixerAu,
                                  kAudioUnitProperty_BusCount,
                                  kAudioUnitScope_Input,
                                  0,
                                  &busCount,
                                  size));
    
}

static OSStatus setReverb(AudPlayer* player, UInt32 busIndex) {
    UInt32    renderFlags3d;
    UInt32    outSize = sizeof(renderFlags3d);

    // get the current render flags for this bus
    OSStatus res = AudioUnitGetProperty (player->mixerAu,
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
    return AudioUnitSetProperty(player->mixerAu,
                              kAudioUnitProperty_3DMixerRenderingFlags,
                              kAudioUnitScope_Input, busIndex,
                              &renderFlags3d, 
                              sizeof(renderFlags3d));

}

static OSStatus setObjectCoordinates(AudPlayer* player, UInt32 mixerBus, Float32 inAzimuth,
                                     Float32 inDistance) {
    OSStatus res = AudioUnitSetParameter(player->mixerAu,
                          k3DMixerParam_Azimuth,
                          kAudioUnitScope_Input,
                          mixerBus,
                          inAzimuth,
                          0);
    if (res != noErr) {
        return res;
    }
    
    return AudioUnitSetParameter(player->mixerAu,
                          k3DMixerParam_Distance,
                          kAudioUnitScope_Input,
                          mixerBus,
                          inDistance,
                          0);
}

static OSStatus prepareInput(AudSource* source) {
    UInt32 formatPropSize = sizeof(source->format);
    OSStatus res = AudioFileGetProperty(source->file,
                                        kAudioFilePropertyDataFormat,
                                        &formatPropSize,
                                        &source->format);
    if (res != noErr) {
        return res;
    }

    
    res = AudioUnitSetProperty(source->au,
                               kAudioUnitProperty_ScheduledFileIDs,
                               kAudioUnitScope_Global,
                               0,
                               &source->file,
                               sizeof(source->file));
    if (res != noErr) {
        return res;
    }

    
    UInt64 numPackets;
    UInt32 numPacketsPropSize = sizeof(numPackets);
    res = AudioFileGetProperty(source->file,
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
    region.mAudioFile = source->file;
    region.mLoopCount = 1;
    region.mStartFrame = 0;
    region.mFramesToPlay = (UInt32)numPackets * source->format.mFramesPerPacket;

    res = AudioUnitSetProperty(source->au,
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

    res = AudioUnitSetProperty(source->au,
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