define(function (require) {

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

    function getUUID() {
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

    // send auth
    document.getElementById("auth").onclick = function () {
        if (!conn) {
            return false;
        }
        let data = {
            "id": getUUID(),
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
    document.getElementById("messageForm").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        let sequence = getSequenceId()
        let data = {
            "id": getUUID(),
            "sender": getUserIdSender(),
            "receiver": getUserIdReceiver(),
            "text": msg.value,
            "time": new Date().getTime(),
        }

        // append logs
        var item = document.createElement("div");
        item.innerHTML = msg.value;
        appendLog(item);

        let bs = messagePackage(msgTypeText, sequence, data.sender, data.receiver, msg.value)
        conn.send(bs);
        msg.value = "";
        return false;
    };

    if (window["WebSocket"]) {
        // conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn = new WebSocket("ws://127.0.0.1:3000/ws");
        conn.binaryType = "arraybuffer"
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (event) {
            var text = ""
            if (typeof event.data === "string") {
                text = event.data
            }
            if (event.data instanceof ArrayBuffer) {
                var bs = new Uint8Array(event.data)
                var data = bs.subarray(2, bs.byteLength)
                text = new TextDecoder().decode(data);
            }
            var messages = text.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
});
