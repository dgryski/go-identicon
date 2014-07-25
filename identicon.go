/*
Copyright (c) 2013, Damian Gryski <damian@gryski.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

   * Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.

   * Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package identicon creates simple visual hashes of data as PNGs.
/*
For more information, see https://en.wikipedia.org/wiki/Identicon and https://github.com/blog/1586-identicons
*/
package identicon

import (
	"bytes"
	"hash"
	"image"
	"image/color"
	"image/png"

	"github.com/dchest/siphash"
)

// Renderer allows rendering of data as a PNG identicon
type Renderer interface {
	// Render generates a PNG from data
	Render(data []byte) []byte
}

type identicon struct {
	key    []byte
	sqSize int
	rows   int
	cols   int
	h      hash.Hash64
}

const xborder = 35
const yborder = 35
const maxX = 420
const maxY = 420

// New5x5 creates a new 5-by-5 identicon renderer using 'key' as the hash salt
func New5x5(key []byte) Renderer {
	return &identicon{
		sqSize: 70,
		rows:   5,
		cols:   5,
		h:      siphash.New(key),
	}
}

// New7x7 creates a new 7-by-7 identicon renderer using 'key' as the hash salt
func New7x7(key []byte) Renderer {
	return &identicon{
		sqSize: 50,
		rows:   7,
		cols:   7,
		h:      siphash.New(key),
	}
}

func (icon *identicon) Render(data []byte) []byte {

	icon.h.Reset()
	icon.h.Write(data)
	h := icon.h.Sum64()

	nrgba := color.NRGBA{
		R: uint8(h),
		G: uint8(h >> 8),
		B: uint8(h >> 16),
		A: 0xff,
	}
	h >>= 24

	img := image.NewPaletted(image.Rect(0, 0, maxX, maxY), color.Palette{color.NRGBA{0xf0, 0xf0, 0xf0, 0xff}, nrgba})

	sqx := 0
	sqy := 0

	pixels := make([]byte, icon.sqSize)
	for i := 0; i < icon.sqSize; i++ {
		pixels[i] = 1
	}

	for i := 0; i < icon.rows*(icon.cols+1)/2; i++ {

		if h&1 == 1 {

			for i := 0; i < icon.sqSize; i++ {
				x := xborder + sqx*icon.sqSize
				y := yborder + sqy*icon.sqSize + i
				offs := img.PixOffset(x, y)
				copy(img.Pix[offs:], pixels)

				x = xborder + (icon.cols-1-sqx)*icon.sqSize
				offs = img.PixOffset(x, y)
				copy(img.Pix[offs:], pixels)
			}
		}

		h >>= 1
		sqy++
		if sqy == icon.rows {
			sqy = 0
			sqx++
		}
	}

	var buf bytes.Buffer

	png.Encode(&buf, img)

	return buf.Bytes()
}
