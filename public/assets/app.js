
var app = new Vue({
    el:'#app',
    data: {
        ws: null,
        serveUrl: "ws://127.0.0.1:8080/ws",
        roomInput:null,
        rooms : [],
        user: {
            name:""
        },
        users: []
    },
    mounted: function() {
        this.connectToWebsocket()
    },
    methods:{
        connectToWebsocket() {
            this.ws = new WebSocket(this.serveUrl + "?name=" + this.user.name);
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
                const room = this.findRoom(msg.target);
                if (typeof room !== "undefined") {
                    room.message.push(msg);
                }
            }
        },
        sendMessage(room) {
            if(room.newMessage !== "") {
                this.ws.send(JSON.stringify({
                    action: 'send-message',
                    target: room.name,
                    message:this.newMessage}));
                room.newMessage = ""
            }
        },
        findRoom(roomName) {
            for (let i=0; i < this.rooms.length ; i ++){
                if (this.rooms[i].name === roomName) {
                    return this.rooms[i]
                }
            }
        },
        joinRoom() {
            this.ws.send(JSON.stringify({
                action: 'join-room',
                message: this.roomInput,
            }));
            this.messages = [],
            this.rooms.push({
                "name":this.roomInput,
                "messages":[],
            });
            this.roomInput= "";
        },
        leaveRoom(room) {
            this.ws.send(JSON.stringify({
                action: 'leave-room',
                message: room.name,
            }));

            for (let i = 0; i < this.rooms.length ; i++) {
                if (this.rooms[i].name === room.name) {
                    this.rooms.splice(i,1);
                    break;
                }
            }
        }
    }

})