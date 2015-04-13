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
//        NSArray* arguments = [[NSProcessInfo processInfo] arguments];
//        if (arguments.count != 2) {
//            NSLog(@"expected only one argument");
//            return 127;
//        }
//        NSString* soundFile = arguments[1];
        NSString* soundFile = @"/Users/jeff/repos/progperday/20150413/audemio/Rana_clamitans.mp3";
        AudioPlayback* ap = [[AudioPlayback alloc] init];
        OSStatus res = [ap playSound: soundFile];
        if (res != noErr) {
            NSLog(@"Failed!");
        }
    }
    return 0;
}
