@echo off

set LOGFILE="%~dp0cover_all_win.txt"
set RETCODE=0

cd /d "%~dp0"

pushd console
echo github.com/unirita/cuto/console package tested...
go test -coverprofile cover.out> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd db
echo github.com/unirita/cuto/db package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd query
echo github.com/unirita/cuto/db/query package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd tx
echo github.com/unirita/cuto/db/tx package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd log
echo github.com/unirita/cuto/log package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd master
echo github.com/unirita/cuto/master package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd config
echo github.com/unirita/cuto/master/config package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd jobnet
echo github.com/unirita/cuto/master/jobnet package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd parser
echo github.com/unirita/cuto/master/jobnet/parser package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd
pushd remote
echo github.com/unirita/cuto/master/remote package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd message
echo github.com/unirita/cuto/message package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd servant
echo github.com/unirita/cuto/servant package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd config
echo github.com/unirita/cuto/servant/config package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd job
echo github.com/unirita/cuto/servant/job package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
pushd remote
echo github.com/unirita/cuto/servant/remote package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd show
echo github.com/unirita/cuto/show package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
pushd gen
echo github.com/unirita/cuto/show/gen package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd util
echo github.com/unirita/cuto/util package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

pushd flowgen
pushd converter
echo github.com/unirita/cuto/flowgen/converter package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd realtime
pushd network
echo github.com/unirita/cuto/realtime/network package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd
popd

pushd utctime
echo github.com/unirita/cuto/utctime package tested...
go test -coverprofile cover.out>> %LOGFILE%
if %errorlevel% neq 0 (
  echo NG.
  set RETCODE=1
)
popd

if "%1" neq "" type %LOGFILE%
if "%1" equ "" pause

exit %RETCODE%
