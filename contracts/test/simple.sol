// SPDX-License-Indentifier: MIT
pragma solidity ^0.8.0;

contract simple {
    string public data;

    constructor(string memory msg) public {
        data = msg;
    }

    function SetMessage(string memory msg) public {
        data = msg;
    }

    function GetMessage() view public returns(string memory) {
        return data;
    }
}
