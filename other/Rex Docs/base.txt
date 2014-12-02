PACKAGE DOCUMENTATION

package base
    import "dctech/rex/commands/base"



FUNCTIONS

func Command_Break(script *rex.Script, params []*rex.Value)
    A "soft" return, break will never return more than one level.

	break [value]

    Calling break inside a loop or if command will return from the current
    BLOCK not the command itself, this makes break good for ensuring if
    returns a specific value and/or "returning" a value to loop. Returns
    value or unchanged.

func Command_BreakIf(script *rex.Script, params []*rex.Value)
    A conditional version of break.

	breakif condition [value]

    Returns value or unchanged.

func Command_BreakLoop(script *rex.Script, params []*rex.Value)
    Forces a return until it hits a loop or foreach command or the script
    exits.

	breakloop [value]

    Returns value or unchanged.

func Command_BreakLoopIf(script *rex.Script, params []*rex.Value)
    A conditional version of breakloop.

	breakloopif condition [value]

    Returns value or unchanged.

func Command_Copy(script *rex.Script, params []*rex.Value)
    Copies a value.

	copy value

    Note that this command may not be useful for some types. The new value
    has invalid position info. Returns the new value.

func Command_Error(script *rex.Script, params []*rex.Value)
    Manipulates the error flag.

	error [value]

    If you pass no parameters the error flag will be returned, to set or
    unset the flag pass a boolean value. Returns unchanged or the value of
    the error flag.

func Command_Eval(script *rex.Script, params []*rex.Value)
    Runs a code block or (if the value is not code) converts the value to a
    string and compiles/runs that.

	eval code

    Does not halt any exit states except break and return. Returns result of
    running code.

func Command_Exists(script *rex.Script, params []*rex.Value)
    Returns true if an index exists in a map or array.

	exists value index

    Returns true or false.

func Command_Exit(script *rex.Script, params []*rex.Value)
    Exit the script.

	exit [value]

    Returns value or unchanged.

func Command_For(script *rex.Script, params []*rex.Value)
    Runs code (end-start)/step times.

	for start end step code

    Params for code:

	index

    If code returns false for exits early. code MUST be a block created via
    a block declaration! Does not stop returns, but does work with
    breakloop. A step value of 0 will default to 1. Negative step values are
    allowed, in which case the loop will count down instead of up. Returns
    the return value of the last command in the code.

func Command_ForEach(script *rex.Script, params []*rex.Value)
    Runs code for each entry in a map or array value.

	foreach objectvalue code

    Params for code:

	key value

    If code returns false foreach aborts. code MUST be a block created via a
    block declaration! Does not stop returns, but does work with breakloop.
    Returns the return value of the last command in code.

func Command_If(script *rex.Script, params []*rex.Value)
    If the condition is true run true code else if false code exists call
    false code.

	if condition truecode [falsecode]

    Returns the return value of the last command in the code it runs or
    unchanged.

func Command_IsNil(script *rex.Script, params []*rex.Value)
    Checks if a value is nil.

	isnil value

    Returns true or false.

func Command_Len(script *rex.Script, params []*rex.Value)
    Fetches the element count of an Indexable.

	len value

    Returns the element count.

func Command_Loop(script *rex.Script, params []*rex.Value)
    Runs code for as long as the code returns true.

	loop code

    Returns the return value of the last command in the code it runs, always
    false unless loop exited with ret (In which case the return value is
    unusable by the command calling loop anyway).

func Command_ModVal(script *rex.Script, params []*rex.Value)
    Copies the second value over the first, allows you to do some weird
    pointer-like things. In general do not use unless you know what you are
    doing.

	modval a b

    Returns unchanged.

func Command_Nop(script *rex.Script, params []*rex.Value)
    Does nothing.

	nop

    Returns unchanged.

func Command_OnError(script *rex.Script, params []*rex.Value)
    Runs code if the error flag is true, the error flag is cleared before
    the code is run.

	onerror code

    Returns unchanged.

func Command_Ret(script *rex.Script, params []*rex.Value)
    Return from current command.

	ret [value]

    Some commands will be bypassed like if and loop, for example calling ret
    from inside a loop will not return from the loop, it will return from
    the command that called loop. See break. Returns value or unchanged.

func Command_Type(script *rex.Script, params []*rex.Value)
    Reads or checks value types.

	type value [typ_string]

    Returns the type as a string (if called without a type string) else
    returns true or false.


