/*
For copyright/license see header in file "doc.go"
*/

package raptor

import "fmt"
import "errors"
import "bytes"
import "encoding/binary"
import "compress/flate"

// The functions in this file are for end users to allow storing scripts in binary format.

// LoadFile will load a file and add its code to a script. Automatically handles both source and binary files.
// The error is the result from LoadBinary if it is called and fails.
func LoadFile(filename string, code []byte, script *Script) error {
	if ValidateBinary(code) {
		bin, err := LoadBinary(code)
		if err != nil {
			return err
		}
		script.Code.AddCode(bin)
		return nil
	}
	script.Code.Add(NewLexer(string(code), NewPosition(1, 1, filename)))
	return nil
}

// LoadBinary reads the file signature and calls the proper loader function.
func LoadBinary(code []byte) (*Code, error) {
	sig := string(code[0:8])

	switch sig {
	case "Raptor 1": // Normal Binary
		return LoadBinaryV1(code)
	case "Raptor 2": // Small Binary
		return LoadBinaryV2(code)
	case "Raptor 3": // Normal Binary, Compressed string table
		return LoadBinaryV3(code)
	case "Raptor 4": // Small Binary, Compressed string table
		return LoadBinaryV4(code)
	}
	panic("File passed to LoadBinary has invalid signature: " + sig)
}

// ValidateBinary validates the file signature.
func ValidateBinary(code []byte) bool {
	sig := string(code[0:8])

	switch sig {
	case "Raptor 1": // Normal Binary
		return true
	case "Raptor 2": // Small Binary
		return true
	case "Raptor 3": // Normal Binary, Compressed string table
		return true
	case "Raptor 4": // Small Binary, Compressed string table
		return true
	}
	return false
}

// LoadBinaryV1 loads version 1 Raptor binaries ("normal" binaries).
// See WriteBinaryV1 for format spec.
func LoadBinaryV1(code []byte) (script *Code, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			script = nil
		}
	}()

	// Validate signature
	sig := string(code[0:8])
	if sig != "Raptor 1" {
		panic("File passed to LoadBinaryV1 has non-version 1 signature: " + sig)
	}

	// Read header
	headersec := bytes.NewBuffer(code[8:12])
	var tokens int16
	err := binary.Read(headersec, binary.BigEndian, &tokens)
	if err != nil {
		panic(err)
	}
	var strings int16
	err = binary.Read(headersec, binary.BigEndian, &strings)
	if err != nil {
		panic(err)
	}

	// Read code
	codesec := bytes.NewBuffer(code[12 : tokens*2+12])
	tokenlist := make([]int16, tokens)

	err = binary.Read(codesec, binary.BigEndian, &tokenlist)
	if err != nil {
		panic(err)
	}

	// Read string table
	stringsec := bytes.NewBuffer(code[tokens*2+12:])
	stringtbl := make([]string, 0, 20)
	for i := 0; i < int(strings); i++ {
		var strlen int16
		err = binary.Read(stringsec, binary.BigEndian, &strlen)
		if err != nil {
			panic(err)
		}

		str := make([]byte, strlen)
		len, err := stringsec.Read(str)
		if err != nil {
			panic(err)
		}
		if len != int(strlen) {
			panic("Reported string length is not equal to amount read.")
		}

		stringtbl = append(stringtbl, string(str))
	}

	// Generate *Code
	this := new(Code)
	this.Code = make([]uint32, tokens)
	this.StringTable = stringtbl

	for i := range tokenlist {
		this.Code[i] = uint32(tokenlist[i])
	}

	return this, nil
}

