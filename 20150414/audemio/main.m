//
//  main.m
//  audemio
//
//  Created by Jeff Younker on 4/2/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "AudioPlayback.h"

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSString* frog1 = @"/Users/jeff/repos/progperday/20150414/audemio/Rana_clamitans.mp3";
        NSString* frog2 = @"/Users/jeff/repos/progperday/20150414/audemio/rana_sylvatica_chorus.aiff";
        AudioPlayback* ap = [[AudioPlayback alloc] init];
        OSStatus res = [ap playSound: @[frog1, frog2]];
        if (res != noErr) {
            NSLog(@"Failed!");
        }
    }
    return 0;
}
