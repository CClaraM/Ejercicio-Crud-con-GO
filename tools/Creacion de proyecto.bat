@echo off
echo Antes de crear el proyecto se debe instalar la libreria MSYS2 y actualizar GCC
echo esto puedes hacerlo descargandolo desde https://www.msys2.org/, y siguiendo las instrucciones
echo presione una tecla para continuar
pause
echo Se crear proyecto
mkdir go_crud_fyne
cd go_crud_fyne
go mod init go_crud_fyne

echo Se instala dependencias
go get fyne.io/fyne/v2@latest
go get github.com/go-sql-driver/mysql@latest

echo se forza instalacion de dependencias
go get github.com/go-gl/gl/v2.1/gl
go get github.com/go-gl/glfw/v3.3/glfw

echo Se organiza el proyecto
go mod tidy