// WriteBinaryV1 writes version 1 Raptor binaries.
//
// Format:
// Version 1, Big Endian encoding
// Signature, 8 bytes, "Raptor 1"
// Token count, 2 bytes
// String count, 2 bytes
// Code stream, 2 bytes per token, Token count * 2 bytes long
// String table, variable entry size, 2 byte record length, record length bytes of data, repeats String count times
//
// Notes:
// Should work for all (normal) scripts.
// Max string size is 65535 bytes.
// Max unique string count is 65525.
// Max token count is 65535.
func WriteBinaryV1(script *Code) (code []byte, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			code = nil
		}
	}()

	out := new(bytes.Buffer)

	// Write Header
	out.Write([]byte("Raptor 1"))

	if len(script.Code) > 0xFFFF {
		panic("Too many tokens to encode in a version 1 binary.")
	}
	err := binary.Write(out, binary.BigEndian, int16(len(script.Code)))
	if err != nil {
		panic(err)
	}

	if len(script.StringTable) > 0xFFFF-10 {
		panic("String table too long to encode in a version 1 binary.")
	}
	err = binary.Write(out, binary.BigEndian, int16(len(script.StringTable)))
	if err != nil {
		panic(err)
	}

	// Write Token Stream
	for i := range script.Code {
		err := binary.Write(out, binary.BigEndian, int16(script.Code[i]))
		if err != nil {
			panic(err)
		}
	}

	// Write String Table
	for i := range script.StringTable {
		str := []byte(script.StringTable[i])
		if len(str) > 0xFFFF {
			panic("String too long to encode in a version 1 binary.")
		}
		err := binary.Write(out, binary.BigEndian, int16(len(str)))
		if err != nil {
			panic(err)
		}
		out.Write(str)
	}

	return out.Bytes(), nil
}

// LoadBinaryV2 loads version 2 Raptor binaries ("small" binaries).
// See WriteBinaryV2 for format spec.
func LoadBinaryV2(code []byte) (script *Code, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			script = nil
		}
	}()

	// Validate signature
	sig := string(code[0:8])
	if sig != "Raptor 2" {
		panic("File passed to LoadBinaryV2 has non-version 2 signature: " + sig)
	}

	// Read header
	headersec := bytes.NewBuffer(code[8:11])
	var tokens uint16
	err := binary.Read(headersec, binary.BigEndian, &tokens)
	if err != nil {
		panic(err)
	}

	var strings uint8
	err = binary.Read(headersec, binary.BigEndian, &strings)
	if err != nil {
		panic(err)
	}

	// Read code
	codesec := bytes.NewBuffer(code[11 : tokens+11])
	tokenlist := make([]uint8, tokens)

	err = binary.Read(codesec, binary.BigEndian, &tokenlist)
	if err != nil {
		panic(err)
	}

	// Read string table
	stringsec := bytes.NewBuffer(code[tokens+11:])
	stringtbl := make([]string, 0, 20)
	for i := 0; i < int(strings); i++ {
		var strlen int16
		err = binary.Read(stringsec, binary.BigEndian, &strlen)
		if err != nil {
			panic(err)
		}

		str := make([]byte, strlen)
		len, err := stringsec.Read(str)
		if err != nil {
			panic(err)
		}
		if len != int(strlen) {
			panic("Reported string length is not equal to amount read.")
		}

		stringtbl = append(stringtbl, string(str))
	}

	// Generate *Code
	this := new(Code)
	this.Code = make([]uint32, tokens)
	this.StringTable = stringtbl

	for i := range tokenlist {
		this.Code[i] = uint32(tokenlist[i])
	}

	return this, nil
}

