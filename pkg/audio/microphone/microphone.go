// Copyright 2022 Symbl.ai SDK contributors. All Rights Reserved.
// Use of this source code is governed by an Apache-2.0 license that can be found in the LICENSE file.
// SPDX-License-Identifier: Apache-2.0

package microphone

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/gordonklaus/portaudio"
	klog "k8s.io/klog/v2"
)

func Initialize() {
	portaudio.Initialize()
}

func Teardown() {
	portaudio.Terminate()
}

func New(cfg AudioConfig) (*Microphone, error) {
	klog.V(6).Infof("Microphone.New ENTER\n")

	m := &Microphone{
		stopChan: make(chan struct{}),
		intBuf:   make([]int16, 1024),
		muted:    false,
	}

	portaudio.Initialize()

	stream, err := portaudio.OpenDefaultStream(cfg.InputChannels, 0, float64(cfg.SamplingRate), len(m.intBuf), m.intBuf)
	if err != nil {
		klog.V(1).Infof("OpenDefaultStream failed. Err: %v\n", err)
		klog.V(6).Infof("Microphone.New LEAVE\n")
		return nil, err
	}

	// housekeeping
	m.stream = stream

	klog.V(3).Infof("OpenDefaultStream succeded\n")
	klog.V(6).Infof("Microphone.New LEAVE\n")

	return m, nil
}

func (m *Microphone) Start() error {
	err := m.stream.Start()
	if err != nil {
		klog.V(1).Infof("Mic failed to start. Err: %v\n", err)
		return err
	}

	klog.V(3).Infof("Start() succeded\n")
	return nil
}

func (m *Microphone) Read() ([]int16, error) {
	err := m.stream.Read()
	if err != nil {
		klog.V(1).Infof("stream.Read failed. Err: %v\n", err)
		return nil, err
	}

	buf := make([]int16, 1024)
	byteCopied := copy(buf, m.intBuf)
	klog.V(7).Infof("stream.Read bytes copied: %d\n", byteCopied)
	return buf, nil
}

func (m *Microphone) Stream(w io.Writer) error {
	for {
		select {
		case <-m.stopChan:
			return nil
		default:
			err := m.stream.Read()
			if err != nil {
				klog.V(1).Infof("stream.Read failed. Err: %v\n", err)
				return err
			}

			byteCount, err := w.Write(m.int16ToLittleEndianByte(m.intBuf))
			if err != nil {
				klog.V(1).Infof("w.Write failed. Err: %v\n", err)
				return err
			}
			klog.V(7).Infof("io.Writer succeeded. Bytes written: %d\n", byteCount)
		}
	}

	return nil
}

func (m *Microphone) Mute() {
	m.mute.Lock()
	m.muted = true
	m.mute.Unlock()
}

func (m *Microphone) Unmute() {
	m.mute.Lock()
	m.muted = false
	m.mute.Unlock()
}

func (m *Microphone) Stop() error {
	err := m.stream.Stop()
	if err != nil {
		klog.V(1).Infof("stream.Stop failed. Err: %v\n", err)
		return err
	}

	close(m.stopChan)
	<-m.stopChan

	return nil
}

func (m *Microphone) int16ToLittleEndianByte(f []int16) []byte {
	m.mute.Lock()
	isMuted := m.muted
	m.mute.Unlock()

	if isMuted {
		klog.V(7).Infof("Mic is MUTED!\n")
		f = make([]int16, len(f))
	}

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		klog.V(1).Infof("binary.Write failed. Err %v\n", err)
	}

	return buf.Bytes()
}
