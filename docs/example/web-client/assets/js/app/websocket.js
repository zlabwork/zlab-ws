define(function (require) {

    const HOST = "ws://127.0.0.1:3000/ws"
    const Long = require("long")
    const CryptoJS = require("crypto-js");
    const nul = 0x00
    const lf = 0x0A // 换行符
    const us = 0x1F // 单元分隔符
    const msgTypeAuth = 0x03
    const msgTypeText = 0x20
    const epoch = 1288834974657

    const headSize = 4 // 32 Bit
    const bodyHeadSize = 24 // 192 Bit

    const deviceUUID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    const testToken = "dev-token"

    var conn;
    var msg = document.getElementById("msg");
    var log = document.getElementById("log");
    var secretKey = CryptoJS.enc.Hex.parse("ffffffffffffffffffffffffffffffff")

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

    // 7 bit + 41 bit + 16 bit
    function getSequenceId(isPrivate) {
        var s = 0 // TODO: 7 bit
        if (isPrivate === true) {
            s = 0b0100000
        } else {
            s = 0b1000000
        }
        let t = new Date().getTime() - epoch
        let n = Math.floor(Math.random() * 65535);
        return Long.fromValue(s).shiftLeft(41).or(t).shiftLeft(16).or(n)
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

    function toUint16(bf, offset, length) {
        let view = new DataView(bf.buffer, offset, length);
        return view.getUint16(0, false);
    }

    function toHexString(bytes) {
        return Array.from(bytes, function (byte) {
            return ('0' + (byte & 0xFF).toString(16)).slice(-2);
        }).join('')
    }

    function wordArrayToUint8Array(wordArray) {
        var len = wordArray.words.length,
            u8_array = new Uint8Array(len << 2),
            offset = 0, word, i
        ;
        for (i = 0; i < len; i++) {
            word = wordArray.words[i];
            u8_array[offset++] = word >> 24;
            u8_array[offset++] = (word >> 16) & 0xff;
            u8_array[offset++] = (word >> 8) & 0xff;
            u8_array[offset++] = word & 0xff;
        }
        return u8_array.subarray(0, wordArray.sigBytes);
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

            let sequence = getSequenceId(true)
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

                var length = event.data.byteLength
                var arrayData = new Uint8Array(event.data, 0)
                var offset = 0

                while (offset < length) {
                    var l = toUint16(arrayData, offset + 2, 2)
                    if (l === 0) {
                        return;
                    }
                    var data = new Uint8Array(event.data, offset, l)
                    var head = data.subarray(0, headSize)
                    if (head[0] !== nul) {
                        return;
                    }
                    var body = data.subarray(headSize, l)

                    // Decrypt
                    var iv = window.btoa(String.fromCharCode(...body.subarray(0, 16)))
                    var ciphertext = window.btoa(String.fromCharCode(...body.subarray(16)))
                    var options = {
                        iv: CryptoJS.enc.Base64.parse(iv),
                        mode: CryptoJS.mode.CFB,
                        padding: CryptoJS.pad.NoPadding
                    };
                    var plaintext = CryptoJS.AES.decrypt(ciphertext, secretKey, options);
                    body = wordArrayToUint8Array(plaintext)

                    // slice
                    var mid = body.subarray(0, 16)
                    var send = Long.fromBytes(body.subarray(8, 16)).toString()
                    var recv = Long.fromBytes(body.subarray(16, bodyHeadSize)).toString()
                    var content = body.subarray(bodyHeadSize, body.byteLength)
                    var text = new TextDecoder().decode(content)

                    // logs
                    // console.log("send:" + send + ", recv:" + recv)
                    console.log("Message Id: " + toHexString(mid))

                    // insert
                    var item = document.createElement("div");
                    item.innerHTML = '<div class="textReceive">' + text + '</div>';
                    appendLog(item);

                    // offset
                    offset += l
                }
            };
        } else {
            var item = document.createElement("div");
            item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
            appendLog(item);
        }
    }

    init()
});
