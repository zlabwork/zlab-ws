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
            location: '../../node_module/crypto-js-4.1.1',
            main: 'index'
        }
    ]
});
```


## 使用 js 
```js
require(["crypto-js/aes", "crypto-js/sha256"], function (AES, SHA256) {
    console.log(SHA256("Message"));
});
```
