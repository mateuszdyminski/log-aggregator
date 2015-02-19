/* jshint devel:true */

(function(chat) {

    var LogLevel = {
        "I": "info",
        "W": "warning",
        "E": "error",
        "F": "fatal"
    }

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
        var message = JSON.parse(msg.data);
        console.debug(message.Host);
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
        hostEl.appendChild(document.createTextNode(message.Host));
        messageEl.appendChild(hostEl);
        messageEl.appendChild(document.createTextNode(message.Content));
        setClass(message, messageEl);
        document.getElementById('chat').appendChild(messageEl);
    }

    function setClass(message, element) {
        var levelType = LogLevel[message.Level];
        if (levelType !== undefined) {
            var classAttr = document.createAttribute("class");
            classAttr.value = LogLevel[message.Level];
            element.setAttributeNode(classAttr);
        }
    }
})(this.chat = {});

document.getElementById('message').onkeydown = function(event) {
    if (event.which == 13 || event.keyCode == 13) {
        chat.send();
        return false;
    }
    return true;
};
