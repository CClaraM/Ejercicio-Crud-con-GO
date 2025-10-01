@echo off
REM Cambia esta ruta según donde tengas instalado MSYS2
SET MSYS2_GCC_PATH=C:\msys64\ucrt64\bin\

REM Obtiene el PATH actual del usuario
FOR /F "tokens=* USEBACKQ" %%i IN (`echo %PATH%`) DO SET CURRENT_PATH=%%i

REM Agrega la ruta de GCC al PATH permanentemente para el usuario
setx PATH "%MSYS2_GCC_PATH%;%CURRENT_PATH%"

echo Ruta de GCC agregada permanentemente al PATH del usuario.
echo Abre un nuevo CMD y verifica con:
echo gcc --version

pause