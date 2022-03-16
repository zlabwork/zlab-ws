define(function (require) {

    // @link https://github.com/dcodeIO/Long.js
    // @link https://www.cnblogs.com/luludongxu/p/13366521.html
    // @link https://blog.csdn.net/humanbeng/article/details/122010117

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

    function getUserIdSender() {
        return Number(document.getElementById("senderId").value);
    }

    function getUserIdReceiver() {
        return Number(document.getElementById("receiverId").value);
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
    // text to bytes with message type
    function messagePackage(msgType, text) {
        // string 转 bytes
        let textBytes = new TextEncoder().encode(text);
        let bs = new Uint8Array(textBytes.byteLength + headSize + bodyHeadSize)
        // head
        bs[0] = nul
        bs[1] = msgType
        for (var i = 0; i < textBytes.byteLength; i++) {
            bs[i + headSize + bodyHeadSize] = textBytes[i];
        }
        return bs
    }

    function messagePackageNew(msgType, sequence, sender, receiver, text) {
        // string 转 bytes
        let textBytes = new TextEncoder().encode(text);
        let bs = new Uint8Array(textBytes.byteLength + headSize + bodyHeadSize)
        // TODO: int64 to bytes
        // head
        bs[0] = nul
        bs[1] = msgType
        for (var i = 0; i < textBytes.byteLength; i++) {
            bs[i + headSize + bodyHeadSize] = textBytes[i];
        }
        return bs
    }

    // auth
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
        // TODO:
        let bs = messagePackage(msgTypeAuth, JSON.stringify(data))
        conn.send(bs);
    }

    document.getElementById("messageForm").onsubmit = function () {
        if (!conn) {
            return false;
        }
        if (!msg.value) {
            return false;
        }

        let data = {
            "id": getUUID(),
            "sender": getUserIdSender(),
            "receiver": getUserIdReceiver(),
            "text": msg.value,
            "time": new Date().getTime(),
        }
        // let bs = messagePackage(msgTypeText, JSON.stringify(data))
        let bs = messagePackageNew(msgTypeText, 123456, 111111, 222222, msg.value)
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
