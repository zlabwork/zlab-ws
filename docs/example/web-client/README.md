## 参考文档
https://github.com/requirejs/requirejs  
https://github.com/volojs/create-template  
https://github.com/requirejs/example-multipage  

## 配置 require.config
```js
require.config({
    baseUrl: 'assets/js',
    paths: {
        jquery: 'lib/jquery-3.6.0.min'
    },
    packages: [
        {
            name: 'crypto-js',
            location: '../../node_modules/crypto-js',
            main: 'index'
        },
        {
            name: 'long',
            location: '../../node_modules/long/umd',
            main: 'index'
        }
    ]
});

```


## 使用 js 
```js
// https://www.npmjs.com/package/crypto-js
require(["crypto-js/aes", "crypto-js/sha256"], function (AES, SHA256) {
    console.log(SHA256("Message"));
});

// https://www.npmjs.com/package/long
require(["long"], function (Long) {
    var value1 = Long.fromString("1234567890123456789")
    var value2 = new Long(0xFFFFFFFF, 0x7FFFFFFF);
    console.log(value1.toBytes());
    console.log(value2.toString());
});
```

```js
// Uint8Array to Uint16

// var bf = new Uint8Array([1, 2]) // 仅适用方法一
var bf = new Buffer([1, 2]) // 适用两种方法

// 方法一
var view = new DataView(bf.buffer, 0);
var i = view.getUint16(0, false);
console.log(i)

// 方法二
var i = bf.readUInt16BE(0)
console.log(i);

```
