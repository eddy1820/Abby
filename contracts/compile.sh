#!/bin/bash

# 編譯 Solidity 合約
solc --abi --bin contracts/SimpleStorage.sol -o build/

# 生成 Go 綁定代碼
abigen --bin=build/SimpleStorage.bin --abi=build/SimpleStorage.abi --pkg=contracts --out=contracts/SimpleStorage.go
