// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MyToken is ERC20 {
    uint256 public constant MAX_SUPPLY = 100_000_000 * 10 ** 18;
    uint256 public constant PRICE = 0.01 ether;

    constructor(string memory name, string memory symbol) ERC20(name, symbol) {

    }

    function mint(uint256 amount) external payable {
        require(totalSupply() + amount <= MAX_SUPPLY, "not enough tokens");

        uint256 fee = (amount / (10 ** 18)) * PRICE;
        require(msg.value >= fee, "not enough fee");

        _mint(msg.sender, amount);
        payable(msg.sender).transfer(msg.value - fee);
    }
}
