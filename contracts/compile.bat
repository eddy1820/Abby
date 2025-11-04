@echo off

REM 創建 build 目錄（如果不存在）
if not exist "build" mkdir build

REM 編譯 Solidity 合約
solcjs --abi --bin SimpleStorage.sol -o build/

REM 重命名生成的文件
cd build
del SimpleStorage.abi SimpleStorage.bin 2>nul
rename "SimpleStorage_sol_SimpleStorage.abi" "SimpleStorage.abi"
rename "SimpleStorage_sol_SimpleStorage.bin" "SimpleStorage.bin"
cd ..

REM 生成 Go 綁定代碼
C:\Users\Eddy\go\bin\abigen --bin=build/SimpleStorage.bin --abi=build/SimpleStorage.abi --pkg=contracts --out=SimpleStorage.go
type SimpleStorage.go
