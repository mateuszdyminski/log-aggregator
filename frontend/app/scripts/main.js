/* jshint devel:true */

(function(chat) {
    var loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    } else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.hostname + ":8090/wsapi/ws";

    var socket = new WebSocket(new_uri);
    var id = -1;

    socket.onopen = function(event) {
        console.log('Server connection open.');
    };

    socket.onmessage = function(msg) {
        console.log(msg);
        var message = JSON.parse(msg.data);
        appendMessage(message);
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

    function appendMessage(message) {
        var messageEl = document.createElement('li');
        var hostEl = document.createElement('span');
        hostEl.appendChild(document.createTextNode(message.Key));
        messageEl.appendChild(hostEl);

        messageEl.appendChild(document.createTextNode(message.Message));

        document.getElementById('chat').appendChild(messageEl);
    }
})(this.chat = {});

document.getElementById('message').onkeydown = function(event) {
    if (event.which == 13 || event.keyCode == 13) {
        chat.send();
        return false;
    }
    return true;
};
