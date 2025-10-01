@echo off
REM ============================================
REM Configuración de variables de entorno para Go
REM ============================================

setx GOROOT "C:\Program Files\Go" /M
setx GOPATH "%USERPROFILE%\go" /M
setx PATH "%PATH%;C:\Program Files\Go\bin;%USERPROFILE%\go\bin" /M

echo ============================================
echo Variables de entorno configuradas:
echo GOROOT = C:\Program Files\Go
echo GOPATH = %USERPROFILE%\go
echo PATH   = actualizado con Go
echo ============================================
echo.
echo Reinicia la terminal o la PC para aplicar los cambios.
pause