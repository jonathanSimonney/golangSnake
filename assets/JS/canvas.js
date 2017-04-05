Array.prototype.containsArray = function(val) {
    var hash = {};
    for(var i=0; i<this.length; i++) {
        hash[this[i]] = i;
    }
    return hash.hasOwnProperty(val);
};

// Warn if overriding existing method
if(Array.prototype.equals)
    console.warn("Overriding existing Array.prototype.equals. Possible causes: New API defines the method, there's a framework conflict or you've got double inclusions in your code.");
// attach the .equals method to Array's prototype to call it on any array
Array.prototype.equals = function (array) {
    // if the other array is a falsy value, return
    if (!array)
        return false;

    // compare lengths - can save a lot of time
    if (this.length != array.length)
        return false;

    for (var i = 0, l=this.length; i < l; i++) {
        // Check if we have nested arrays
        if (this[i] instanceof Array && array[i] instanceof Array) {
            // recurse into the nested arrays
            if (!this[i].equals(array[i]))
                return false;
        }
        else if (this[i] != array[i]) {
            // Warning - two different object instances will never be equal: {x:20} != {x:20}
            return false;
        }
    }
    return true;
};
// Hide method from for-in loops
//Object.defineProperty(Array.prototype, "equals", {enumerable: false});

function Coordinate(x, y){
    this.X = x;
    this.Y = y;
}

function Snake(FirstName, Coordinates, Color){
    this.FirstName = FirstName;
    this.Coordinates = Coordinates;
    this.Color = Color;
    var myObject = this;

    this.changeValue = function (FirstName, Coordinates, Color) {
        myObject.FirstName = FirstName;
        myObject.Coordinates = Coordinates;
        myObject.Color = Color;
    };

    this.drawSnake = function(canvas){
        if (canvas.color.indexOf(this.Color) !== -1){
            return 'drawing impossible! Please choose another color';
        }

        Coordinates = myObject.Coordinates;
        for (var i in Coordinates){
            canvas.fillCell(Coordinates[i], this.Color);
        }
    };

    this.move = function(direction){
        console.log(direction.key, "new movement");
        //var newCoordinates;
        switch (direction.key){
            case 'z' :
                //console.log(myObject.Coordinates, myObject);
                websocket.send(JSON.stringify({Code : 0, Data : 'z'}));
                //newCoordinates = new Coordinate(myObject.Coordinates[0][0], myObject.Coordinates[0][1]-1);//todo comment
                break;
            case 'q':
                websocket.send(JSON.stringify({Code : 0, Data : 'q'}));
                //newCoordinates = new Coordinate(myObject.Coordinates[0][0]-1, myObject.Coordinates[0][1]);//todo comment
                break;
            case 's':
                websocket.send(JSON.stringify({Code : 0, Data : 's'}));
                //newCoordinates = new Coordinate(myObject.Coordinates[0][0], myObject.Coordinates[0][1]+1);//todo comment
                break;
            case 'd':
                websocket.send(JSON.stringify({Code : 0, Data : 'd'}));
                //newCoordinates = new Coordinate(myObject.Coordinates[0][0]+1, myObject.Coordinates[0][1]);//todo comment
                break;
            default:
                return 'invalid direction supplied. No move will be made.';
        }

        /*myObject.Coordinates.unshift(newCoordinates);//todo uncomment and use coords as object
        if (canvas.appleList.containsArray(newCoordinates)){
            //canvas.createApple(myObject.Coordinates);
            console.log(canvas.appleList);
            canvas.appleList = canvas.appleList.filter(function(item) {
                return !item.equals(newCoordinates);
            });
            console.log(canvas.appleList);
        }else{
            myObject.Coordinates.pop();
        }
        canvas.drawAnew();*/
    };

    window.addEventListener('keypress', this.move);
}

function Canvas(htmlElement, cellWidth, color){
    this.htmlElement = htmlElement;
    this.canvas = htmlElement.getContext('2d');
    this.cellWidth = cellWidth;
    this.color = color;
    this.width = Math.floor(this.htmlElement.width/cellWidth);
    this.height = Math.floor(this.htmlElement.height/cellWidth);
    this.appleList= [];
    var myObject = this;

    this.drawAnew = function(){
        myObject.canvas.fillStyle = color[0];
        var olderFillStyle = color[0];

        for (var i = 0; i < myObject.width; i++){
            if (myObject.canvas.fillStyle !== olderFillStyle){
                toggleFillStyle(myObject);
            }
            for (var j = 0; j < myObject.height; j++){
                toggleFillStyle(myObject);
                if (j === 0){
                    olderFillStyle = myObject.canvas.fillStyle;
                }

                var save = myObject.canvas.fillStyle;

                for (var iterator in myObject.appleList){
                    if (myObject.appleList[iterator].X === i && myObject.appleList[iterator].Y === j){
                        myObject.canvas.fillStyle = 'green';
                        break;
                    }
                }

                myObject.canvas.fillRect(i*myObject.cellWidth,j*myObject.cellWidth,myObject.cellWidth,myObject.cellWidth);
                myObject.canvas.fillStyle = save;
            }
        }
        snake1.drawSnake(myObject);
    };

    this.fillCell = function(coordinate, color){
        myObject.canvas.fillStyle = color;
        myObject.canvas.fillRect(coordinate.X*this.cellWidth,coordinate.Y*this.cellWidth,this.cellWidth,this.cellWidth);
    };

    this.getRandomCoordInCanvas = function(){
        var x = Math.floor(Math.random() * myObject.width);
        var y = Math.floor(Math.random() * myObject.height);

        return new Coordinate(x, y);
    };

    this.drawAnew();
}

function toggleFillStyle(canvasObject){
    canvasObject.canvas.fillStyle = canvasObject.canvas.fillStyle === canvasObject.color[0] ? canvasObject.color[1] : canvasObject.color[0];
}

function linkSocketListener(socket){
    socket.onopen = function() {
        console.log("Socket opened");
        socket.send(JSON.stringify({Code : 1, Data : 'initial connection'}));
    };
    socket.onmessage = function (e) {
        var message = JSON.parse(e.data);
        if (message.Code === 400){
            console.log(snake1);
            snake1.changeValue(message.ArraySnake[0].Name, message.ArraySnake[0].Position, message.ArraySnake[0].Color);
            console.log(snake1);
            //snake1 = message.ArraySnake[0];
            canvas.appleList = message.ArrayApple;
            canvas.drawAnew();
        }
    };
    socket.onclose = function () {
        console.log("Socket closed");
    }
}

var canvas;
var websocket;
var snake1 = new Snake();

window.onload = function () {
    var htmlElementCanvas = document.getElementsByTagName('canvas')[0];
    if (htmlElementCanvas.getContext) {
        websocket = new WebSocket('ws://localhost:8081/single');
        linkSocketListener(websocket);

        canvas = new Canvas(htmlElementCanvas, 50, ['#ff0000', '#0000ff']);
        snake1.drawSnake(canvas);

        canvas.drawAnew();
    }
};