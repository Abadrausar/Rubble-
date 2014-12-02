#AutoIt3Wrapper_Icon=.\rubble\Rubble.ico
#AutoIt3Wrapper_UseX64=n
#AutoIt3Wrapper_Res_Description=Rubble GUI Launcher
#AutoIt3Wrapper_Res_Fileversion=1.0
#AutoIt3Wrapper_Res_LegalCopyright=Copyright 2013-2014 Milo Christiansen
#AutoIt3Wrapper_Res_Language=1033
#AutoIt3Wrapper_Res_Comment=Rubble Version (at last GUI build): 3.11

#include <ButtonConstants.au3>
#include <GUIConstantsEx.au3>
#include <TreeViewConstants.au3>
#include <WindowsConstants.au3>

#include <File.au3>
#include <Array.au3>

Opt("TrayMenuMode",1)
Opt("TrayIconHide",1)

Dim $AddonList[1]
$AddonList[0] = 0
Dim $AddonItemList[1]
$AddonItemList[0] = 0

Global $Main = GUICreate("Rubble GUI", 437, 459, -1, -1, BitOR($GUI_SS_DEFAULT_GUI,$WS_SIZEBOX,$WS_THICKFRAME))
Global $hAddons = GUICtrlCreateTreeView(8, 8, 420, 352, BitOR($GUI_SS_DEFAULT_TREEVIEW,$TVS_CHECKBOXES), $WS_EX_CLIENTEDGE)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKTOP+$GUI_DOCKBOTTOM)
Global $hRun = GUICtrlCreateButton("Run Rubble!", 8, 416, 422, 34)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKBOTTOM+$GUI_DOCKHEIGHT)
Global $hOtherOptions = GUICtrlCreateInput("", 8, 392, 423, 21)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKBOTTOM+$GUI_DOCKHEIGHT)
GUICtrlCreateLabel("Extra Rubble Commandline Options", 8, 368, 418, 17)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKBOTTOM+$GUI_DOCKHEIGHT)

Dim $Rubble = IniRead(@ScriptDir & "\gui.ini", "general", "rubble", ".")
Dim $AddonDir = $Rubble & "\addons"
If FileExists($Rubble & "\rubble.ini") Then
	$AddonDir = IniRead($Rubble & "\rubble.ini", "rubble", "addonsdir", $AddonDir)
EndIf

RunWait($Rubble & "/rubble.exe" & ' -addonlist', $Rubble , @SW_HIDE)

If Not FileExists($AddonDir & "\addonlist.ini") Then
	MsgBox(16, "Loading Error", "Rubble failed to generate/update the addonlist.ini, are you sure you set this thing up right?")
	Exit 1
EndIf

Dim $AddonList[1]
$AddonList[0] = 0
Dim $AddonItemList[1]
$AddonItemList[0] = 0

Local $section = IniReadSection($AddonDir & "\addonlist.ini", "addons")
For $X = 1 To $section[0][0]
	_ArrayAdd($AddonList, $section[$X][0])
	_ArrayAdd($AddonItemList, GUICtrlCreateTreeViewItem($section[$X][0], $hAddons))
	$AddonList[0] += 1
	$AddonItemList[0] +=1

	If $section[$X][1] == "true" Then
		GUICtrlSetState($AddonItemList[$AddonItemList[0]], $GUI_CHECKED)
	EndIf
Next

GUICtrlSetData($hOtherOptions, IniRead(@ScriptDir & "\gui.ini", "general", "otheroptions", ""))

GUISetState(@SW_SHOW)

Dim $msg
While 1
	$msg = GUIGetMsg()
	Switch $msg
		Case $GUI_EVENT_CLOSE
			IniWrite(@ScriptDir & "\gui.ini", "general", "otheroptions", GUICtrlRead($hOtherOptions))
			Exit
		Case $hRun
			Local $addons = ""
			For $X = 1 To $AddonItemList[0]
				If BitAnd(GUICtrlRead($AddonItemList[$X]),$GUI_CHECKED) Then
					$addons &= ";" & $AddonList[$X]
				EndIf
			Next
			$addons = StringTrimLeft($addons,1)
			Run(@ComSpec & ' /k ""./rubble.exe" -addons="' & $addons & '" ' & GUICtrlRead($hOtherOptions) & '"', $Rubble)
			IniWrite(@ScriptDir & "\gui.ini", "general", "otheroptions", GUICtrlRead($hOtherOptions))
			IniWrite(@ScriptDir & "\gui.ini", "general", "rubble", $Rubble)
			Exit
	EndSwitch
WEnd
