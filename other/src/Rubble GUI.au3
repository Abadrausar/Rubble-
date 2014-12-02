#AutoIt3Wrapper_Icon=.\rubble\Rubble.ico
#AutoIt3Wrapper_UseX64=n
#AutoIt3Wrapper_Res_Description=Rubble GUI Launcher
#AutoIt3Wrapper_Res_Fileversion=1.0
#AutoIt3Wrapper_Res_LegalCopyright=Copyright 2013 Milo Christiansen
#AutoIt3Wrapper_Res_Language=1033

#include <ButtonConstants.au3>
#include <GUIConstantsEx.au3>
#include <TreeViewConstants.au3>
#include <WindowsConstants.au3>

#include <File.au3>

Opt("TrayMenuMode",1)
Opt("TrayIconHide",1)

Dim $AddonList[1]
$AddonList[0] = 0
Dim $AddonItemList[1]
$AddonItemList[0] = 0

Dim $Main = GUICreate("Rubble GUI", 466, 239, -1, -1)
Dim $hAddons = GUICtrlCreateTreeView(5, 5, 300, 200, BitOR($GUI_SS_DEFAULT_TREEVIEW,$TVS_CHECKBOXES), $WS_EX_CLIENTEDGE)
Dim $hRun = GUICtrlCreateButton("Run Rubble!", 310, 5, 150, 50)
Dim $hOtherOptions = GUICtrlCreateInput("", 5, 210, 455, 21)
GUISetState(@SW_SHOW)

_LoadAddons()
_Init()

Dim $msg
While 1
	$msg = GUIGetMsg()
	Switch $msg
		Case $GUI_EVENT_CLOSE
			_Exit()
		Case $hRun
			Local $addons = ""
			For $X = 1 To $AddonItemList[0]
				If BitAnd(GUICtrlRead($AddonItemList[$X]),$GUI_CHECKED) Then
					$addons &= ";" & $AddonList[$X]
				EndIf
			Next
			$addons = StringTrimLeft($addons,1)
			Run(@ComSpec & ' /k rubble.exe -addons="' & $addons & '" ' & GUICtrlRead($hOtherOptions))
			_Exit()
	EndSwitch
WEnd

Func _LoadAddons()
	$AddonList = _FileListToArray("./addons", "*", 2)
	If @error Then
		Dim $AddonList[1]
		$AddonList[0] = 0
		Dim $AddonItemList[1]
		$AddonItemList[0] = 0
		Return
	EndIf

	ReDim $AddonItemList[$AddonList[0] + 1]
	$AddonItemList[0] = $AddonList[0]
	For $X = 1 To $AddonList[0]
		$AddonItemList[$X] = GUICtrlCreateTreeViewItem($AddonList[$X], $hAddons)
	Next
EndFunc

Func _Init()
	If FileExists(@ScriptDir & "\gui.ini") Then
		Local $count = 0

		While 1
			Local $val = IniRead(@ScriptDir & "\gui.ini", "addons", $count, "ERROR")
			If $val == "ERROR" Then
				ExitLoop
			EndIf

			For $X = 1 To $AddonList[0]
				If $AddonList[$X] == $val Then
					GUICtrlSetState($AddonItemList[$X], $GUI_CHECKED)
					ExitLoop
				EndIf
			Next
			$count += 1
		WEnd

		GUICtrlSetData($hOtherOptions, IniRead(@ScriptDir & "\gui.ini", "general", "otheroptions", ""))
	EndIf
EndFunc

Func _Exit()
	FileClose(FileOpen(@ScriptDir & "\gui.ini", 2))
	Local $count = 0
	For $X = 1 To $AddonItemList[0]
		If BitAnd(GUICtrlRead($AddonItemList[$X]),$GUI_CHECKED) Then
			IniWrite(@ScriptDir & "\gui.ini", "addons", $count, $AddonList[$X])
			$count += 1
		EndIf
	Next
	IniWrite(@ScriptDir & "\gui.ini", "general", "otheroptions", GUICtrlRead($hOtherOptions))
	Exit
EndFunc
