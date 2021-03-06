PACKAGE DOCUMENTATION

package rexdfraw
    import "dctech/rexdfraw"



FUNCTIONS

func Command_Dump(script *rex.Script, params []*rex.Value)
    Dumps a parsed raw file to text.

	df:raw:dump file

    This command is for use with the indexable returned by df:raw:parse.
    Returns a string version of the file.

func Command_Parse(script *rex.Script, params []*rex.Value)
    Parses a raw file.

	df:raw:parse text

    This command is much slower, but also far more flexible, than
    df:raw:walk. Returns an indexable of all the tags in the file (each tag
    is also indexable).

func Command_Tag(script *rex.Script, params []*rex.Value)
    Creates a new raw tag.

	df:raw:tag id [params...]

    Returns a tag of the same format as used by df:raw:parse.

func Command_Walk(script *rex.Script, params []*rex.Value)
    Parses the raws and runs code for every tag. Return false to abort
    early, respects breakloop.

	df:raw:walk file code

    Parameters for code:

	tag

    code MUST be a block created via a block declaration! Returns the file
    (as a string).


