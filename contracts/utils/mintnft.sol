// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC721 {
    function totalSupply() external view returns(uint256);
    function mint(uint256 amount) external payable;
    function transferFrom(address from, address to, uint256 tokenId) external;
}

contract ERC721Mint {
    constructor(address ERC721, address owner) payable {
        uint256 t = IERC721(ERC721).totalSupply();
        IERC721(ERC721).mint{value: 0.05 ether}(5);

        for (uint i=0; i<5; i++) {
            IERC721(ERC721).transferFrom(address(this), owner, t+i);
        }
        
        selfdestruct(payable(owner));
    }
}

contract MintFactory {
    address owner;

    constructor() {
        owner = msg.sender;
    }

    function deploy(address ERC721, uint count) external payable {
        for (uint i=0; i<count; i++) {
            new ERC721Mint{value: 0.05 ether}(ERC721, owner);
        }
    }
}

