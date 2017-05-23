function linkJoinEvent(){
    for (var i in document.querySelectorAll("#join>button")) {
        document.querySelectorAll("#join>button")[i].onclick = function (e) {
            websocket.send(JSON.stringify({kind : 'connect', slot : parseInt(this.name), name: 'change', color: 'red'}));
            disable(document.querySelectorAll("#join>button"));
            window.addEventListener('keypress', move);
            startButton.disabled = parseInt(e.target.name) === 0;
        }
    }

    document.getElementById('start').onclick = function () {
        websocket.send(JSON.stringify({kind : 'start'}));
    }
}

function move(direction){
    switch (direction.key){
        case 'z' :
            //console.log(myObject.Coordinates, myObject);
            websocket.send(JSON.stringify({kind : 'move', key : 'up'}));
            //newCoordinates = new Coordinate(myObject.Coordinates[0][0], myObject.Coordinates[0][1]-1);//todo comment
            break;
        case 'q':
            websocket.send(JSON.stringify({kind : 'move', key : 'left'}));
            //newCoordinates = new Coordinate(myObject.Coordinates[0][0]-1, myObject.Coordinates[0][1]);//todo comment
            break;
        case 's':
            websocket.send(JSON.stringify({kind : 'move', key : 'down'}));
            //newCoordinates = new Coordinate(myObject.Coordinates[0][0], myObject.Coordinates[0][1]+1);//todo comment
            break;
        case 'd':
            websocket.send(JSON.stringify({kind : 'move', key : 'right'}));
            //newCoordinates = new Coordinate(myObject.Coordinates[0][0]+1, myObject.Coordinates[0][1]);//todo comment
            break;
        default:
            return 'invalid direction supplied. No move will be made.';
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