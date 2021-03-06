PACKAGE DOCUMENTATION

package rubble
    import "rubble"



FUNCTIONS

func Command_Abort(script *rex.Script, params []*rex.Value)
    Causes rubble to abort with an error, use for correctable errors like
    configuration problems.

	rubble:abort msg

    Returns unchanged.

func Command_CallTemplate(script *rex.Script, params []*rex.Value)
    Calls a Rubble template.

	rubble:calltemplate name [params...]

    Returns the templates return value.

func Command_Compress(script *rex.Script, params []*rex.Value)
    Compresses a string using the Rubble encoding.

	rubble:compress string

    Returns the encoded and compressed text.

func Command_ConfigVar(script *rex.Script, params []*rex.Value)
    Manages Rubble variables.

	rubble:configvar name [value]

    Returns unchanged or the variables value.

func Command_CurrentFile(script *rex.Script, params []*rex.Value)
    Returns the name of the current file.

	rubble:currentfile

    Returns the file name.

func Command_Decompress(script *rex.Script, params []*rex.Value)
    Decompresses a Rubble encoded string.

	rubble:decompress string

    Returns the decoded and decompressed text.

func Command_ExpandVars(script *rex.Script, params []*rex.Value)
    Expands Rubble variables.

	rubble:expandvars raws

    Returns the raws with all Rubble variables expanded.

func Command_FileTag(script *rex.Script, params []*rex.Value)
    Manages file tags.

	rubble:filetag filename tag [value]

    Returns the tag's state (if called without a value) or returns
    unchanged.

func Command_GFileTag(script *rex.Script, params []*rex.Value)
    Manages global script file tags.

	rubble:gfiletag filename tag [value]

    Returns the tag's state (if called without a value) or returns
    unchanged.

func Command_NewAddon(script *rex.Script, params []*rex.Value)
    Creates an empty addon object and adds it to the addon list.

	rubble:newaddon name

    Does nothing if the addon already exists. Returns the addon.

func Command_NewFile(script *rex.Script, params []*rex.Value)
    Adds a new file to an addon.

	rubble:newfile addon name contents

    Fails silently if the addon does not exist. Returns unchanged.

func Command_Parse(script *rex.Script, params []*rex.Value)
    Parses Rubble code.

	rubble:stageparse code [stage]

    Note that how code is parsed depends on the parse stage. Valid values
    for stage are:

	-1 (or just leave it off) to use the current stage
	3 for preparse
	4 for parse
	5 for postparse

    The other stage numbers are not valid for the stage parser. Returns the
    result of running code through the stage parser.

func Command_Patch(script *rex.Script, params []*rex.Value)
    Applies a patch to a string.

	rubble:patch string patch

    Returns the patched text.

func Command_Template(script *rex.Script, params []*rex.Value)
    Defines a Rubble script template.

	rubble:template name code

    code MUST be a block created via a block declaration! Parameter names,
    count, and default values is determined by the block meta-data. Returns
    unchanged.

func Command_UserTemplate(script *rex.Script, params []*rex.Value)
    Defines a Rubble user template.

	rubble:usertemplate name [params...] code

    Returns unchanged.


SUBDIRECTORIES

	encoder
	interface
	rblutil
	std_lib

