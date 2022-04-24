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
    var seq = 100 // TODO::store

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
        return seq++
    }

    function getRandomBytes(len) {
        const arr = new Uint8Array(len)
        window.crypto.getRandomValues(arr)
        return arr
    }

    function getAesOption(iv) {
        return {
            iv: CryptoJS.enc.Base64.parse(iv),
            mode: CryptoJS.mode.CFB,
            padding: CryptoJS.pad.NoPadding
        };
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
    function messagePackage(type, sequence, send, session, data) {

        var i = 0

        // 1. Converts to bytes
        session = Long.fromValue(session).toBytesBE()
        sequence = Long.fromValue(sequence).toBytesBE()
        send = Long.fromValue(send).toBytesBE()

        // 2. string to bytes
        let bodyData = new TextEncoder().encode(data);
        let totalSize = headSize + bodyHeadSize + bodyData.byteLength
        let bs = new Uint8Array(totalSize)
        bs[0] = nul
        bs[1] = type

        // 3. TODO: length & AES
        // let sizeBytes = Long.fromValue(totalSize).toBytesBE()
        // for (i = 0; i < sizeBytes.length; i++) {
        //     bs[i + 2] = sizeBytes[i + 6];
        // }

        // 4. session
        for (i = 0; i < session.length; i++) {
            bs[i + headSize] = session[i];
        }

        // 5. sequence
        for (i = 0; i < sequence.length; i++) {
            bs[i + headSize + 8] = sequence[i];
        }

        // 6. sender
        for (i = 0; i < send.length; i++) {
            bs[i + headSize + 16] = send[i];
        }

        // 7. data
        for (i = 0; i < bodyData.byteLength; i++) {
            bs[i + headSize + bodyHeadSize] = bodyData[i];
        }

        // Encrypt
        let iv = getRandomBytes(16)
        let plaintext = CryptoJS.enc.Hex.parse(toHexString(bs.subarray(headSize)))
        var cipher = CryptoJS.AES.encrypt(
            plaintext,
            secretKey,
            getAesOption(window.btoa(String.fromCharCode(...iv)))
        )
        cipher = wordArrayToUint8Array(cipher.ciphertext)
        let allString = bs.subarray(0, headSize) + "," + iv + "," + cipher
        let allBytes = new Uint8Array(allString.split(","))
        let sizeBytes = Long.fromValue(allBytes.length).toBytesBE()
        for (i = 0; i < 2; i++) {
            allBytes[i + 2] = sizeBytes[i + 6];
        }

        return allBytes
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
                "send": getUserIdSender(),
                "sid": getUserIdReceiver(),
                "text": msg.value,
                "time": new Date().getTime(),
            }

            // append logs
            var item = document.createElement("div");
            item.innerHTML = '<div class="textSend">' + msg.value + '</div>';
            appendLog(item);

            let bs = messagePackage(msgTypeText, sequence, data.send, data.sid, msg.value)
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
                    var plaintext = CryptoJS.AES.decrypt(ciphertext, secretKey, getAesOption(iv));
                    body = wordArrayToUint8Array(plaintext)

                    // slice
                    var mid = body.subarray(8, bodyHeadSize)
                    var send = Long.fromBytes(body.subarray(16, bodyHeadSize)).toString()
                    var recv = Long.fromBytes(body.subarray(0, 8)).toString()
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
