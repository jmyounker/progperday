//
//  AudioPlayback.m
//  audemio
//
//  Created by Jeff Younker on 4/2/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#include <unistd.h>

#import "AudioPlayback.h"

@implementation AudioPlayback

- (void) playSound: (NSString*) file {
    // Construct URL to sound file
    NSURL *soundUrl = [NSURL fileURLWithPath:file];

    // Create audio player object and initialize with URL to sound
    AVAudioPlayer* player;
    NSError* err;
    player = [[AVAudioPlayer alloc] initWithContentsOfURL:soundUrl error:&err];
    if (err) {
        NSLog(@"failed to construct player: %@", [err localizedDescription]);
        return;
    }

    // Determine how long the sound is
    AVURLAsset* audioAsset = [AVURLAsset URLAssetWithURL:soundUrl options:nil];
    CMTime duration = audioAsset.duration;
    float durationSeconds = CMTimeGetSeconds(duration);                                ;

    NSNumberFormatter *formatter = [[NSNumberFormatter alloc] init];
    formatter.roundingIncrement = [NSNumber numberWithDouble:0.01];

    NSLog(@"playing sound for %@ seconds",
          [formatter stringFromNumber:[NSNumber numberWithDouble: durationSeconds]]);

    [player play];
    // If the program ends then the sound stops playing.
    sleep(durationSeconds);

    NSLog(@"sound played");
    
}

@end