@echo off
cd /d "%~dp0"
call sleep.bat

echo %1
echo %2
echo %3
if "%1" neq "abc" (
  echo Invalid Argument#1 - %1
  exit /b 101
)
if %2 neq "e f g" (
  echo Invalid Argument#2 - %2
  exit /b 102
)
if "%3" neq "h" (
  echo Invalid Argument#3 - %3
  exit /b 103
)

exit /b 0
