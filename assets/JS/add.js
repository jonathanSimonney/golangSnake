function linkSocketListener(socket, initDataSent){
    socket.onopen = function() {
        socket.send(initDataSent);
    };
    socket.onmessage = function (e) {
        console.log(e);
    };
    socket.onclose = function () {
        console.log("Socket closed");
    }
}

window.onload = function(){
    document.getElementById('myForm').onsubmit = function(){
        var objectData = {
            'Name' : this['name'].value,
            'X'    : parseInt(this['xPosition'].value),
            'Y'    : parseInt(this['yPosition'].value)
        };

        console.log(JSON.stringify(objectData));

        var websocket = new WebSocket('ws://localhost:8081/');

        linkSocketListener(websocket, JSON.stringify(objectData));

        return false;
    }
}