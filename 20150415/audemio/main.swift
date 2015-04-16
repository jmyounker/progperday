//
//  main.swift
//  audemio
//
//  Created by Jeff Younker on 4/15/15.
//  Copyright (c) 2015 Jeff Younker. All rights reserved.
//

import Foundation
import AudioToolbox
import AVFoundation

func main() -> Int {
    let frog1 = "/Users/jeff/repos/progperday/20150414/audemio/Rana_clamitans.mp3"
    let frog2 = "/Users/jeff/repos/progperday/20150414/audemio/rana_sylvatica_chorus.aiff"
    var ap = AudioPlayback(soundFiles: [frog1, frog2])
    ap.playSound()
    return 0;
}

func checked(status: OSStatus) -> OSStatus? {
    if status == noErr {
        return nil
    }
    return status
}

class AudioPlayback {
    var graph = AUGraph()
    let sounds: Array<Sound>
    var mixerAu = AudioUnit()
    var outputAu = AudioUnit()

    init(soundFiles: Array<String>) {
        var sounds: Array<Sound> = [];
        for i in 0...soundFiles.count - 1 {
            sounds.append(Sound(fileName: soundFiles[i], busIndex: UInt32(i)))
        }
        self.sounds = sounds
    }

    
    func playSound() -> OSStatus {
        if let err = checked(self.createAuGraph()) {
            return err
        }

        if let err = checked(self.createAuGraph()) {
            return err
        }
        
        for sound in self.sounds {
            if let err = checked(sound.open()) {
                return err
            }
        }
        
        if let err = checked(AUGraphStart(self.graph)) {
            return err
        }
        
        NSLog("playing sound...")
        
        let interval = Float32(0.25)
        let intervalCount = Float32(21) / interval
        var azm = Float32(90) // directly to the right
        for var i = Float32(0); i < intervalCount; i++ {
            setObjectCoordinates(self.mixerAu, 0, 0, i * interval/4)
            setObjectCoordinates(self.mixerAu, 1, azm, 1)
            azm -= 180 / intervalCount
            usleep(useconds_t(interval * 1000.0 * 1000.0))
        }
        
        NSLog("playback complete.")

        AUGraphStop(self.graph)
        AUGraphUninitialize(self.graph)
        AUGraphClose(self.graph)
        for sound in self.sounds {
            sound.cleanup()
        }
        
        return noErr
    }
    
    func createAuGraph() -> OSStatus {
        if let err = checked(NewAUGraph(&self.graph)) {
            NSLog("error: could not open augraph")
            return err
        }
        
        var outputCd = AudioComponentDescription()
        outputCd.componentType = OSType(kAudioUnitType_Output)
        outputCd.componentSubType = OSType(kAudioUnitSubType_DefaultOutput)
        outputCd.componentManufacturer = OSType(kAudioUnitManufacturer_Apple)
        var outputNode = AUNode()
        if let err = checked(AUGraphAddNode(self.graph, &outputCd, &outputNode)) {
            NSLog("error: could not add output to augraph")
            return err
        }

        var soundNodes: Array<AUNode> = [];
        for sound in self.sounds {
            var fileCd = AudioComponentDescription()
            fileCd.componentType = OSType(kAudioUnitType_Generator)
            fileCd.componentSubType = OSType(kAudioUnitSubType_AudioFilePlayer)
            fileCd.componentManufacturer = OSType(kAudioUnitManufacturer_Apple)
            var fileNode = AUNode()
            if let err = checked(AUGraphAddNode(self.graph, &fileCd, &fileNode)) {
                NSLog("error: could not add sound file to augraph")
                return err
            }
            soundNodes.append(fileNode)
        }

        var mixerCd = AudioComponentDescription()
        mixerCd.componentType = OSType(kAudioUnitType_Mixer)
        mixerCd.componentSubType = OSType(kAudioUnitSubType_3DMixer)
        mixerCd.componentManufacturer = OSType(kAudioUnitManufacturer_Apple)
        var mixerNode = AUNode()
        if let err = checked(AUGraphAddNode(self.graph, &mixerCd, &mixerNode)) {
            NSLog("error: could not add mixer to augraph")
            return err
        }

        if let err = checked(AUGraphOpen(self.graph)) {
            NSLog("error: could not open augraph")
            return err
        }

        for (sound, soundNode) in Zip2(self.sounds, soundNodes) {
            if let err = checked(AUGraphNodeInfo(self.graph, soundNode, nil, &sound.au)) {
                NSLog("error: could not get player audio unit")
                return err
            }
        }

        if let err = checked(AUGraphNodeInfo(self.graph, mixerNode, nil, &self.mixerAu)) {
            NSLog("error: could not get mixer audio unit")
            return err
        }

        if let err = checked(AUGraphNodeInfo(self.graph, outputNode, nil, &self.outputAu)) {
            NSLog("error: could not get output audio unit")
            return err
        }
        
        if let err = checked(setMixerBusCount(self.mixerAu, 64)) {
            NSLog("error: could not set mixer bus count")
            return err
        }

        var reverbSetting: UInt32 = 1
        let reverbSettingSize: UInt32 = 4
        if let err = checked(AudioUnitSetProperty(self.mixerAu,
            OSType(kAudioUnitProperty_UsesInternalReverb),
            OSType(kAudioUnitScope_Global),
            0,
            &reverbSetting,
            reverbSettingSize)) {
                NSLog("error: could not turn on reverb")
                return err
        }

        for (sound, soundNode) in Zip2(self.sounds, soundNodes) {
            if let err = checked(AUGraphConnectNodeInput(self.graph, soundNode, 0, mixerNode, sound.busIndex)) {
                NSLog("error: could not wire sound into mixer")
                return err
            }
        }

        if let err = checked(AUGraphConnectNodeInput(self.graph, mixerNode, 0, outputNode, 0)) {
            NSLog("error: could not wire mixer into output")
            return err
        }

        if let err = checked(AUGraphInitialize(self.graph)) {
            NSLog("error: could not initialize augraph")
            return err
        }
        
        for sound in self.sounds {
            setReverb(self.mixerAu, sound.busIndex)
        }
        
        return noErr
    }
}

