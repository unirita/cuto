@echo off

set CUR=%CD%
cd /d "%~dp0"
cd ..\..
set VERCUR=%CD%

if "%1" equ "" (
  echo Nothing ID
  exit 101
)
echo ID[ "%1" ]
if "%2" equ "" (
  echo Nothing start-date
  exit 102
)
echo start[ "%2" ]

rem ### master environment %windir% ###
if "%TESTENV1%" neq "%windir%" (
  echo Invalid TESTENV1[ %TESTENV1% ]
  exit 103
)
rem ### servant environment %OS% ###
if "%TESTENV2%" neq "%OS%" (
  echo Invalid TESTENV2[ %TESTENV2% ]
  exit 104
)

if "%VERCUR%" neq "%CUR%" (
  echo Invalid Current-dir[ %CD% ]
  exit 105
)

exit 0
