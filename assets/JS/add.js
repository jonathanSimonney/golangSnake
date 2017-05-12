function linkJoinEvent(){
    for (var i in document.querySelectorAll("#join>button")) {
        document.querySelectorAll("#join>button")[i].onclick = function () {
            websocket.send(JSON.stringify({kind : 'connect', slot : parseInt(this.name), name: 'change', color: 'black'}));
            disable(document.querySelectorAll("#join>button"));
            pageIndex = this.name;
        }
    }

    document.getElementById('start').onclick = function () {
        websocket.send(JSON.stringify({kind : 'start'}));
    }
}






/*function linkSocketListener(socket){
    socket.onopen = function() {
        console.log("Socket opened");
        socket.send(JSON.stringify({Code : 0, Input : 'initial connection'}));
    };
    socket.onmessage = function (e) {
        console.log(e);
    };
    socket.onclose = function () {
        console.log("Socket closed");
    }
}

window.onload = function(){
    var websocket = new WebSocket('ws://localhost:8081/');

    linkSocketListener(websocket);
    document.getElementById('myForm').onsubmit = function(){
        var objectData = {
            'Name' : this['name'].value,
            'X'    : parseInt(this['xPosition'].value),
            'Y'    : parseInt(this['yPosition'].value)
        };

        console.log(JSON.stringify({Code : 1, Data : objectData}));
        websocket.send(JSON.stringify({Code : 1, Data : objectData}));

        return false;
    }
}*/