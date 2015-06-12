' processNameに指定したプロセス名が存在する場合は1を返す。
' 存在しない場合は0を返す。
processName = "master.exe"
If GetProcessId(processName) <> 0 Then
  WScript.Quit(1)
Else
  WScript.Quit(0)
End If

'-------------------------------------------------------------------------------
' 指定されたプロセス名のIDを取得する
Function GetProcessId(ProcessName)
    Dim Service,QfeSet,Qfe,r
    
    Set Service = WScript.CreateObject("WbemScripting.SWbemLocator").ConnectServer
    Set QfeSet = Service.ExecQuery("Select * From Win32_Process Where Caption='" & ProcessName & "'")
    
    r = 0
    For Each Qfe In QfeSet
        r = Qfe.ProcessId
        Exit For
    Next

    GetProcessId = r
End Function
