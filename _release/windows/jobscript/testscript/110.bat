@echo off

if "%1" neq "123" (
  echo Invalid Argument. - %1
  exit 12
)
exit
