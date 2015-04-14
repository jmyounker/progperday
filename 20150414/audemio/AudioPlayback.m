//
//  AudioPlayback.m
//  audemio
//
//  Created by Jeff Younker on 4/13/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#include <unistd.h>

#import "MixerUtils.h"
#import "Sound.h"

#import "AudioPlayback.h"

#define kNumberPlaybackBuffers	3
#define kBuffTime 0.5
#define kTotalBuffDrainTime kNumberPlaybackBuffers * kBuffTime

@implementation AudioPlayback {
    Sound* source; // reference to your output file
    AUGraph graph; // audio graph that plays events
    AudioUnit mixerAu; // audio unit for 3d mixer
}

- (OSStatus) playSound: (NSString*) file {
    self->source = [[Sound alloc] initWithAudioFile:file busIndex:0];

    OSStatus res = [self createAuGraph];
    if (res != noErr) {
        return res;
    }

    res = [self->source open];
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphStart(self->graph);
    if (res != noErr) {
        return res;
    }
    NSLog(@"playing sound...");
    
    for (Float32 i=0; i<21*4; i++) {
        setObjectCoordinates(self->mixerAu, 0, 0, i*0.25);
        usleep(0.25 * 1000.0 * 1000.0);
    }

    NSLog(@"playback complete...");
    
cleanup:
    [self cleanup];
    
    return noErr;
}

- (void) cleanup {
    AUGraphStop(self->graph);
    AUGraphUninitialize(self->graph);
    AUGraphClose(self->graph);
    [self->source cleanup];
}

- (OSStatus) createAuGraph {
    OSStatus res = NewAUGraph(&self->graph);
    if (res != noErr) {
        return res;
    }

    AudioComponentDescription outputCd = {0};
    outputCd.componentType = kAudioUnitType_Output;
    outputCd.componentSubType = kAudioUnitSubType_DefaultOutput;
    outputCd.componentManufacturer = kAudioUnitManufacturer_Apple;
    
    AUNode outputNode;
    res = AUGraphAddNode(self->graph, &outputCd, &outputNode);

    if (res != noErr) {
        return res;
    }

    AudioComponentDescription fileCd = {0};
    fileCd.componentType = kAudioUnitType_Generator;
    fileCd.componentSubType = kAudioUnitSubType_AudioFilePlayer;
    fileCd.componentManufacturer = kAudioUnitManufacturer_Apple;

    AUNode fileNode;
    res = AUGraphAddNode(self->graph, &fileCd, &fileNode);
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
    res = AUGraphAddNode(self->graph, &mixerCd, &mixerNode);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphOpen(self->graph);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphNodeInfo(self->graph, fileNode, NULL, [self->source audioUnit]);
    if (res != noErr) {
        return res;
    }

    res = AUGraphNodeInfo(self->graph, mixerNode, NULL, &self->mixerAu);
    if (res != noErr) {
        return res;
    }

    res = setMixerBusCount(self->mixerAu, 64);
    if (res != noErr) {
        return res;
    }

    UInt32 reverbSetting = 1; // turn it on;
    res = AudioUnitSetProperty(self->mixerAu,
                                  kAudioUnitProperty_UsesInternalReverb,
                                  kAudioUnitScope_Global,
                                  0,
                                  &reverbSetting,
                                  sizeof(reverbSetting));
    
    res = AUGraphConnectNodeInput(self->graph, fileNode, 0, mixerNode, 0);
    if (res != noErr) {
        return res;
    }

    res = AUGraphConnectNodeInput(self->graph, mixerNode, 0, outputNode, 0);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphInitialize(self->graph);
    if (res != noErr) {
        return res;
    }

    setReverb(self->mixerAu, 0);
    return noErr;
}


@end