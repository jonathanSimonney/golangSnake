Objective of this repository is to make a kinda slither online which will be able to 
let multiple player play together online.

To use it, the server.go shall be launched as it is the server and the different
 player should connect using the .html and assets files. 
 
With go installed, you will need to
```
go get golang.org/x/net/websocket
go install golang.org/x/net/websocket
```

Since this is my demo app for heroku, if you too want to deploy it to heroku, don't forget to scale your app with one dyno (with :
```heroku ps:scale web=1```