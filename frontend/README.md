# Websockets client project


## Preparations

### Requirements

- NodeJS & npm

### Installation of required packages

Inside of directory:
- npm install -g bower
- npm install
- bower install

### Run
- gulp serve


## Web sockets

### Web socket initialization

var socket = new WebSocket('ws://host:8080/', ['protocol-name', 'can-be-array-or-single-string']);


### Web socket methods

Events:
- socket.onopen(event);
- socket.onmessage(event);
- socket.onclose(event);
- socket.onclose(event);

Sending:
- socket.send(string);
- socket.send(Blob);
- socket.send(ArrayBuffer);

Closing:
- socket.close();

## Useful links:
- https://developer.mozilla.org/en-US/docs/WebSockets
- http://socket.io/
