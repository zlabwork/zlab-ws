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
        },
        {
            name: 'uuid',
            location: '../../node_modules/uuid/dist/umd',
            main: 'uuid.min'
        }
    ]
});

require(['app/websocket']);
