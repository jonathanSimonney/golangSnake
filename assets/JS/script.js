var canvas;//todo change so everything is dynamic!!! , changer code par kind.
var websocket;
var arraySnake = [];
var pageIndex;
var cellWidth = 10;

window.onload = function () {
    var htmlElementCanvas = document.getElementsByTagName('canvas')[0];
    if (htmlElementCanvas.getContext) {
        websocket = new WebSocket('ws://localhost:8081/');
        linkSocketListener(websocket);

        canvas = new Canvas(htmlElementCanvas, cellWidth, ['#ff0000', '#0000ff']);

        canvas.drawAnew();
    }

    linkJoinEvent();
};