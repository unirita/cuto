@echo off

if "%1" neq "123" (
  echo Invalid Argument. - %1
  exit /b 12
)
exit /b
