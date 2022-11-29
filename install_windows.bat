@echo off
::  prompt user to get elevated priviledge for installation
::-------------------------------------
REM  --> Check for permissions
>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"

REM --> If error flag set, we do not have admin.
if '%errorlevel%' NEQ '0' (
    echo Requesting administrative privileges...
    goto UACPrompt
) else ( goto gotAdmin )

:UACPrompt
    echo Set UAC = CreateObject^("Shell.Application"^) > "%temp%\getadmin.vbs"
    set params = %*:"="
    echo UAC.ShellExecute "cmd.exe", "/c %~s0 %params%", "", "runas", 1 >> "%temp%\getadmin.vbs"

    "%temp%\getadmin.vbs"
    del "%temp%\getadmin.vbs"
    exit /B

:gotAdmin
    pushd "%CD%"
    CD /D "%~dp0"
::--------------------------------------

::ENTER YOUR CODE BELOW:

:: Delete windows service service if it exists and create new one
SC QUERY WOLServerService > NUL
IF ERRORLEVEL 1060 GOTO MISSING
ECHO service exist stopping then deleting service
sc stop WOLServerService 
sc delete WOLServerService 
GOTO NOSERVICE

:MISSING
ECHO service does not exist creating new one...

:NOSERVICE

:: kill any existing client
QPROCESS "wolwebservice.exe">NUL 2> nul
IF %ERRORLEVEL% EQU 0 taskkill /IM wolservice.exe


timeout 2

::Copy binaries to target destination
::Check of C:\wolservice\ folder exists if not create it

IF EXIST C:\wolservice\ GOTO FOLDEREXISTS
MKDIR C:\wolservice\
GOTO INSTALL
:FOLDEREXISTS
::Delete old wolwebservice.exe
del C:\wolservice\wolwebservice.exe
del C:\wolservice\efiDLL.dll
:INSTALL
echo f | xcopy /f /y build\dist\wolwebservice c:\wolservice\wolwebservice.exe
echo f | xcopy /f /y build\dist\efiDLL.dll c:\wolservice\

echo "Creating windows service"
::Create windows service
sc create WOLServerService binPath="C:\wolservice\wolwebservice.exe"
sc config WOLServerService start= auto
sc start WOLServerService

timeout 5