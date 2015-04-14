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
    AUGraph graph; // audio graph that plays events
    AudioUnit mixerAu; // audio unit for 3d mixer
}

- (OSStatus) playSound: (NSArray*) files {
    NSMutableArray* srcs = [NSMutableArray arrayWithObjects: nil];
    for (int i=0; i<[files count]; i++) {
        [srcs addObject: [[Sound alloc] initWithAudioFile:files[i] busIndex:i]];
    };
    self.sources = srcs;
         
    OSStatus res = [self createAuGraph];
    if (res != noErr) {
        return res;
    }

    for (Sound* sound in self.sources) {
        res = [sound open];
        if (res != noErr) {
            return res;
        }
    }
    
    res = AUGraphStart(self->graph);
    if (res != noErr) {
        return res;
    }
    NSLog(@"playing sound...");

    // For the hard-coded current values, make Rana clamitans fades into the distance behind/ahead,
    // while Rana sylvatica swirls from right to left.
    Float32 interval = 0.25;
    Float32 intervalCount = 21 / interval;
    Float32 azm2 = 90;
    for (Float32 i=0; i<intervalCount; i++) {
        setObjectCoordinates(self->mixerAu, 0, 0, i*interval/4);
        setObjectCoordinates(self->mixerAu, 1, azm2, 1);
        azm2 -= 180.0 / intervalCount;
        usleep(interval * 1000.0 * 1000.0);
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
    for (Sound* sound in self.sources) {
        [sound cleanup];
    }
}

- (OSStatus) createAuGraph {
    OSStatus res = NewAUGraph(&self->graph);
    if (res != noErr) {
        NSLog(@"error: could not create augraph");
        return res;
    }

    AudioComponentDescription outputCd = {0};
    outputCd.componentType = kAudioUnitType_Output;
    outputCd.componentSubType = kAudioUnitSubType_DefaultOutput;
    outputCd.componentManufacturer = kAudioUnitManufacturer_Apple;
    
    AUNode outputNode;
    res = AUGraphAddNode(self->graph, &outputCd, &outputNode);

    if (res != noErr) {
        NSLog(@"error: could not add output to augraph");
        return res;
    }
    
    AUNode soundNodes[[self.sources count]];
    for (int i=0; i<[self.sources count]; i++) {
        AudioComponentDescription fileCd = {0};
        fileCd.componentType = kAudioUnitType_Generator;
        fileCd.componentSubType = kAudioUnitSubType_AudioFilePlayer;
        fileCd.componentManufacturer = kAudioUnitManufacturer_Apple;
    
        res = AUGraphAddNode(self->graph, &fileCd, &soundNodes[i]);
        if (res != noErr) {
            NSLog(@"error: could not add sound to augraph");
            return res;
        }
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
        NSLog(@"error: could not add node to augraph");
        return res;
    }
    
    res = AUGraphOpen(self->graph);
    if (res != noErr) {
        NSLog(@"error: could not open augraph");
        return res;
    }
    
    for (int i=0; i<[self.sources count]; i++) {
        Sound* sound = self.sources[i];
        res = AUGraphNodeInfo(self->graph, soundNodes[i], NULL, [sound audioUnit]);
        if (res != noErr) {
            NSLog(@"error: could not get sound audio unit");
            return res;
        }
    }

    res = AUGraphNodeInfo(self->graph, mixerNode, NULL, &self->mixerAu);
    if (res != noErr) {
        NSLog(@"error: could not get mixer audio unit");
        return res;
    }

    res = setMixerBusCount(self->mixerAu, 64);
    if (res != noErr) {
        NSLog(@"error: could not set mixer bus count");
        return res;
    }

    UInt32 reverbSetting = 1; // turn it on;
    res = AudioUnitSetProperty(self->mixerAu,
                                  kAudioUnitProperty_UsesInternalReverb,
                                  kAudioUnitScope_Global,
                                  0,
                                  &reverbSetting,
                                  sizeof(reverbSetting));
    
    for (int i=0; i<[self.sources count]; i++) {
        Sound* sound = self.sources[i];
        res = AUGraphConnectNodeInput(self->graph, soundNodes[i], 0, mixerNode, sound.busIndex);
        if (res != noErr) {
            NSLog(@"error: could not wire sound into augraph");
            return res;
        }
    }
        
    res = AUGraphConnectNodeInput(self->graph, mixerNode, 0, outputNode, 0);
    if (res != noErr) {
        return res;
    }
    
    res = AUGraphInitialize(self->graph);
    if (res != noErr) {
        NSLog(@"error: could not initialize augraph");
        return res;
    }

    setReverb(self->mixerAu, 0);

    return noErr;
}


@end