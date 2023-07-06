import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

contract EquipsRegistry {

    event Equip(address indexed owner, bytes32 indexed payload);

    constructor(){}

    function addressToString(address _address) internal pure returns (string memory) {
        bytes32 value = bytes32(uint256(uint160(_address)));
        bytes memory alphabet = "0123456789abcdef";
        
        bytes memory str = new bytes(42);
        str[0] = "0";
        str[1] = "x";
        for (uint256 i = 0; i < 20; i++) {
            str[2 + i * 2] = alphabet[uint8(value[i + 12] >> 4)];
            str[3 + i * 2] = alphabet[uint8(value[i + 12] & 0x0f)];
        }
        return string(str);
    }

    function uintToString(uint256 _value) internal pure returns (string memory) {
        if (_value == 0) {
            return "0";
        }

        uint256 temp = _value;
        uint256 digits;

        while (temp != 0) {
            digits++;
            temp /= 10;
        }

        bytes memory buffer = new bytes(digits);

        while (_value != 0) {
            digits--;
            buffer[digits] = bytes1(uint8(48 + (_value % 10)));
            _value /= 10;
        }

        return string(buffer);
    }

    function stringToBytes32WithAddress(address addr, string memory str, string memory salt, string memory index) public pure returns (bytes32) {
        bytes memory addrBytes = abi.encodePacked(addr);
        bytes memory strBytes = bytes(str);
        bytes memory saltBytes = bytes(salt);
        bytes memory indexBytes = bytes(index);
        bytes memory divBytes = bytes(":");
        
        bytes memory combinedBytes = new bytes(addrBytes.length + divBytes.length*3 + strBytes.length + divBytes.length + saltBytes.length + indexBytes.length);
        uint256 idx = 0;
        
        for (uint256 i = 0; i < addrBytes.length; i++) {
            combinedBytes[idx] = addrBytes[i];
            idx++;
        }
        
        for (uint256 i = 0; i < divBytes.length; i++) {
            combinedBytes[idx] = divBytes[i];
            idx++;
        }
        
        for (uint256 i = 0; i < strBytes.length; i++) {
            combinedBytes[idx] = strBytes[i];
            idx++;
        }
        
        for (uint256 i = 0; i < divBytes.length; i++) {
            combinedBytes[idx] = divBytes[i];
            idx++;
        }
        
        for (uint256 i = 0; i < saltBytes.length; i++) {
            combinedBytes[idx] = saltBytes[i];
            idx++;
        }

        for (uint256 i = 0; i < divBytes.length; i++) {
            combinedBytes[idx] = divBytes[i];
            idx++;
        }

        for (uint256 i = 0; i < indexBytes.length; i++) {
            combinedBytes[idx] = indexBytes[i];
            idx++;
        }
        
        bytes32 paddedBytes;
        assembly {
            paddedBytes := mload(add(combinedBytes, 32))
        }
        
        return paddedBytes;
    }

    function concatenate(address addr, uint256 value, uint salt, uint index) public pure returns (bytes32) {
        string memory addrStr = addressToString(addr);
        string memory valueStr = uintToString(value);
        string memory saltStr = uintToString(salt);
        string memory indexStr = uintToString(index);

        return stringToBytes32WithAddress(addr, valueStr, saltStr, indexStr);
    }

    function equip(address[] memory token_addresses, uint[] memory token_ids, bytes32[] memory equips, uint salt) public {
        require(token_addresses.length == token_ids.length, "Mismatch");
        
        for(uint i = 0; i < token_ids.length; i++){
            require (
                keccak256(abi.encode(concatenate(token_addresses[i], token_ids[i], salt, i)))
                ==
                keccak256(abi.encode(equips[i])),
                "Mismatch"
            );

            emit Equip(msg.sender, equips[i]);
        }
    }

    function equipWithOwnership(address[] memory token_addresses, uint[] memory token_ids, bytes32[] memory equips, uint salt) public {
        require(token_addresses.length == token_ids.length, "Mismatch");
        
        for(uint i = 0; i < token_ids.length; i++){
            
            require(IERC721(token_addresses[i]).ownerOf(token_ids[i]) == msg.sender, "Not Owner");

            require (
                keccak256(abi.encode(concatenate(token_addresses[i], token_ids[i], salt, i)))
                ==
                keccak256(abi.encode(equips[i])),
                "Mismatch"
            );

            emit Equip(msg.sender, equips[i]);
        }
    }
}