// WriteBinaryV2 writes version 2 Raptor binaries.
//
// Format:
// Version 2, Big Endian encoding
// Signature, 8 bytes, "Raptor 2"
// Token count, 2 bytes
// String count, 1 byte
// Code stream, 1 byte per token, Token count bytes long
// String table, variable entry size, 2 byte record length, record length bytes of data, repeats String count times
//
// Notes:
// This format produces smaller binaries at the cost of a smaller maximum string table.
// Max string size is 65535 bytes.
// Max unique string count is 245.
// Max token count is 65535.
func WriteBinaryV2(script *Code) (code []byte, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			code = nil
		}
	}()

	out := new(bytes.Buffer)

	// Write Header
	out.Write([]byte("Raptor 2"))
	if len(script.Code) > 0xFFFF {
		panic("Too many tokens to encode in a version 2 binary.")
	}
	err := binary.Write(out, binary.BigEndian, uint16(len(script.Code)))
	if err != nil {
		panic(err)
	}

	if len(script.StringTable) > 0xFF-10 {
		panic("String table too long to encode in a version 2 binary.")
	}
	err = binary.Write(out, binary.BigEndian, uint8(len(script.StringTable)))
	if err != nil {
		panic(err)
	}

	// Write Token Stream
	for i := range script.Code {
		err := binary.Write(out, binary.BigEndian, uint8(script.Code[i]))
		if err != nil {
			panic(err)
		}
	}

	// Write String Table
	for i := range script.StringTable {
		str := []byte(script.StringTable[i])
		if len(str) > 0xFFFF {
			panic("String too long to encode in a version 2 binary.")
		}
		err := binary.Write(out, binary.BigEndian, int16(len(str)))
		if err != nil {
			panic(err)
		}
		out.Write(str)
	}

	return out.Bytes(), nil
}

// LoadBinaryV3 loads version 3 Raptor binaries ("normal" binaries with compressed string table).
// See WriteBinaryV3 for format spec.
func LoadBinaryV3(code []byte) (script *Code, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			script = nil
		}
	}()

	// Validate signature
	sig := string(code[0:8])
	if sig != "Raptor 3" {
		panic("File passed to LoadBinaryV3 has non-version 3 signature: " + sig)
	}

	// Read header
	headersec := bytes.NewBuffer(code[8:12])
	var tokens int16
	err := binary.Read(headersec, binary.BigEndian, &tokens)
	if err != nil {
		panic(err)
	}
	var strings int16
	err = binary.Read(headersec, binary.BigEndian, &strings)
	if err != nil {
		panic(err)
	}

	// Read code
	codesec := bytes.NewBuffer(code[12 : tokens*2+12])
	tokenlist := make([]int16, tokens)

	err = binary.Read(codesec, binary.BigEndian, &tokenlist)
	if err != nil {
		panic(err)
	}

	// Read string table
	stringsec := flate.NewReader(bytes.NewBuffer(code[tokens*2+12:]))
	stringtbl := make([]string, 0, 20)
	for i := 0; i < int(strings); i++ {
		var strlen int16
		err = binary.Read(stringsec, binary.BigEndian, &strlen)
		if err != nil {
			panic(err)
		}

		str := make([]byte, strlen)
		len, err := stringsec.Read(str)
		if err != nil {
			panic(err)
		}
		if len != int(strlen) {
			panic("Reported string length is not equal to amount read.")
		}

		stringtbl = append(stringtbl, string(str))
	}
	stringsec.Close()

	// Generate *Code
	this := new(Code)
	this.Code = make([]uint32, tokens)
	this.StringTable = stringtbl

	for i := range tokenlist {
		this.Code[i] = uint32(tokenlist[i])
	}

	return this, nil
}

// WriteBinaryV3 writes version 3 Raptor binaries.
//
// Format:
// Version 3, is just version 1 with a DEFLATE compressed string table.
func WriteBinaryV3(script *Code) (code []byte, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			code = nil
		}
	}()

	out := new(bytes.Buffer)

	// Write Header
	out.Write([]byte("Raptor 3"))

	if len(script.Code) > 0xFFFF {
		panic("Too many tokens to encode in a version 3 binary.")
	}
	err := binary.Write(out, binary.BigEndian, int16(len(script.Code)))
	if err != nil {
		panic(err)
	}

	if len(script.StringTable) > 0xFFFF-10 {
		panic("String table too long to encode in a version 3 binary.")
	}
	err = binary.Write(out, binary.BigEndian, int16(len(script.StringTable)))
	if err != nil {
		panic(err)
	}

	// Write Token Stream
	for i := range script.Code {
		err := binary.Write(out, binary.BigEndian, int16(script.Code[i]))
		if err != nil {
			panic(err)
		}
	}

	// Write String Table
	strsec, _ := flate.NewWriter(out, flate.BestCompression)
	for i := range script.StringTable {
		str := []byte(script.StringTable[i])
		if len(str) > 0xFFFF {
			panic("String too long to encode in a version 1 binary.")
		}
		err := binary.Write(strsec, binary.BigEndian, int16(len(str)))
		if err != nil {
			panic(err)
		}
		strsec.Write(str)

	}
	err = strsec.Close()
	if err != nil {
		panic(err)
	}

	return out.Bytes(), nil
}

