define(function (require) {

    const HOST = "ws://127.0.0.1:3000/ws"
    const Long = require("long")
    const nul = 0x00
    const lf = 0x0A // 换行符
    const us = 0x1F // 单元分隔符
    const msgTypeAuth = 0x03
    const msgTypeText = 0x20

    const headSize = 4 // 32 Bit
    const bodyHeadSize = 24 // 192 Bit

    const deviceUUID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    const testToken = "dev-token"

    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");

    function getQuery(args) {
        var query = window.location.search.substring(1);
        var vars = query.split("&");
        for (var i = 0; i < vars.length; i++) {
            var pair = vars[i].split("=");
            if (pair[0] == args) {
                return pair[1];
            }
        }
        return (false);
    }

    function getSequenceId() {
        let n = Math.floor(Math.random() * 65535);
        let ts = new Date().getTime()
        return Long.fromValue(ts).shiftLeft(16).or(n)
    }

    function getUserIdSender() {
        return document.getElementById("senderId").value.trim();
    }

    function getUserIdReceiver() {
        return document.getElementById("receiverId").value.trim();
    }

    function getUuid() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }
    }

    // head: | 8 Bit Unused | 8 Bit Type | 16 Bit Length ｜
    // body: | 64 Bit SequenceID | 64 Bit SenderID | 64 Bit ReceiverID | Data ｜
    function messagePackage(type, sequence, sender, receiver, data) {

        var i = 0

        // 1. Converts to bytes
        sequence = Long.fromValue(sequence).toBytesBE()
        sender = Long.fromValue(sender).toBytesBE()
        receiver = Long.fromValue(receiver).toBytesBE()

        // 2. string to bytes
        let bodyData = new TextEncoder().encode(data);
        let totalSize = headSize + bodyHeadSize + bodyData.byteLength
        let bs = new Uint8Array(totalSize)
        bs[0] = nul
        bs[1] = type

        // 3. TODO: length & AES
        let sizeBytes = Long.fromValue(totalSize).toBytesBE()
        for (i = 0; i < sizeBytes.length; i++) {
            bs[i + 2] = sizeBytes[i + 6];
        }

        // 4. sequence
        for (i = 0; i < sequence.length; i++) {
            bs[i + headSize] = sequence[i];
        }

        // 5. sender
        for (i = 0; i < sender.length; i++) {
            bs[i + headSize + 8] = sender[i];
        }

        // 6. receiver
        for (i = 0; i < receiver.length; i++) {
            bs[i + headSize + 16] = receiver[i];
        }

        // 7. data
        for (i = 0; i < bodyData.byteLength; i++) {
            bs[i + headSize + bodyHeadSize] = bodyData[i];
        }

        return bs
    }

    function init() {

        // argv
        let sender = getQuery("sender")
        let receiver = getQuery("receiver")
        if (receiver != false) {
            document.getElementById("receiverId").value = receiver;
        }
        if (sender != false) {
            document.getElementById("senderId").value = sender;
            document.title = "User: " + sender
        }

        // send auth
        document.getElementById("auth").onclick = function () {
            if (!conn) {
                return false;
            }
            let data = {
                "id": getUuid(),
                "sender": getUserIdSender(),
                "token": testToken,
                "uuid": deviceUUID,
                "version": "v1.0",
                "os": "ios",
                "dateTime": new Date().getTime()
            }
            let bs = messagePackage(msgTypeAuth, 0, 0, 0, JSON.stringify(data))
            conn.send(bs);
        }

        // send message
        document.getElementById("submit").onclick = function () {
            if (!conn) {
                return false;
            }
            if (!msg.value) {
                return false;
            }

            let sequence = getSequenceId()
            let data = {
                "id": getUuid(),
                "sender": getUserIdSender(),
                "receiver": getUserIdReceiver(),
                "text": msg.value,
                "time": new Date().getTime(),
            }

            // append logs
            var item = document.createElement("div");
            item.innerHTML = '<div class="textSend">' + msg.value + '</div>';
            appendLog(item);

            let bs = messagePackage(msgTypeText, sequence, data.sender, data.receiver, msg.value)
            conn.send(bs);
            msg.value = "";
            return false;
        };

        if (window["WebSocket"]) {
            // conn = new WebSocket("ws://" + document.location.host + "/ws");
            conn = new WebSocket(HOST);
            conn.binaryType = "arraybuffer"
            conn.onclose = function (evt) {
                var item = document.createElement("div");
                item.innerHTML = "<b>Connection closed.</b>";
                appendLog(item);
            };
            conn.onmessage = function (event) {

                if (event.data instanceof ArrayBuffer === false) {
                    return
                }

                var head = new Uint8Array(event.data, 0, headSize);
                var body = new Uint8Array(event.data, headSize)

                // slice
                var mid = body.subarray(0, 16)
                var send = body.subarray(8, 16)
                var recv = body.subarray(16, bodyHeadSize)
                var data = body.subarray(bodyHeadSize, body.byteLength)
                var receiver = Long.fromBytes(recv).toString()
                var sender = Long.fromBytes(send).toString()
                var text = new TextDecoder().decode(data)

                // logs
                console.log("Message Id: " + mid)

                // insert
                var item = document.createElement("div");
                item.innerHTML = '<div class="textReceive">' + text + '</div>';
                appendLog(item);
            };
        } else {
            var item = document.createElement("div");
            item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
            appendLog(item);
        }
    }

    init()
});
