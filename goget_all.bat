@echo off

set GOPATH=%~dp0

go get github.com/coopernurse/gorp
if %errorlevel% neq 0 goto err

go get github.com/BurntSushi/toml
if %errorlevel% neq 0 goto err

go get github.com/cihub/seelog
if %errorlevel% neq 0 goto err

go get github.com/mattn/go-sqlite3
if %errorlevel% neq 0 goto err

exit /b 0

:err
echo Failed go get.
exit /b 1
