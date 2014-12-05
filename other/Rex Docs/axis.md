PACKAGE DOCUMENTATION

package axisrex
    import "dctech/axis/axisrex"



FUNCTIONS

func Command_Del(script *rex.Script, params []*rex.Value)
    Delete a AXIS file or directory.

	axis:del ds path

    Returns true or false.

func Command_Exists(script *rex.Script, params []*rex.Value)
    Does a AXIS file exist?

	axis:exists ds path

    Returns true or false.

func Command_GetChild(script *rex.Script, params []*rex.Value)
    Get an AXIS DataSource from an AXIS Collection by path.

	axis:getchild collection path

    Returns the DataSource or an error message. May set the Error flag.

func Command_IsDir(script *rex.Script, params []*rex.Value)
    Is a AXIS file a directory?

	axis:isdir ds path

    Returns true or false.

func Command_Mount(script *rex.Script, params []*rex.Value)
    Mount an AXIS DataSource on an AXIS Collection.

	axis:mount collection path ds

    Returns unchanged.

func Command_NewDir(script *rex.Script, params []*rex.Value)
    Create a new empty AXIS logical directory.

	axis:newdir

    Returns a reference to the new directory.

func Command_NewFile(script *rex.Script, params []*rex.Value)
    Create a new read/write AXIS logical file.

	axis:newfile content

    Returns a reference to the new file.

func Command_Read(script *rex.Script, params []*rex.Value)
    Read from a AXIS file.

	axis:read ds path

    Returns file contents or an error message. May set the Error flag.

func Command_WalkDirs(script *rex.Script, params []*rex.Value)
    Iterate over all the directories in a AXIS directory.

	axis:walkdirs ds path code

    Runs code for every directory found, params:

	path

    code MUST be a block created via a block declaration! Returns unchanged.

func Command_WalkFiles(script *rex.Script, params []*rex.Value)
    Iterate over all the files in a AXIS directory.

	axis:walkfiles ds path code

    Runs code for every directory found, params:

	path

    code MUST be a block created via a block declaration! Returns unchanged.

func Command_Write(script *rex.Script, params []*rex.Value)
    Write to a AXIS file.

	axis:write ds path content

    If the directories in the path do not exist axis:write tries to create
    them. Returns unchanged or an error message. May set the Error flag.


