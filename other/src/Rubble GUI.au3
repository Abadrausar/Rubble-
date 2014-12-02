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
#include <Array.au3>

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
GUICtrlSetData($hOtherOptions, IniRead(@ScriptDir & "\gui.ini", "general", "otheroptions", ""))

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
			Run(@ComSpec & ' /k rubble.exe -addons="' & $addons & '" ' & GUICtrlRead($hOtherOptions))
			IniWrite(@ScriptDir & "\gui.ini", "general", "otheroptions", GUICtrlRead($hOtherOptions))
			Exit
	EndSwitch
WEnd

Func _LoadAddons()
	RunWait('rubble.exe -addonlist', @ScriptDir, @SW_HIDE)

	If Not FileExists(@ScriptDir & "\addons\addonlist.ini") Then
		MsgBox(16, "Loading Error", "Rubble failed to generate/update the addonlist.ini, are you sure you set this thing up right?")
		Exit 1
	EndIf

	Dim $AddonList[1]
	$AddonList[0] = 0
	Dim $AddonItemList[1]
	$AddonItemList[0] = 0

	Local $section = IniReadSection(@ScriptDir & "\addons\addonlist.ini", "addons")
	For $X = 1 To $section[0][0]
		_ArrayAdd($AddonList, $section[$X][0])
		_ArrayAdd($AddonItemList, GUICtrlCreateTreeViewItem($section[$X][0], $hAddons))
		$AddonList[0] += 1
		$AddonItemList[0] +=1

		If $section[$X][1] == "true" Then
			GUICtrlSetState($AddonItemList[$AddonItemList[0]], $GUI_CHECKED)
		EndIf
	Next
EndFunc
