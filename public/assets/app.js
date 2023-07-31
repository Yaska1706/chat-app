
var app = new Vue({
    el:'#app',
    data: {
        ws: null,
        serveUrl: "ws://127.0.0.1:8080/ws"
    },
    mounted: function() {
        this.connectToWebsocket()
    },
    methods:{
        connectToWebsocket() {
            this.ws = new WebSocket(this.serveUrl);
            this.ws.addEventListener('open', (event) => {this.onWebSocketOpen(event)});
        },
        onWebSocketOpen(){
            console.log("connected to Room")
        }
    }

})