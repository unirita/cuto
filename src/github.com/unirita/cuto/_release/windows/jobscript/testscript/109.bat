@echo off
cd /d "%~dp0"
echo %TEST05%
if "%TEST05%" neq "‚ ,‚¢,‚¤" (
  echo Invalid TEST05.
  exit /b 12
)
exit /b
