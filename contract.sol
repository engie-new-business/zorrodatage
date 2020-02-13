pragma solidity ^0.6.1;

contract Contract {

    mapping(bytes32 => uint) public fingerprints;

    event Registered(address from, bytes32 shasum);

    function register(bytes32 shasum) public {
        require(fingerprints[shasum] == 0, "shasum already registered");
        fingerprints[shasum] = now;
        emit Registered(msg.sender, shasum);
    }

    function lookup(bytes32 shasum) public view returns (uint) {
        return fingerprints[shasum];
    }
}