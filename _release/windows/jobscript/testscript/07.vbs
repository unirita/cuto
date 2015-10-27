Set objWshShell = WScript.CreateObject("WScript.Shell")
objPath = objWshShell.ExpandEnvironmentStrings("%PATH%")
WScript.Echo objPath
lngPos = InStr(1, objPath, "C:\cuto07", 1)
If lngPos = 0 Then
	WScript.Echo "[err] Invalid path."
	Wscript.Quit 1
End If

Set objParm = WScript.Arguments
If objParm.Count <> 1 Then
	WScript.Echo "[err] Invalid Argument Count. - " & objParm.Count
	Wscript.Quit 2
End If

If StrComp(objParm.Item(0), "a b", 1) <> 0 Then
	WScript.Echo "[err] Invalid Argument#1 - " & objParm.Item(0)
	Wscript.Quit 3
End If

Wscript.Quit 0
 