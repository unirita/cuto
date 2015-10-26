Option Explicit
Dim objParm
Set objParm = Wscript.Arguments
WScript.Echo "Argument1=" & objParm(0)
WScript.Echo "Argument2=" & objParm(1)
