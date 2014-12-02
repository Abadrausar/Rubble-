package main

import "dctech/nca4"
import "dctech/nca4/base"
import "dctech/nca4/bit"
import "dctech/nca4/cmp"
import "dctech/nca4/conio"
import "dctech/nca4/csv"
import "dctech/nca4/env"
import "dctech/nca4/fileio"
import "dctech/nca4/ini"
import "dctech/nca4/math"
import "dctech/nca4/stack"
import "dctech/nca4/str"

var GlobalNCAState *nca4.State

func InitNCA() {
	state := nca4.NewState()
	
	state.NoRecover = !Recover

	// Load all commands
	base.Setup(state)
	bit.Setup(state)
	cmp.Setup(state)
	conio.Setup(state)
	csv.Setup(state)
	env.Setup(state)
	fileio.Setup(state)
	ini.Setup(state)
	math.Setup(state)
	stack.Setup(state)
	str.Setup(state)
	
	GlobalNCAState = state
}

// Add commands here
