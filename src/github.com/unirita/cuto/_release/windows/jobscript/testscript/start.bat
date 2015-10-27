@echo off

set CUR=%CD%
cd /d "%~dp0"
cd ..\..
set VERCUR=%CD%

if "%1" equ "" (
  echo Nothing ID
  exit /b 101
)
echo ID[ "%1" ]
if "%2" equ "" (
  echo Nothing start-date
  exit /b 102
)
echo start[ "%2" ]

rem ### master environment %windir% ###
if "%TESTENV1%" neq "%windir%" (
  echo Invalid TESTENV1[ %TESTENV1% ]
  exit /b 103
)
rem ### servant environment %OS% ###
if "%TESTENV2%" neq "%OS%" (
  echo Invalid TESTENV2[ %TESTENV2% ]
  exit /b 104
)

if "%VERCUR%" neq "%CUR%" (
  echo Invalid Current-dir[ %CD% ]
  exit /b 105
)

exit /b 0
