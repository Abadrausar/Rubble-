/*
Copyright 2013 by Milo Christiansen

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

package axis

type ErrTyp int
const (
	ErrBadPath ErrTyp = iota
	ErrNotFound
	ErrNotDir
	ErrNotFile
	ErrReadOnly
)

func NewError(typ ErrTyp, path string) error {
	return &PathError{
		Path: path,
		Typ: typ,
	}
}

type PathError struct {
	Path string
	Typ ErrTyp
}

func (err *PathError) Error() string {
	switch err.Typ {
	case ErrBadPath:
		return "Bad Path: " + err.Path
	case ErrNotFound:
		return "File/Dir Not Found: " + err.Path
	case ErrNotDir:
		return "Not a Directory: " + err.Path
	case ErrNotFile:
		return "Not a File: " + err.Path
	case ErrReadOnly:
		return "File/Dir Read-only: " + err.Path
	default:
		return "Invalid Error Code: " + err.Path
	}
}
