@if(0)==(0) ECHO OFF
cscript.exe //nologo //E:JScript "%~f0" %* %1
rem goto :EOF
exit /b %errorlevel%
@end

var args = WScript.arguments;
if (args.length < 1) {
  WScript.sleep(1000)
} else {
  WScript.sleep(args(0)*1000)
}
