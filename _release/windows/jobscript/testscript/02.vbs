Option Explicit
On Error Resume Next

Dim objWshShell
Dim objRc
Dim objSd
Dim objEd
Dim objOut

Set objWshShell = WScript.CreateObject("WScript.Shell")
objRc = objWshShell.ExpandEnvironmentStrings("%RC%")
objSd = objWshShell.ExpandEnvironmentStrings("%SD%")
objEd = objWshShell.ExpandEnvironmentStrings("%ED%")
objOut = objWshShell.ExpandEnvironmentStrings("%OUT%")

If objRc <> 0 Then
	WScript.Echo "[err] Invalid RC : " & objRc
Else
	WScript.Echo objRc
End If

If Len(objSd) = 0 Then
	Wscript.Echo "[err] Invalid SD : " & objSd
Else
	WScript.Echo objSd
End If

If Len(objEd) = 0 Then
	Wscript.Echo "[err] Invalid ED : " & objEd
Else
	WScript.Echo objEd
End If

If Len(objOut) <> 0 Then
	Wscript.Echo "[err] Invalid OUT : " & objOut
Else
	WScript.Echo objOut
End If

Set objWshShell = Nothing

Wscript.Quit 0
