PACKAGE DOCUMENTATION

package str
    import "dctech/rex/commands/str"



FUNCTIONS

func Command_Add(script *rex.Script, params []*rex.Value)
    Appends two or more strings together.

	str:add a b [c...]

    Returns the result of appending all parameters together.

func Command_Char(script *rex.Script, params []*rex.Value)
    Gets char at pos.

	str:char a pos

    Operand pos is converted to a 64 bit integer. Invalid strings are given
    the value 0 If the position is out of range returns unchanged and sets
    the error flag. Returns the character.

func Command_Cmp(script *rex.Script, params []*rex.Value)
    Compare two strings.

	str:cmp a b

    Returns true or false.

func Command_Find(script *rex.Script, params []*rex.Value)
    Search for a substring in a string.

	str:find string substring

    Returns the position of the substring or -1 if the substring is not
    present.

func Command_Fmt(script *rex.Script, params []*rex.Value)
    Formats a string.

	str:fmt format values...

    valid format "verbs" for fmtstr:

	%% literal percent
	%s the raw string
	%d as a decimal number
	%x as a lowercase hexadecimal number
	%X as an upercase hexadecimal number

    Returns the formatted string.

func Command_Left(script *rex.Script, params []*rex.Value)
    Return x characters from the left side of a string.

	str:left str x

    Returns the new string or (if the index is out of range) returns the
    string and sets the error flag.

func Command_Len(script *rex.Script, params []*rex.Value)
    Gets the length of a string.

	str:len a

    Returns the length.

func Command_Mid(script *rex.Script, params []*rex.Value)
    Returns count characters from a string.

	str:mid str start count

    Returns the new string or (if start or count is out of range) returns as
    close a result as possible and sets the error flag.

func Command_Replace(script *rex.Script, params []*rex.Value)
    Replaces search with replace in source.

	str:replace source search replace occurrence

    Occurrence gives a number of times to carry out the replacement, use -1
    to replace all. Returns the new string.

func Command_Right(script *rex.Script, params []*rex.Value)
    Return x characters from the right side of a string.

	str:right str x

    Returns the new string or (if the index is out of range) returns the
    string and sets the error flag.

func Command_ToLower(script *rex.Script, params []*rex.Value)
    Converts str to lower case.

	str:tolower str

    Returns the string with all letters converted to lower case.

func Command_ToUpper(script *rex.Script, params []*rex.Value)
    Converts str to upper case.

	str:toupper str

    Returns the string with all letters converted to upper case.

func Command_TrimLeft(script *rex.Script, params []*rex.Value)
    Remove x characters from the left side of a string.

	str:trimleft str x

    Returns the new string or (if the index is out of range) returns the
    string and sets the error flag.

func Command_TrimRight(script *rex.Script, params []*rex.Value)
    Remove x characters from the right side of a string.

	str:trimright str x

    Returns the new string or (if the index is out of range) returns the
    string and sets the error flag.

func Command_TrimSpace(script *rex.Script, params []*rex.Value)
    Trims leading and trailing whitespace from a string.

	str:trimspace str

    Returns str with leading and trailing whitespace removed.


