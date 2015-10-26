Option Explicit
Dim objParm
Set objParm = Wscript.Arguments
WScript.StdErr.WriteLine "Argument1=" & objParm(0)
