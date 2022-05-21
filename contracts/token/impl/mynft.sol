// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// address: 0x3320325420e04E52C47Df416529b2680aaeD7c4A
contract MyNFT is ERC721Enumerable, Ownable {
    uint256 public constant MAX_SUPPLY = 10000;
    uint256 public constant PRICE = 0.1 ether;

    string private baseUri;

    constructor(string memory name, string memory symbol) ERC721(name, symbol) {

    }

    function mint(uint256 amount) external payable {
        require(totalSupply() + amount <= MAX_SUPPLY, "not enough tokens");

        uint256 fee = amount * PRICE;
        require(msg.value >= fee, "not enough fee");

        uint256 tokenId = totalSupply();

        for (uint i=0; i<amount; i++) {
            _safeMint(msg.sender, tokenId + i + 1);
        }
        
        payable(msg.sender).transfer(msg.value - fee);
    }

    function setBaseURI(string memory url) external onlyOwner {
        baseUri = url;
    }

    function _baseURI() internal view virtual override returns (string memory) {
        return baseUri;
    }
}
