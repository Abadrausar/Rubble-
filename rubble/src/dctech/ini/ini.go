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

// A simple ini parser.
package ini

import "fmt"
import "strings"

// An ini key/value pair.
type Key struct {
	Value   string
	Comment string
}

// An ini section.
// All keys are listed, in first added first order, in KeyOrder. 
// No key is ever removed from KeyOrder so you can not assume that a key listed there in in Content.
// Keys are written out in the order they are listed in KeyOrder.
type Section struct {
	Content  map[string]*Key
	KeyOrder []string
	Comment  string
}

// Create a new section structure and init its fields.
func NewSection() *Section {
	ret := new(Section)
	ret.Content = make(map[string]*Key, 20)
	ret.KeyOrder = make([]string, 0, 20)
	return ret
}

// An ini File.
// All sections are listed, in first added first order, in SectionOrder. 
// No section is ever removed from SectionOrder so you can not assume that a section listed there in in Content.
// Sections are written out in the order they are listed in SectionOrder.
type File struct {
	Content      map[string]*Section
	SectionOrder []string
}

// Create a new file structure and init its fields.
func NewFile() *File {
	ret := new(File)
	ret.Content = make(map[string]*Section, 20)
	ret.SectionOrder = make([]string, 0, 20)
	return ret
}

const (
	SectionNotExist = iota
	KeyNotExist
	MalformedInput
	TokenExpected
)

type Err struct {
	Type int
	Line int
	Dat1 string
	Dat2 string
}

// Possible values:
//  The section: %s Does not exist.
//  The key: %s Does not exist.
//  Malformed input on line: %d
//  Invalid Token: Found: %s. Expected: %s. On line: %d
func (this *Err) Error() string {
	switch this.Type {
	case SectionNotExist:
		return fmt.Sprintf("The section: %s Does not exist.", this.Dat1)
	case KeyNotExist:
		return fmt.Sprintf("The key: %s Does not exist.", this.Dat1)
	case MalformedInput:
		return fmt.Sprintf("Malformed input on line: %d", this.Line)
	case TokenExpected:
		return fmt.Sprintf("Invalid Token: Found: %s. Expected: %s. On line: %d", this.Dat1, this.Dat2, this.Line)
	}
	panic("Invalid error type.")
}

func newErr(kind int, dat1, dat2 string, line int) error {
	err := new(Err)
	err.Type = kind
	err.Dat1 = dat1
	err.Dat2 = dat2
	err.Line = line
	return err
}

// Parse an ini file
func Parse(text string) (ret *File, err error) {
	err = nil
	ret = NewFile()

	defer func() {
		if x := recover(); x != nil {
			switch i := x.(type) {
			case error:
				err = i
			default:
				err = fmt.Errorf("%v", i)
			}
		}
	}()

	lex := NewLexer(text, 1)

	com := acumulateComments(lex)
	for {
		if !lex.CheckLookAhead(TknSection) {
			break
		}
		com = parseSection(lex, ret, com)
	}
	return
}

func parseSection(lex *Lexer, file *File, comment string) string {
	section := NewSection()

	section.Comment = comment
	lex.GetToken(TknSection)
	sectionname := lex.Current.Lexeme
	retcom := ""
	for {
		com := acumulateComments(lex)
		if !lex.CheckLookAhead(TknKey) {
			retcom = com
			break
		}
		parseKey(lex, section, com)
	}
	createSection(sectionname, file)
	file.Content[sectionname] = section
	return retcom
}

func parseKey(lex *Lexer, section *Section, comment string) {
	key := new(Key)
	key.Comment = comment
	lex.GetToken(TknKey)
	keyname := lex.Current.Lexeme
	lex.GetToken(TknValue)
	key.Value = lex.Current.Lexeme
	createKey(keyname, section)
	section.Content[keyname] = key
}

func acumulateComments(lex *Lexer) string {
	ret := ""
	for lex.CheckLookAhead(TknComment) {
		lex.GetToken(TknComment)
		ret += lex.Current.Lexeme + "\n"
	}
	return ret
}

// Format an ini file as text ready to write to disk
func Format(file *File) string {
	ret := make([]byte, 0, 100)
	for _, val := range file.SectionOrder {
		if sec, ok := file.Content[val]; ok {
			ret = formatSection(val, sec, ret)
		}
	}
	return string(ret[:len(ret)-1])
}

func formatComment(comment string) []byte {
	strsplit := strings.Split(comment, "\n")
	ret := make([]byte, 0, 100)

	for _, val := range strsplit {
		ret = append(ret, []byte("#"+val+"\n")...)
	}

	return ret[:len(ret)-2]
}

func formatSection(name string, section *Section, out []byte) []byte {
	out = append(out, formatComment(section.Comment)...)
	out = append(out, []byte("["+name+"]\n")...)

	for _, val := range section.KeyOrder {
		if key, ok := section.Content[val]; ok {
			out = formatKey(val, key, out)
		}
	}
	out = append(out, byte('\n'))

	return out
}

func formatKey(name string, key *Key, out []byte) []byte {
	out = append(out, formatComment(key.Comment)...)
	out = append(out, []byte(name+"="+key.Value+"\n")...)
	return out
}

// Get a value
func (file *File) Get(section, key string) (string, error) {
	if sec, ok := file.Content[section]; ok {
		if val, ok := sec.Content[key]; ok {
			return val.Value, nil
		}
		return "", newErr(KeyNotExist, key, "", 0)
	}
	return "", newErr(SectionNotExist, section, "", 0)
}

// Set a value
func (file *File) Set(section, key, value string) error {
	if sec, ok := file.Content[section]; ok {
		if val, ok := sec.Content[key]; ok {
			val.Value = value
			return nil
		}
		return newErr(KeyNotExist, key, "", 0)
	}
	return newErr(SectionNotExist, section, "", 0)
}

// Create a key (and possibly a section)
func (file *File) Create(section, key string) {
	if _, ok := file.Content[section]; !ok {
		file.Content[section] = NewSection()
		file.SectionOrder = append(file.SectionOrder, section)
	}
	if _, ok := file.Content[section].Content[key]; !ok {
		file.Content[section].Content[key] = new(Key)
		file.Content[section].KeyOrder = append(file.Content[section].KeyOrder, key)
	}
}

// Create a section
func createSection(section string, file *File) {
	if _, ok := file.Content[section]; !ok {
		file.Content[section] = NewSection()
		file.SectionOrder = append(file.SectionOrder, section)
	}
}

// Create a key
func createKey(key string, section *Section) {
	if _, ok := section.Content[key]; !ok {
		section.Content[key] = new(Key)
		section.KeyOrder = append(section.KeyOrder, key)
	}
}
