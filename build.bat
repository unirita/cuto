@echo off

setlocal
set GOPATH=%~dp0

del "%GOPATH%src\cuto\master\master.exe"
del "%GOPATH%src\cuto\servant\servant.exe"
del "%GOPATH%src\cuto\show\show.exe"

echo master building...
cd /d "%GOPATH%src\cuto\master"
go build
if %errorlevel% neq 0 goto err

echo servant building...
cd /d "%GOPATH%src\cuto\servant"
go build
if %errorlevel% neq 0 goto err

echo show utility building...
cd /d "%GOPATH%src\cuto\show"
go build
if %errorlevel% neq 0 goto err

exit /b 0

:err
echo Failed go build.
exit /b 1
