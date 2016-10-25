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

// The speex decoder, to decode the encoded speex frame to PCM samples.
package speex

/*
#cgo CFLAGS: -I${SRCDIR}/../speex-lib/objs/include
#cgo LDFLAGS: ${SRCDIR}/../speex-lib/objs/lib/libspeex.a -lm
#include "speex/speex.h"

typedef struct {
	SpeexBits bits;
	const SpeexMode* mode;
	void* state;

	int frame_size;
	int sample_rate;
} speexdec_t;

int speexdec_init(speexdec_t* h, int sample_rate, int channels) {
	h->mode = 0;
	h->state = 0;
	h->frame_size = h->sample_rate = 0;

    // TODO: support stereo speex.
    if (channels != 1) {
        return -1;
    }

    int frame_size;
    const SpeexMode* mode;
    if (1) {
        int spx_mode;
        switch (sample_rate) {
            case 8000:  spx_mode = 0; break;
            case 16000: spx_mode = 1; break;
            case 32000: spx_mode = 2; break;
            default: return -1;
        }

        mode = speex_lib_get_mode(spx_mode);
        if (!mode) {
            return -1;
        }
    }
    h->mode = mode;

    void* state = speex_decoder_init(mode);
    if (!state) {
        return -1;
    }

    h->state = state;
    speex_bits_init(&h->bits);

    if (1) {
        spx_int32_t N;
        speex_decoder_ctl(state, SPEEX_GET_FRAME_SIZE, &N);
        h->frame_size = N;

        speex_decoder_ctl(state, SPEEX_GET_SAMPLING_RATE, &N);
        h->sample_rate = N;
    }

	return 0;
}

void speexdec_close(speexdec_t* h) {
	if (h->state) {
		speex_decoder_destroy(h->state);
		speex_bits_destroy(&h->bits);
	}
	h->state = 0;
}


int speexdec_decode(speexdec_t* h, char* frame, int nb_frame, char* pcm, int* pnb_pcm, int* isDone) {
	// the output pcm must equals to the frames(each is 16bits).
	if (*pnb_pcm != h->frame_size * sizeof(spx_int16_t)) {
		return -1;
	}
	
	if (speex_bits_remaining(&h->bits) < 5 ||
		speex_bits_peek_unsigned(&h->bits, 5) == 0xF) {
		
		speex_bits_read_from(&h->bits, frame, nb_frame);
	}

	spx_int16_t* output = (spx_int16_t*)pcm;
	int ret = speex_decode_int(h->state, &h->bits, output);

	// 0 for no error, -1 for end of stream, -2 corrupt stream
	if (ret <= -2) {
		return ret;
	}

	if (ret == -1) {
		*pnb_pcm = 0;
		return 0;
	}

	if (speex_bits_remaining(&h->bits) < 5 &&
		speex_bits_peek_unsigned(&h->bits, 5) == 0xF) {
		
		*isDone = 1;
	}

	return 0;
}

int speexdec_frame_size(speexdec_t* h) {
	return h->frame_size;
}

int speexdec_sample_rate(speexdec_t* h) {
	return h->sample_rate;
}
*/
import "C"

import (
	"fmt"
	"unsafe"
	"bytes"
)

type SpeexDecoder struct {
	m C.speexdec_t
}

func NewSpeexDecoder() *SpeexDecoder {
	return &SpeexDecoder{}
}

// @remark only support mono speex(channels must be 1).
func (v *SpeexDecoder) Init(sampleRate, channels int) (err error) {
	if channels != 1 {
		return fmt.Errorf("only support mono(1), actual is %v", channels)
	}

	r := C.speexdec_init(&v.m, C.int(sampleRate), C.int(channels))
	if int(r) != 0 {
		return fmt.Errorf("init decoder failed, err=%v", int(r))
	}

	return
}

func (v *SpeexDecoder) Close() {
	C.speexdec_close(&v.m)
}

// @return pcm is nil when EOF.
func (v *SpeexDecoder) Decode(frame []byte) (pcm []byte, err error) {
	p := (*C.char)(unsafe.Pointer(&frame[0]))
	pSize := C.int(len(frame))
	
	pIsDone := C.int(0)
	
	var result bytes.Buffer
	
	for {
		if pIsDone == 1 {
			break
		}
		// each sample is 16bits(2bytes),
		// so we alloc the output to frame_size*2.
		nbPcmBytes := v.FrameSize()*2
		
		pcmTmp := make([]byte, nbPcmBytes)
		pPcm := (*C.char)(unsafe.Pointer(&pcmTmp[0]))
		pNbPcm := C.int(nbPcmBytes)
		
		r := C.speexdec_decode(&v.m, p, pSize, pPcm, &pNbPcm, &pIsDone)
		if int(r) != 0 {
			return nil,fmt.Errorf("decode failed, err=%v", int(r))
		}
		if int(pNbPcm) <= 0 {
			return nil,nil
		}
		if int(pNbPcm) != nbPcmBytes {
			return nil,fmt.Errorf("invalid pcm size %v", int(pNbPcm))
		}
		result.Write(pcmTmp)
	}
	pcm = result.Bytes()
	
	return
}

func (v *SpeexDecoder) FrameSize() int {
	return int(C.speexdec_frame_size(&v.m))
}

func (v *SpeexDecoder) SampleRate() int {
	return int(C.speexdec_sample_rate(&v.m))
}

func (v *SpeexDecoder) Channels() int {
	return 1
}
