//
//  MixerUtils.cpp
//  audemio
//
//  Created by Jeff Younker on 4/14/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//


#include "MixerUtils.h"

const Float64 kDefaultSampleRate = 44100.0;

// #import "CoreAudio/PublicUtility/CAStreamBasicDescription.h"

OSStatus setMixerBusCount(AudioUnit mixer, UInt32 inBusCount) {
    UInt32 busCount = inBusCount;
    UInt32 size = sizeof(busCount);
    return (AudioUnitSetProperty (mixer,
                                  kAudioUnitProperty_BusCount,
                                  kAudioUnitScope_Input,
                                  0,
                                  &busCount,
                                  size));
    
}

OSStatus setReverb(AudioUnit mixer, UInt32 busIndex) {
    UInt32    renderFlags3d;
    UInt32    outSize = sizeof(renderFlags3d);
    
    // get the current render flags for this bus
    OSStatus res = AudioUnitGetProperty (mixer,
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
    return AudioUnitSetProperty(mixer,
                                kAudioUnitProperty_3DMixerRenderingFlags,
                                kAudioUnitScope_Input, busIndex,
                                &renderFlags3d,
                                sizeof(renderFlags3d));
    
}

OSStatus setObjectCoordinates(AudioUnit mixer, UInt32 mixerBus, Float32 inAzimuth,
                                     Float32 inDistance) {
    OSStatus res = AudioUnitSetParameter(mixer,
                                         k3DMixerParam_Azimuth,
                                         kAudioUnitScope_Input,
                                         mixerBus,
                                         inAzimuth,
                                         0);
    if (res != noErr) {
        return res;
    }
    
    return AudioUnitSetParameter(mixer,
                                 k3DMixerParam_Distance,
                                 kAudioUnitScope_Input,
                                 mixerBus,
                                 inDistance,
                                 0);
}

OSStatus getDesiredRenderChannelCount(AudioUnit outputAu, UInt32* numChannels) {
    // get the HAL device id form the output AU
    AudioDeviceID  deviceID;
    UInt32 returnValue = 2; // return stereo by default
    
    UInt32 deviceSize =  sizeof(deviceID);
    
    //get the current device
    OSStatus res = AudioUnitGetProperty(outputAu,
                                        kAudioOutputUnitProperty_CurrentDevice,
                                        kAudioUnitScope_Output,
                                        1,
                                        &deviceID,
                                        &deviceSize);
    if (res != noErr) {
        return res;
    }
    
    //Get the users speaker configuration
    const AudioObjectPropertyAddress channelLayoutAddr = {
        kAudioDevicePropertyPreferredChannelLayout,
        kAudioObjectPropertyScopeGlobal,
        0 };

    UInt32 layoutSize;
    res = AudioObjectGetPropertyDataSize(deviceID,
                                                  &channelLayoutAddr,
                                                  0,
                                                  NULL,
                                                  &layoutSize);
    
    if (res != noErr) {
        *numChannels = 2;
        return noErr;
    }
    
    AudioChannelLayout *layout = NULL;
    layout = (AudioChannelLayout *) calloc(1, layoutSize);

    
    if (layout != NULL) {
        res = AudioObjectGetPropertyData(deviceID,
                                         &channelLayoutAddr,
                                         0,
                                         NULL,
                                         &layoutSize,
                                         layout);
        if (res != noErr) {
            return res;
        }
        if (layout->mChannelLayoutTag == kAudioChannelLayoutTag_UseChannelDescriptions) {
            // no channel layout tag is returned,
            //so walk through the channel descriptions and count
            // the channels that are associated with a speaker
            if (layout->mNumberChannelDescriptions == 2) {
                returnValue = 2; // there is no channel info for stereo
            } else {
                returnValue = 0;
                for (UInt32 i = 0; i < layout->mNumberChannelDescriptions; i++) {
                    if (layout->mChannelDescriptions[i].mChannelLabel !=
                        kAudioChannelLabel_Unknown) {
                        returnValue++;
                    }
                }
            }
        } else {
            switch (layout->mChannelLayoutTag) {
                case kAudioChannelLayoutTag_AudioUnit_5_0:
                case kAudioChannelLayoutTag_AudioUnit_5_1:
                case kAudioChannelLayoutTag_AudioUnit_6:
                    returnValue = 5;
                    break;
                case kAudioChannelLayoutTag_AudioUnit_4:
                    returnValue = 4;
                    break;
                default:
                    returnValue = 2;
            }
        }
        
        free(layout);
    }

    *numChannels = returnValue;
    return noErr;
}

void setCanonicalAudioStreamBasic(UInt32 nChannels, bool interleaved, AudioStreamBasicDescription* asbd) {
    // note: leaves sample rate untouched
    asbd->mFormatID = kAudioFormatLinearPCM;
    asbd->mFormatFlags = kAudioFormatFlagsNativeFloatPacked | kAudioFormatFlagIsNonInterleaved;
    UInt32 sampleSize = sizeof(float);
    asbd->mBitsPerChannel = 8 * sampleSize;
    asbd->mChannelsPerFrame = nChannels;
    asbd->mFramesPerPacket = 1;
    asbd->mSampleRate = kDefaultSampleRate;
    if (interleaved)
        asbd->mBytesPerPacket = asbd->mBytesPerFrame = nChannels * sampleSize;
    else {
        asbd->mBytesPerPacket = asbd->mBytesPerFrame = sampleSize;
        asbd->mFormatFlags |= kAudioFormatFlagIsNonInterleaved;
    }
}

OSStatus configureGraphForChannelLayout(AudioUnit mixer, AudioUnit output) {
    
    // get the channel count that should be set
    // for the mixer's output stream format
    UInt32 mixerChannelCount;
    OSStatus res = getDesiredRenderChannelCount(output, &mixerChannelCount);
    if (res != noErr) {
        return res;
    }
    
    // set the stream format
    AudioStreamBasicDescription format;
    UInt32 outSize = sizeof(format);
    res = AudioUnitGetProperty(output,
                               kAudioUnitProperty_StreamFormat,
                               kAudioUnitScope_Output,
                               0,
                               &format,
                               &outSize);
    if (res != noErr) {
        return res;
    }
    
    // not interleaved
    setCanonicalAudioStreamBasic(mixerChannelCount, false, &format);
    outSize = sizeof(format);
    res = AudioUnitSetProperty (output,
                                kAudioUnitProperty_StreamFormat,
                                kAudioUnitScope_Input,
                                0,
                                &format,
                                outSize);
    if (res != noErr) {
        return res;
    }
    
    return AudioUnitSetProperty (mixer,
                                 kAudioUnitProperty_StreamFormat,
                                 kAudioUnitScope_Output,
                                 0,
                                 &format,
                                 outSize);
}

