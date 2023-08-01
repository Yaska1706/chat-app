
var app = new Vue({
    el:'#app',
    data: {
        ws: null,
        serveUrl: "ws://127.0.0.1:8080/ws",
        messages: [],
        newMessage: "",
    },
    mounted: function() {
        this.connectToWebsocket()
    },
    methods:{
        connectToWebsocket() {
            this.ws = new WebSocket(this.serveUrl);
            this.ws.addEventListener('open', (event) => {this.onWebSocketOpen(event)});
            this.ws.addEventListener('message',(event) => {this.handleNewMessage(event)})
        },
        onWebSocketOpen(){
            console.log("connected to Room")
        },
        handleNewMessage(event) {
            let data = event.data;
            data = data.Split(/\r?\n/);
            for (let i = 0; i < data; i++) {
                let msg = JSON.parse(data[i]);
                this.messages.push(msg)
            }
        },
        sendMessage() {
            if(this.newMessage !== "") {
                this.ws.send(JSON.stringify({message:this.newMessage}));
                this.newMessage = ""
            }
        }
    }

})