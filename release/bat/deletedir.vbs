' Argument#1 削除対象のフォルダ
' Argument#2 残すフォルダ数
Set objParm = Wscript.Arguments
If objParm.Count <> 2 Then
  WScript.Quit(1)
End If

Set objFSO = WScript.CreateObject("Scripting.FileSystemObject")
Set objFolder = objFSO.GetFolder(objParm(0))
Set subfolders = objFolder.subfolders

lngCount = subfolders.Count - objParm(1)
If lngCount < 1 Then
  WScript.Quit(0)
End If

counter = 1
For Each subfolder in objFolder.subfolders
  objFSO.DeleteFolder subfolder, True
  counter = counter + 1
  If counter > lngCount Then
    Exit For
  End If
Next

WScript.Quit(0)
