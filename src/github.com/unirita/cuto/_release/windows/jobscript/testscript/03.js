var objNamed = WScript.Arguments.Named

if (objNamed.Length != 2) {
  WScript.Echo("[err] Invalid Argument number = "+objNamed.Length);
  WScript.Quit(1);
}

if (objNamed.Item("PREFIX") != "abc") {
  WScript.Echo("[err] Invalid Argument PREFIX = "+objNamed.Item("PREFIX"));
  WScript.Quit(2);
}
if (objNamed.Item("SUFFIX") != "xyz") {
  WScript.Echo("[err] Invalid Argument SUFFIX = "+objNamed.Item("SUFFIX"));
  WScript.Quit(3);
}

WScript.Quit(0);
