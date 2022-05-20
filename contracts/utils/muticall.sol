// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MutiCall {
    function multiCall(address[] calldata targets, bytes[] calldata data) external returns (bytes[] memory) {
        require(targets.length == data.length, "Not match length with target and data");
        
        bytes[] memory results = new bytes[](data.length);
        for (uint i=0; i<targets.length; i++) {
            (bool success, bytes memory result) = targets[i].call(data[i]);
            require(success, "call target failed");
            results[i] = result;
        }

        return results;
    }
}
