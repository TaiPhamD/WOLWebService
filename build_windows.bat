@SETLOCAL
@REM alternative command to generate for vs2019:cd build; cmake .. -G "Visual Studio 16 2019" -A x64
@REM delete output folder 

@REM add msys build tools to path
SET PATH=C:\msys64\mingw64\bin;C:\msys64\usr\bin;%PATH%
if exist build rd /s /q build
mkdir -p build

@REM Building util_windows.cpp into a static lib that's compatible with CGO
g++ -O2 -std=c++17 -Wall -c -march=x86-64 server\util\util_windows.cpp -o build\util_windows.o -static-libgcc -static-libstdc++
ar rc build\libutil_windows.a build\util_windows.o
