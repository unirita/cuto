@echo off

setlocal

del "%GOPATH%\src\cuto\master\master.exe"
del "%GOPATH%\src\cuto\servant\servant.exe"
del "%GOPATH%\src\cuto\show\show.exe"

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
cd /d "%GOPATH%\src\cuto\master"
go build
if %errorlevel% neq 0 goto err

echo servant building...
cd /d "%GOPATH%\src\cuto\servant"
go build
if %errorlevel% neq 0 goto err

echo show utility building...
cd /d "%GOPATH%\src\cuto\show"
go build
if %errorlevel% neq 0 goto err

exit /b 0

:err
echo Failed go build.
exit /b 1
