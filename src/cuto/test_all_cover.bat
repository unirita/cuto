@echo off

set LOGFILE="%~dp0cover_all.txt"
set RETCODE=0

cd /d "%~dp0"

pushd console
echo console package tested...
go test -coverprofile cover.out> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd db
echo db package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd query
echo db/query package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd tx
echo db/tx package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd log
echo log package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd master
echo master package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd config
echo master/config package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd jobnet
echo master/jobnet package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd parser
echo master/jobnet/parser package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd
pushd remote
echo master/remote package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd message
echo message package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd servant
echo servant package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd config
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd job
echo servant/job package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd remote
echo servant/remote package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd show
echo show package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd gen
echo show/gen package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd util
echo util package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

if "%1" neq "" type %LOGFILE%
if "%1" equ "" pause

exit /b %RETCODE%
