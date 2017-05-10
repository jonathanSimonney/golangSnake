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
    this.x = x;
    this.y = y;
}

function Snake(FirstName, Coordinates, Color, isYours){
    this.firstName = FirstName;
    this.coordinates = Coordinates;
    this.color = Color;
    var myObject = this;

    this.changeValue = function (FirstName, Coordinates, Color) {
        myObject.firstName = FirstName;
        myObject.coordinates = Coordinates;
        myObject.color = Color;
    };

    this.drawSnake = function(canvas){
        console.log('hello. Snake drawing!');
        if (canvas.color.indexOf(this.color) !== -1){
            console.log('drawing impossible! Please choose another color');
            return 'drawing impossible! Please choose another color';
        }

        Coordinates = myObject.coordinates;
        for (var i in Coordinates){
            canvas.fillCell(Coordinates[i], this.color);
        }
    };

    this.move = function(direction){
        console.log(direction.key, "new movement");
        console.log(JSON.stringify({Code : 0, Data : 'z'}));
        //var newCoordinates;
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
    if (isYours){
        console.log('listener added!!!');
        window.addEventListener('keypress', this.move);
    }
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
                    if (myObject.appleList[iterator].x === i && myObject.appleList[iterator].y === j){
                        myObject.canvas.fillStyle = 'green';
                        break;
                    }
                }

                myObject.canvas.fillRect(i*myObject.cellWidth,j*myObject.cellWidth,myObject.cellWidth,myObject.cellWidth);
                myObject.canvas.fillStyle = save;
            }
        }

        for (i in arraySnake){
            if (arraySnake.hasOwnProperty(i)){
                console.log(arraySnake, i, myObject);
                arraySnake[i].drawSnake(myObject);
            }
        }
    };

    this.fillCell = function(coordinate, color){
        myObject.canvas.fillStyle = color;
        myObject.canvas.fillRect(coordinate.x*this.cellWidth,coordinate.y*this.cellWidth,this.cellWidth,this.cellWidth);
    };

    /*this.getRandomCoordInCanvas = function(){
        var x = Math.floor(Math.random() * myObject.width);
        var y = Math.floor(Math.random() * myObject.height);

        return new Coordinate(x, y);
    };*/

    this.drawAnew();
}

function toggleFillStyle(canvasObject){
    canvasObject.canvas.fillStyle = canvasObject.canvas.fillStyle === canvasObject.color[0] ? canvasObject.color[1] : canvasObject.color[0];
}

function makeUpdate(messageSnakeArray, AppleArray){
    for (var i in messageSnakeArray) {
        if (messageSnakeArray.hasOwnProperty(i)) {
            if (i in arraySnake) {
                arraySnake[i].changeValue(messageSnakeArray[i].name, messageSnakeArray[i].body, messageSnakeArray[i].color);
            } else {
                var isYours = false;
                if (i === pageIndex) {
                    isYours = true;
                }

                arraySnake.push(new Snake(messageSnakeArray[i].name, messageSnakeArray[i].body, messageSnakeArray[i].color, isYours));
            }
        }
    }
    canvas.appleList = AppleArray;
    canvas.drawAnew();
}

function linkSocketListener(socket){
    /*socket.onopen = function() {
        console.log("Socket opened");
        console.log(JSON.stringify({Code : 2, Data : 'initial connection'}));
        //socket.send(JSON.stringify({Code : 2, Data : 'initial connection'}));
    };*/
    socket.onmessage = function (e) {
        var message = JSON.parse(e.data);
        console.log(message);
        switch (message.kind){
            case 'init':
                alert('init sent!');
                break;
            case 'update':
                makeUpdate(message.snakes, message.apples);
                break;
            case 'won':
                //win function
                break;
            default :
                alert('invalid kind!');
        }
        /*if (message.Code === 400){
            console.log(message.ArraySnake);
            /*console.log(snake1);
            snake1.changeValue(message.ArraySnake[0].Name, message.ArraySnake[0].Position, message.ArraySnake[0].Color);
            console.log(snake1);
            //snake1 = message.ArraySnake[0];
            for (var i in message.ArraySnake){
                if (message.ArraySnake.hasOwnProperty(i)){
                    if (i in arraySnake){
                        arraySnake[i].changeValue(message.ArraySnake[i].Name, message.ArraySnake[i].Position, message.ArraySnake[i].Color);
                    }else{
                        var isYours = false;
                        if (i === pageIndex){
                            isYours = true;
                        }

                        arraySnake.push(new Snake(message.ArraySnake[i].Name, message.ArraySnake[i].Position, message.ArraySnake[i].Color, isYours));
                    }
                }
            }
            //arraySnake = message.ArraySnake;
            canvas.appleList = message.ArrayApple;
            canvas.drawAnew();
        }else{
            alert(message.Error);
        }*/
    };
    socket.onclose = function () {
        console.log("Socket closed");
    }
}