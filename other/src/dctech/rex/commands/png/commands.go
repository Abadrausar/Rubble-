/*
Copyright 2014 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

// Rex PNG Image Commands.
package png

import "dctech/rex"

import "image"
import "image/color"
import _ "image/png"

import "strings"

// Adds the PNG image commands to the state.
// The PNG image commands are:
//	png:load
func Setup(state *rex.State) (err error) {
	defer func() {
		if !state.NoRecover {
			if x := recover(); x != nil {
				if y, ok := x.(rex.ScriptError); ok {
					err = y
					return
				}
				panic(x)
			}
		}
	}()
	
	mod := state.RegisterModule("png")
	mod.RegisterCommand("load", Command_Load)
	
	return nil
}

// Load a PNG image from a "string".
// 	png:load image_data
// If loading from a file be careful not to mess up the image data.
// Returns a set of nested indexables ([image x y]) pixels are represented by 32 bit alpha premultiplied colors.
func Command_Load(script *rex.Script, params []*rex.Value) {
	if len(params) != 1 {
		rex.ErrorParamCount("png:load", "1")
	}

	// Load image
	imageFile, _, err := image.Decode(strings.NewReader(params[0].String()))
	if err != nil {
		rex.RaiseInternalError(err)
	}
	bounds := imageFile.Bounds()
	
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	columns := make([]*rex.Value, h)
	
	// Convert the pixel data to something easier to read.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		row := make([]*rex.Value, w)
		
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pix := color.RGBAModel.Convert(imageFile.At(x, y)).(color.RGBA)
			var pix32 int64
			pix32 |= int64(pix.R)
			pix32 <<= 8
			pix32 |= int64(pix.G)
			pix32 <<= 8
			pix32 |= int64(pix.B)
			pix32 <<= 8
			pix32 |= int64(pix.A)
			row[x - bounds.Min.X] = rex.NewValueInt64(pix32)
		}
		columns[y - bounds.Min.Y] = rex.NewValueIndex(rex.NewStaticArray(row))
	}
	
	script.RetVal = rex.NewValueIndex(rex.NewStaticArray(columns))
}
