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

package main

import "strings"
import "dctech/raptor"

func SetupBuiltins() {
	NewNativeTemplate("!TEMPLATE", tempTemplate)
	NewNativeTemplate("!SCRIPT_TEMPLATE", tempScriptTemplate)

	NewNativeTemplate("!SCRIPT", tempScript)
	NewNativeTemplate("SCRIPT", tempScript)
	NewNativeTemplate("#SCRIPT", tempScript)
}

func tempTemplate(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to !TEMPLATE.")
	}

	name := params[0]
	text := params[len(params)-1]
	paramNames := params[1 : len(params)-1]

	parsedParams := make([]*TemplateParam, 0, len(paramNames))

	for _, val := range paramNames {
		rtn := new(TemplateParam)
		if strings.Contains(val, "=") {
			parts := strings.SplitN(val, "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val
		parsedParams = append(parsedParams, rtn)
	}

	NewUserTemplate(name, text, parsedParams)

	return ""
}

func tempScriptTemplate(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to !SCRIPT_TEMPLATE.")
	}

	name := params[0]
	text := params[len(params)-1]
	paramNames := params[1 : len(params)-1]

	parsedParams := make([]*TemplateParam, 0, len(paramNames))

	for _, val := range paramNames {
		rtn := new(TemplateParam)
		if strings.Contains(val, "=") {
			parts := strings.SplitN(val, "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val
		parsedParams = append(parsedParams, rtn)
	}

	NewScriptTemplate(name, text, parsedParams)

	return ""
}

func tempScript(params []string) string {
	if len(params) < 1 {
		panic("Wrong number of params to SCRIPT.")
	}

	GlobalRaptorState.Code.Add(params[0])
	GlobalRaptorState.Envs.Add(raptor.NewEnvironment())

	if len(params) > 1 {
		GlobalRaptorState.AddParams(params[1:]...)
	}

	rtn, err := GlobalRaptorState.Run()
	if err != nil {
		panic("Script Error: " + err.Error())
	}

	GlobalRaptorState.Envs.Remove()

	if rtn == nil {
		return ""
	}
	return rtn.String()
}
