var canvas;
var websocket;
var arraySnake = [];
var pageIndex;

window.onload = function () {
    var htmlElementCanvas = document.getElementsByTagName('canvas')[0];
    if (htmlElementCanvas.getContext) {
        websocket = new WebSocket('ws://localhost:8081/single');
        linkSocketListener(websocket);

        canvas = new Canvas(htmlElementCanvas, 50, ['#ff0000', '#0000ff']);

        canvas.drawAnew();
    }

    linkJoinEvent();
};