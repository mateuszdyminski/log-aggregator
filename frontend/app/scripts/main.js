/* jshint devel:true */

(function(chat) {
    var socket = new WebSocket('ws://localhost:8080/', 'echo-protocol');
    var id = -1;

    socket.onopen = function(event) {
        console.log('Server connection open.');
    };

    socket.onmessage = function(msg) {
        var message = JSON.parse(msg.data);
        if ('identifier' === message.type) {
            id = message.data;
        }
        if ('message' === message.type) {
            appendMessage(message.source, message.data);
        }
        console.log(msg);
    }

    socket.onclose = function() {
        console.log('Server connection closed.');
        appendMessage(': (', 'Connection not available');
        socket = undefined;
    }

    socket.onerror = function() {
        console.log('Server connection failure.');
        socket = undefined;
    }

    function appendMessage(source, message) {
        var messageEl = document.createElement('li');
        var sourceEl = document.createElement('span');
        messageEl.appendChild(document.createTextNode(message));
        if (source === id) {
            messageEl.className = 'own';
        }
        document.getElementById('chat').appendChild(messageEl);
    }

    chat.send = function() {
        if (socket != null) {
            var text = document.getElementById('message').value;
            document.getElementById('message').value = '';
            socket.send(text);
        }
    };
})(this.chat = {});

document.getElementById('message').onkeydown = function(event) {
    if (event.which == 13 || event.keyCode == 13) {
        chat.send();
        return false;
    }
    return true;
};
