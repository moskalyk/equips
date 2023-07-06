import fetch from 'cross-fetch';
import {ethers} from 'ethers'

let startTime: any, endTime: any;

const start = () => {
  startTime = new Date();
}

const end = () => {
  endTime = new Date();
  return Math.round(endTime - startTime)
};

function calculateAverage(numbers: any) {
    const sum = numbers.reduce((total: any, num: any) => total + num, 0);
    const average = sum / numbers.length;
    return average;
}

(async () => {

  function bytes32ToStringWithAddress(bytes32) {
    const combinedBytes = ethers.utils.arrayify(bytes32);
    const addrBytes = combinedBytes.slice(0, 20);
    const addressConverted = ethers.utils.getAddress(ethers.utils.hexlify(addrBytes));
    console.log(addressConverted)
    // const string = ethers.utils.parseBytes32String(bytes32);

    let string = "";
    for (let i = 20; i < combinedBytes.length; i++) {
      const byte = combinedBytes[i];
      if (byte === 0) {
        break;
      }
      string += String.fromCharCode(byte);
    }

  
    return { addressConverted, string };
  }

  
  function rightPadBytes(byteArray, length) {
    if (byteArray.length >= length) {
      return byteArray.slice(0, length);
    }
    
    const paddedArray = new Uint8Array(length);
    paddedArray.set(byteArray, 0);
    
    return paddedArray;
  }

  function stringToBytes32WithAddress(address, str, salt, index) {
    const addrBytes = ethers.utils.arrayify(address);
    const strBytes = ethers.utils.toUtf8Bytes(str);
    const saltBytes = ethers.utils.toUtf8Bytes(salt);
    const indexBytes = ethers.utils.toUtf8Bytes(index);

    const div = ethers.utils.toUtf8Bytes(':');
    const combinedBytes = ethers.utils.concat([addrBytes, div, strBytes, div, saltBytes, div, indexBytes]);
    // const paddedBytes = ethers.utils.hexZeroPad(combinedBytes, 32);
    // combinedBytes.padEnd(66, "0");
    return ethers.utils.hexlify(rightPadBytes(combinedBytes, 32));
  }
  
  const address = "0x99AB4d7B127311072e5D159BB30BDf20669aA1a4";
  const str = "6";
  
  const bytes32 = stringToBytes32WithAddress(address, str, "100", "0");
  
  console.log("Address:", address);
  console.log("String:", str);
  console.log('before keccak', bytes32)
  const bytes32Value = ethers.utils.hexZeroPad('0x99ab4d7b127311072e5d159bb30bdf20669aa1a43a363a3100', 32);

  // Encode the string using ABI encoding
  const encodedString = ethers.utils.defaultAbiCoder.encode(['string'], ['0x99ab4d7b127311072e5d159bb30bdf20669aa1a43a363a3100']);
  console.log()
  // Calculate the keccak256 hash of the encoded string
  const hash = ethers.utils.keccak256(encodedString);
  console.log('HASH')
  console.log(hash);

  console.log("bytes32:", bytes32Value);
  console.log("bytes32:", '0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef');

  // Define the types and values to encode
  const types = ['address', 'uint', 'uint'];
  const values = [address, 6, 100];

  // Encode the values using abi.encode
  const encodedData = ethers.utils.defaultAbiCoder.encode(types, values);
  console.log(encodedData);

  const decodedValues = ethers.utils.defaultAbiCoder.decode(types, encodedData);
  console.log(decodedValues);

  const { addressConverted, string } = bytes32ToStringWithAddress(bytes32);
  // console.log(addressConverted)
  // console.log(string)
  // function stringToBytes32(str) {
  //   // Pad the string to 32 bytes (64 hex characters)
  //   const paddedStr = ethers.utils.formatBytes32String(str);
    
  //   // Convert the padded string to a Bytes32 value
  //   const bytes32 = ethers.utils.hexZeroPad(paddedStr, 32);

  //   return bytes32;
  // }

  // const str = "0x99AB4d7B127311072e5D159BB30BDf20669aA1a4:0:100"; // Example string to convert to bytes32

  // const bytes32 = stringToBytes32(str);

  // console.log("String:", str);
  // console.log("bytes32:", bytes32);

    const times: any = []
    for(let i = 0; i < 100; i++){
        start()
        // const res = await fetch('http://localhost:7077/all')
        times.push(end())
    }
    console.log(calculateAverage(times))
})()    