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
        NSArray* arguments = [[NSProcessInfo processInfo] arguments];
        if (arguments.count != 2) {
            NSLog(@"expected only one argument");
            return 127;
        }
        NSString* soundFile = arguments[1];
        AudioPlayback* ap = [[AudioPlayback alloc] init];
        [ap playSound: soundFile];
    }
    return 0;
}
