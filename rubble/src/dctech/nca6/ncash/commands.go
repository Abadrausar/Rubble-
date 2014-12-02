// NCA v6 Shell Command.
package ncash

import "fmt"
import "dctech/nca6"
import "os"
import "io"

// Adds the ncash command to the state. The ncash command gives access to a 
// shell much like NCASH6 that may be broken into at any time.
func Setup(state *nca6.State) {
	state.NewNativeCommand("ncash", CommandNCASH)
}

// Break into the NCA SHell. VERY useful for debuging!
// This command will provide an interactive shell until eof is reached, 
// on windows this may be simulated by pressing CTRL+Z followed by <ENTER>.
// 	ncash
// Returns the return value of the last command to be run.
func CommandNCASH(state *nca6.State, params []*nca6.Value) {
	line := make([]byte, 0, 100)
	curchar := make([]byte, 1, 1)
	escape := false

	fmt.Println("=BEGIN=NCASH=================================")
	fmt.Print(">>>")
	for {
		_, err := os.Stdin.Read(curchar)
		if err == io.EOF {
			fmt.Println("Exiting...")
			break
		} else if err != nil {
			fmt.Println("Read Error:", err, "\nExiting...")
			break
		}

		if curchar[0] == byte('\r') {
			continue
		}

		if curchar[0] == byte('\\') && !escape {
			escape = true
			continue
		}

		if curchar[0] == byte('\n') && !escape {
			state.Code.Add(string(line))
			rtn, err := state.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
			fmt.Println("Ret:", rtn)
			
			line = line[:0]
			fmt.Print(">>>")
			continue
		}

		if curchar[0] != byte('\n') && escape {
			line = append(line, byte('\\'))
		}
		
		if curchar[0] == byte('\n') && escape {
			fmt.Print(">>>")
		}

		escape = false
		line = append(line, curchar...)
	}
	fmt.Println("=END=NCASH===================================")
}
