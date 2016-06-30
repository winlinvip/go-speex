// The MIT License (MIT)
//
// Copyright (c) 2016 winlin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package speex_test

import (
	"fmt"
	"github.com/winlinvip/go-speex/speex"
)

func ExampleSpeexDecoder() {
	var err error
	d := speex.NewSpeexDecoder()
	if err = d.Init(16000, 1); err != nil {
		fmt.Println("init decoder failed, err is", err)
		return
	}
	defer d.Close()

	fmt.Println("FrameSize:", d.FrameSize())
	fmt.Println("SampleRate:", d.SampleRate())

	var pcm []byte
	var frame []byte = []byte{
		0x3d, 0xdc, 0x20, 0x13, 0xf3, 0x00, 0x00, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x80, 0x61, 0xbf, 0xff, 0xf7, 0x6e, 0x3a, 0xa2, 0xff, 0xff, 0xf6, 0x01, 0x37, 0xd7, 0x49, 0x9d, 0xf7, 0xdf, 0xf2,
		0x6f, 0x63, 0xda, 0xcd, 0xa4, 0x18, 0x47, 0xe6, 0x19, 0x47, 0x96, 0xf4, 0x32, 0xe6, 0x21, 0x26, 0x8d, 0x12, 0xee,
		0x6d, 0x7c, 0x5b, 0x3f, 0x3c, 0x5f, 0xd7, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0x6a, 0xba,
		0xb8, 0x4a, 0x74, 0x9a, 0xb4, 0x2d, 0xd8, 0xd8, 0xe1, 0xc3, 0x47, 0x25, 0xe8, 0x05, 0xa3, 0xbb, 0xd7, 0x66, 0x3a,
		0x1b, 0xb7, 0xa4, 0x7d, 0xa2, 0xab, 0xfe, 0xd9, 0x08, 0x2c, 0x47,
	}
	if pcm,err = d.Decode(frame); err != nil {
		fmt.Println("decode failed, err is", err)
		return
	}
	fmt.Println("Frame:", len(frame))
	fmt.Println("PCM:", len(pcm))

	// Output:
	// FrameSize: 320
	// SampleRate: 16000
	// Frame: 106
	// PCM: 640
}