// LoadBinaryV4 loads version 4 Raptor binaries ("small" binaries).
// See WriteBinaryV4 for format spec.
func LoadBinaryV4(code []byte) (script *Code, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			script = nil
		}
	}()

	// Validate signature
	sig := string(code[0:8])
	if sig != "Raptor 4" {
		panic("File passed to LoadBinaryV4 has non-version 4 signature: " + sig)
	}

	// Read header
	headersec := bytes.NewBuffer(code[8:11])
	var tokens uint16
	err := binary.Read(headersec, binary.BigEndian, &tokens)
	if err != nil {
		panic(err)
	}

	var strings uint8
	err = binary.Read(headersec, binary.BigEndian, &strings)
	if err != nil {
		panic(err)
	}

	// Read code
	codesec := bytes.NewBuffer(code[11 : tokens+11])
	tokenlist := make([]uint8, tokens)

	err = binary.Read(codesec, binary.BigEndian, &tokenlist)
	if err != nil {
		panic(err)
	}

	// Read string table
	stringsec := flate.NewReader(bytes.NewBuffer(code[tokens+11:]))
	stringtbl := make([]string, 0, 20)
	for i := 0; i < int(strings); i++ {
		var strlen int16
		err = binary.Read(stringsec, binary.BigEndian, &strlen)
		if err != nil {
			panic(err)
		}

		str := make([]byte, strlen)
		len, err := stringsec.Read(str)
		if err != nil {
			panic(err)
		}
		if len != int(strlen) {
			panic("Reported string length is not equal to amount read.")
		}

		stringtbl = append(stringtbl, string(str))
	}
	stringsec.Close()

	// Generate *Code
	this := new(Code)
	this.Code = make([]uint32, tokens)
	this.StringTable = stringtbl

	for i := range tokenlist {
		this.Code[i] = uint32(tokenlist[i])
	}

	return this, nil
}

// WriteBinaryV4 writes version 4 Raptor binaries.
//
// Format:
// Version 4, is just version 2 with a DEFLATE compressed string table.
func WriteBinaryV4(script *Code) (code []byte, oops error) {
	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				oops = i
			case string:
				oops = errors.New(i)
			default:
				oops = errors.New(fmt.Sprint(i))
			}
			code = nil
		}
	}()

	out := new(bytes.Buffer)

	// Write Header
	out.Write([]byte("Raptor 4"))
	if len(script.Code) > 0xFFFF {
		panic("Too many tokens to encode in a version 4 binary.")
	}
	err := binary.Write(out, binary.BigEndian, uint16(len(script.Code)))
	if err != nil {
		panic(err)
	}

	if len(script.StringTable) > 0xFF-10 {
		panic("String table too long to encode in a version 4 binary.")
	}
	err = binary.Write(out, binary.BigEndian, uint8(len(script.StringTable)))
	if err != nil {
		panic(err)
	}

	// Write Token Stream
	for i := range script.Code {
		err := binary.Write(out, binary.BigEndian, uint8(script.Code[i]))
		if err != nil {
			panic(err)
		}
	}

	// Write String Table
	strsec, _ := flate.NewWriter(out, flate.BestCompression)
	for i := range script.StringTable {
		str := []byte(script.StringTable[i])
		if len(str) > 0xFFFF {
			panic("String too long to encode in a version 4 binary.")
		}
		err := binary.Write(strsec, binary.BigEndian, int16(len(str)))
		if err != nil {
			panic(err)
		}
		strsec.Write(str)
	}
	err = strsec.Close()
	if err != nil {
		panic(err)
	}

	return out.Bytes(), nil
}