class Sound {
    let fileName: String
    let busIndex: UInt32
    var file = AudioFileID()
    var format = AudioStreamBasicDescription()
    var au = AudioUnit()
    var duration: Float64?
    
    init(fileName: String, busIndex: UInt32) {
        self.fileName = fileName
        self.busIndex = busIndex
    }
    
    func open() -> OSStatus {
        let soundUrl = NSURL(fileURLWithPath: self.fileName)
        let audioAsset = AVURLAsset(URL: soundUrl, options: nil)
        self.duration = CMTimeGetSeconds(audioAsset.duration)
        if let err = checked(AudioFileOpenURL(soundUrl, Int8(kAudioFileReadPermission), 0, &self.file)) {
            NSLog("error: could not open sound file url")
            return err
        }
        
        var formatPropSize = UInt32(sizeof(AudioStreamBasicDescription))
        if let err = checked(AudioFileGetProperty(self.file, OSType(kAudioFilePropertyDataFormat), &formatPropSize, &self.format)) {
            NSLog("error: could not get file properties")
            return err
        }
        
        if let err = checked(AudioUnitSetProperty(self.au, OSType(kAudioUnitProperty_ScheduledFileIDs), OSType(kAudioUnitScope_Global), 0, &self.file, UInt32(sizeof(AudioFileID)))) {
            NSLog("error: could not set audio unit properties")
            return err
        }
        
        var numPackets: UInt64 = 0
        var numPacketsPropSize = UInt32(sizeof(UInt64))
        if let err = checked(AudioFileGetProperty(self.file, OSType(kAudioFilePropertyAudioDataPacketCount), &numPacketsPropSize, &numPackets)) {
            NSLog("error: could not get number of packets for sound file")
        }
        
        var region = ScheduledAudioFileRegion()
        region.mTimeStamp.mFlags = OSType(kAudioTimeStampSampleTimeValid)
        region.mTimeStamp.mSampleTime = 0
        region.mCompletionProc = nil
        region.mCompletionProcUserData = nil
        region.mAudioFile = self.file
        region.mLoopCount = 1
        region.mStartFrame = 0
        region.mFramesToPlay = UInt32(numPackets) * self.format.mFramesPerPacket

        if let err = checked(AudioUnitSetProperty(self.au,
            OSType(kAudioUnitProperty_ScheduledFileRegion),
            OSType(kAudioUnitScope_Global),
            0,
            &region,
            UInt32(sizeof(ScheduledAudioFileRegion)))) {
                NSLog("error: could not set playback region")
                return err
        }
        
        var startTime = AudioTimeStamp()
        startTime.mFlags = OSType(kAudioTimeStampSampleTimeValid)
        startTime.mSampleTime = -1;

        if let err = checked(AudioUnitSetProperty(self.au,
            OSType(kAudioUnitProperty_ScheduleStartTimeStamp),
            OSType(kAudioUnitScope_Global),
            0,
            &startTime,
            UInt32(sizeof(AudioTimeStamp)))) {
                NSLog("error: could not set start time")
                return err
        }
        
        return noErr
    }
    
    func cleanup() {
        AudioFileClose(self.file)
    }
}

main()
