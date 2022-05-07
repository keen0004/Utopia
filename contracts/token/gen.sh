#!/bin/bash

solc ./IERC20.sol --abi --optimize --overwrite --output-dir ./ >/dev/null
abigen --abi ./IERC20.abi --pkg token --type ERC20 --out ./erc20.go

solc ./IERC721.sol --abi --optimize --overwrite --output-dir ./ >/dev/null
abigen --abi ./IERC721.abi --pkg token --type ERC721 --out ./erc721.go
