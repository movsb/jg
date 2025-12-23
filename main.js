/** @type {ArrayBuffer} */
const blob = http.get('http://baidu.com').blob();
console.log(blob.byteLength);
