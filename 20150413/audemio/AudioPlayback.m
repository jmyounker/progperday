//
//  AudioPlayback.m
//  audemio
//
//  Created by Jeff Younker on 4/2/15.
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
    
    UInt32 propSize = sizeof(player.format);
    res = AudioFileGetProperty(player.audioFile,
                               kAudioFilePropertyDataFormat,
                               &propSize,
                               &player.format);
    if (res != noErr) {
        return res;
    }
    
    AudioQueueRef queue;
    res = AudioQueueNewOutput(&player.format,
                              pullAudioCallback,
                              &player,
                              NULL,
                              NULL,
                              0,
                              &queue);
    

    UInt32 buffSize; // in bytes
    res = sizesFromTime(player.audioFile, player.format, kBuffTime, &buffSize, &player.numPacketsToRead);
    if (res != noErr) {
        return res;
    }
    
    if (isVBR(player.format)) {
        player.packetDescs = (AudioStreamPacketDescription*)malloc(sizeof(AudioStreamPacketDescription) * player.numPacketsToRead);
    } else {
        player.packetDescs = NULL;
    }
    
    res = copyEncoderCookieToQueue(player.audioFile, queue);
    if (res != noErr) {
        return res;
    }

    AudioQueueBufferRef buffs[kNumberPlaybackBuffers];
    player.isDone = false;
    player.packetPosition = 0;
    int i;
    for (i = 0; i < kNumberPlaybackBuffers; i++) {
        res = AudioQueueAllocateBuffer(queue, buffSize, &buffs[i]);
        if (res != noErr) {
            return res;
        }
        
        pullAudioCallback(&player, queue, buffs[i]);

        if (player.isDone) {
            break;
        }
    }

    res = AudioQueueStart(queue, NULL);
    if (res != noErr) {
        return 1;
    }
    
    printf("playing sound.\n");
    while (!player.isDone) {
        CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0.25, false);
    }

    printf("queuing done.\n");

    CFRunLoopRunInMode(kCFRunLoopDefaultMode, kTotalBuffDrainTime, false);

    res = AudioQueueStop(queue, true);
    if (res != noErr) {
        return res;
    }

    printf("playback done.\n");
    

cleanup:
    AudioQueueDispose(queue, true);
    AudioFileClose(player.audioFile);

    return noErr;
}

static bool isVBR(AudioStreamBasicDescription format) {
    return (format.mBytesPerPacket == 0 || format.mFramesPerPacket == 0);
}


static OSStatus sizesFromTime(AudioFileID audioFile, AudioStreamBasicDescription desc,
                                Float64 time, UInt32* bufSize, UInt32* numPackets) {
    OSStatus err;
    
    UInt32 maxPacketSize;
    UInt32 propSize = sizeof(maxPacketSize);
    OSStatus res = AudioFileGetProperty(audioFile,
                                        kAudioFilePropertyPacketSizeUpperBound,
                                        &propSize,
                                        &maxPacketSize);
    if (res != noErr) {
        return err;
    }
    
    static const int maxBufSize = 0x10000; // limit size to 64K
    static const int minBufSize = 0x4000; // limit size to 16K

    UInt32 defaultBufSize = MAX(maxBufSize, maxPacketSize);

    if (desc.mFramesPerPacket) {
        Float64 numPacketsPerTime = desc.mSampleRate / desc.mFramesPerPacket;
        Float64 numPackets = numPacketsPerTime * time;
        *bufSize = numPackets * maxPacketSize;
    } else {
        *bufSize = defaultBufSize;
    }

    *bufSize = MIN(*bufSize, defaultBufSize);
    *bufSize = MAX(*bufSize, minBufSize);

    *numPackets = *bufSize / maxPacketSize;
    
    return noErr;
}

static OSStatus copyEncoderCookieToQueue(AudioFileID audioFile, AudioQueueRef queue) {
    UInt32 propSize;
    OSStatus res = AudioFileGetPropertyInfo(audioFile,
                                            kAudioFilePropertyMagicCookieData,
                                            &propSize,
                                            NULL);
    if (res == noErr && propSize > 0) {
        Byte* cookie = (UInt8*)malloc(sizeof(UInt8) * propSize);
        res = AudioFileGetProperty(audioFile, kAudioFilePropertyMagicCookieData, &propSize, cookie);
        if (res != noErr) {
            free(cookie);
            return res;
        }

        res = AudioQueueSetProperty(queue, kAudioQueueProperty_MagicCookie, cookie, propSize);
        if (res != noErr) {
            free(cookie);
            return res;
        }
        free(cookie);
    }
    return noErr;
}

static void pullAudioCallback(void *audPlayer, AudioQueueRef queue, AudioQueueBufferRef buff) {
    AudPlayer *ap = (AudPlayer*)audPlayer;
    if (ap->isDone) {
        return;
    }
    
    UInt32 numBytes;
    UInt32 numPackets = ap->numPacketsToRead;
    OSStatus res = AudioFileReadPackets(ap->audioFile,
                                        false,
                                        &numBytes,
                                        ap->packetDescs,
                                        ap->packetPosition,
                                        &numPackets,
                                        buff->mAudioData);
    if (res != noErr) {
        NSLog(@"READING PACKETS FAILURE");
        exit(2);
    }
    
    if (numPackets > 0) {
        buff->mAudioDataByteSize = numBytes;
        res = AudioQueueEnqueueBuffer(queue,
                                buff,
                                (ap->packetDescs ? numPackets : 0),
                                ap->packetDescs);
        if (res != noErr) {
            NSLog(@"AUDIO ENQUEUING FAILURE");
            exit(2);
        }
        
        ap->packetPosition += numPackets;
    } else {
        res = AudioQueueStop(queue, false);
        if (res != noErr) {
            NSLog(@"AUDIO STOP FAILURE");
            exit(2);
        }
        
        ap->isDone = true;
    }
}

@end