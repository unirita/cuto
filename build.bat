@echo off

setlocal

del "%GOPATH%\bin\master.exe"
del "%GOPATH%\bin\servant.exe"
del "%GOPATH%\bin\show.exe"
del "%GOPATH%\bin\flowgen.exe"

if "%1" neq "" goto BUILD
:UNIT_TEST
rem *****************
rem Unit test
rem *****************
cd /d %GOPATH%\src\cuto
call test_all_cover.bat x
if %errorlevel% neq 0 goto err

:BUILD
rem *****************
rem All build
rem *****************
echo master building...
go install cuto/master
if %errorlevel% neq 0 goto err

echo servant building...
go install cuto/servant
if %errorlevel% neq 0 goto err

echo show utility building...
go install cuto/show
if %errorlevel% neq 0 goto err

echo flowgen utility building...
go install cuto/flowgen
if %errorlevel% neq 0 goto err

exit /b 0

:err
echo Failed go build.
exit /b 1
