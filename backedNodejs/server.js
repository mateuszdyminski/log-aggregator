#!/usr/bin/env node

var WebSocketServer = require('websocket').server;
var http = require('http');
var _ = require('lodash');
var Firebase = require("firebase");

var kafka = require('kafka-node'),
    Consumer = kafka.Consumer,
    client = new kafka.Client("zk:2181"),
    consumer = new Consumer(
        client,
        [
            { topic: 'topic', partition: 0 }
        ],
        {
            autoCommit: false
        });


var server = http.createServer(function(request, response) {
    console.log((new Date()) + ' Received request for ' + request.url);
    response.writeHead(404);
    response.end();
});
server.listen(8080, function() {
    console.log((new Date()) + ' Server is listening on port 8080');
});

wsServer = new WebSocketServer({
    httpServer: server,
    // You should not use autoAcceptConnections for production
    // applications, as it defeats all standard cross-origin protection
    // facilities built into the protocol and the browser.  You should
    // *always* verify the connection's origin and decide whether or not
    // to accept it.
    autoAcceptConnections: false
});

function originIsAllowed(origin) {
    // put logic here to detect whether the specified origin is allowed.
    return true;
}

function sendMsg(source, message) {
    console.log('Sending message to ' + connections.length + ' receivers');
    connections.forEach(function(output) {
        var msg = JSON.stringify({
            type: 'message',
            source: source,
            data: message
        });
        output.sendUTF(msg);
    });
}

var connections = [];
var ids = 0;

wsServer.on('request', function(request) {
    if (!originIsAllowed(request.origin)) {
        // Make sure we only accept requests from an allowed origin
        request.reject();
        console.log((new Date()) + ' Connection from origin ' + request.origin + ' rejected.');
        return;
    }

    // Accept connection
    var connection = request.accept('echo-protocol', request.origin);
    connection.connectionId = ++ids;
    console.log((new Date()) + ' Connection ' + connection.connectionId + ' accepted.');

    connections.push(connection);

    // Send identifier to client
    connection.sendUTF(JSON.stringify({
        type: 'identifier',
        data: connection.connectionId
    }));

    // Send hello message
    sendMsg(connection.connectionId, 'Has joined chat');

    connection.on('message', function(message) {
        if (message.type === 'utf8') {
            console.log('Received Message: ' + message.utf8Data);
            sendMsg(connection.connectionId, message.utf8Data);
        } else if (message.type === 'binary') {
            console.log('Received Binary Message of ' + message.binaryData.length + ' bytes');
            connection.sendBytes(message.binaryData);
        }
    });

    connection.on('close', function(reasonCode, description) {
        console.log((new Date()) + ' Peer ' + connection.remoteAddress + ' disconnected.');
        sendMsg(connection.connectionId, 'Has left chat');
        _.remove(connections, connection);
    });
});

// kafka
consumer.on('message', function(message) {
    console.log('Receiving log from kafka: ' + message);
    for(var conn in connections) {
        sendMsg(conn.connectionId, message);
    }
});
