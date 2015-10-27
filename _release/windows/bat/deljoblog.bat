@echo off

setlocal
rem *** 残すジョブログの日数 ***
set GENERAL_NUM=30
rem ****************************

cd /d "%~dp0"
cd ..
set CUTOROOT=%CD%

cscript /nologo %CUTOROOT%\bat\deletedir.vbs %CUTOROOT%\joblog %GENERAL_NUM%

exit /b %errorlevel%
