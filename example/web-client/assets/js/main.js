require.config({
    packages: [
        {
            name: 'crypto-js',
            location: '../../node_module/crypto-js-4.1.1',
            main: 'index'
        }
    ]
});

require(["crypto-js/aes", "crypto-js/sha256"], function (AES, SHA256) {
    console.log(SHA256("Message"));
});