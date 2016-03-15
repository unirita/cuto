@echo off

setlocal

cd /d "%~dp0"
cd ..
set CUTOROOT=%CD%
set PATH=%CUTOROOT%\bin;%PATH%

cscript /nologo "%CUTOROOT%\bat\extmaster.vbs"
if %errorlevel% neq 0 goto err

cd data
if exist bk_jobnet.csv del bk_jobnet.csv
if exist bk_job.csv del bk_job.csv

sqlite3.exe cuto.sqlite < dbinit.query
if %errorlevel% neq 0 (
  echo Database initialize NG.
) else (
  echo Database initialize OK.
)

exit %errorlevel%

:err
echo master running now.
exit 1
