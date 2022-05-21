// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Transfer is Ownable {
    mapping (bytes32 => bool) private _txlist;

    event etransfer(address token, address to, uint256 value);
    event ewithdraw(address token, address to, uint256 value, bytes salt);
    event ediscard(address token, address to, uint256 value, bytes salt);

    receive() external payable {}

    function getBalance() view external returns(uint256) {
        return address(this).balance;
    }

    function batchTransfer(address[] memory tolist, uint256[] memory values) external onlyOwner {
        require(tolist.length == values.length, "tolist.length != values.length");

        uint256 balance = address(this).balance;
        uint256 need = 0;
        for (uint256 i = 0; i < values.length; i++) {
           need += values[i];
        }

        require(balance >= need, "not enough balance");
        for (uint256 i = 0; i < tolist.length; i++) {
            payable(tolist[i]).transfer(values[i]);

            emit etransfer(address(0x0), tolist[i], values[i]);
        }
    }

    function batchTransferToken(address token, address[] memory tolist, uint256[] memory values) external onlyOwner {
        require(tolist.length == values.length, "tolist.length != values.length");

        IERC20 c = IERC20(token);
        uint256 balance = c.balanceOf(address(this));

        uint256 need = 0;
        for (uint256 i = 0; i < values.length; i++) {
           need += values[i];
        }

        require(balance >= need, "not enough balance");
        for (uint256 i = 0; i < tolist.length; i++) {
            c.transfer(tolist[i], values[i]);

            emit etransfer(token, tolist[i], values[i]);
        }
    }

    function withdraw(address token, address to, uint256 value, bytes memory salt, bytes memory signature) external {
        bytes32 hash = keccak256(abi.encode(token, to, value, salt));
        require(_txlist[hash] == false, "transaction has exist");
        require(ECDSA.recover(hash, signature) == owner(), "not owner signed transaction");

        _txlist[hash] = true;
        if (token == address(0x0)) {
            uint256 balance = address(this).balance;
            require(balance >= value, "not enough balance");

            payable(to).transfer(value);
        } else {
            IERC20 c = IERC20(token);
            uint256 balance = c.balanceOf(address(this));
            require(balance >= value, "not enough balance");

            c.transfer(to, value);
        }

        emit ewithdraw(token, to, value, salt);
    }

    function discard(address token, address to, uint256 value, bytes memory salt, bytes memory signature) external {
        bytes32 hash = keccak256(abi.encode(token, to, value, salt));
        require(ECDSA.recover(hash, signature) == owner(), "not owner signed transaction");

        _txlist[hash] = true;
        emit ediscard(token, to, value, salt);
    }
}
