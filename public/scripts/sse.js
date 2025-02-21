
const pauseButton = document.getElementById('pauseButton');
const resumeButton = document.getElementById('resumeButton');
const roomName = window.room_name;
const userName = window.user_name;
const closeRoom = document.getElementById('close-room');
const csrf = document.getElementById('csrf').value;
pauseButton.addEventListener("click", () => {
    sendRoomAction("player::pause")
})
resumeButton.addEventListener("click", () => {
    sendRoomAction('player::resume')
})
async function sendRoomAction(action, meta) {
    const formData = new URLSearchParams();
    formData.append('action', action);
    formData.append('user_name', userName);
    formData.append('csrf', csrf);
    if(meta) {
        formData.append("meta", JSON.stringify(meta))
    }

    try {
        const req = await fetch(`/app/room/${roomName}/send`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData.toString(),
        })
        if (req.status != 200) {
            throw "Erro de requisicao"
        }
    } catch (error) {
        console.log(error)
    }
}
const addVideo = document.getElementById('add-video');
window.sendRoomAction = sendRoomAction
const eventSource = new EventSource('/sse/room/' + roomName);
eventSource.onmessage = function(event) {
    const data = JSON.parse(event.data)
    if (data.value == "close") {
        redirecionarUsuario("O stream da sala foi fechado pelo admin.");
    }
    if(data.value.includes("player::")) handlePlayerEvents(data)
    else handleRoomEvents(data)
};
const bothUsersReady = () => {
    const {friend, user} = window
    return friend.player == "ready" && user.player == "ready" && friend.sse == "ready" && user.sse == "ready"
}
const deployInputToAddVideo = () => {
    if(addVideo && addVideo.classList.contains("hidden")) {
        addVideo.classList.remove("hidden")
    }
}
if(addVideo) {
    const addVideoButton = document.getElementById('loadVideoButton')
    addVideoButton.addEventListener("click", () => {
        const videoUrl = document.getElementById('videoUrl').value
        if(!window.isValidYouTubeUrl(videoUrl)) {
            console.log("Url inválida")
            document.getElementById("videoUrl").value = ""
        }
        window.sendRoomAction("player::change-video", {...window.user, url: videoUrl})
    })
}
window.updateStateStatus = "disable"
window.updateState = () => {
    window.updateStateStatus = "enable"
    setInterval(() => {
        window.sendRoomAction("update::ping", window.user)
    }, 3000)
}
const isGuest = eventData => eventData.meta.type != window.user.type
const handleRoomEvents = (data) => {
    const {value} = data
    switch (value) {
        case "update::ping":
            if (isGuest(data)) {
                window.friend = data.meta
                if(bothUsersReady()) {
                    deployInputToAddVideo()
                }
            }
            break;
        case "server::ping":
            if (window.user.sse == "waiting") {
                window.user.sse = "ready"
                window.setUpPlayer()
                if (window.updateStateStatus == "disable")
                    window.updateState()
            }
            console.log("Dados", data)
            break;
        default:
            console.log("Tipo nao reconhecido", value)
            break;
    }
}
const handlePlayerEvents = (data) => {
    const {value} = data
    console.log(`Evento Player :: ${value} `, data)
    switch (value) {
        case "player::pause":
            window.player.pauseVideo();
            //console.log(window.player)
            break;
        case "player::resume":
            window.player.playVideo();
            //console.log(window.player)
            break;
        case "player::change-video":
            window.player.loadVideoById(window.getYouTubeVideoId(data.meta.url), 0);  // Altera o vídeo para o novo ID
            window.player.pauseVideo()
            break;
        default:
            console.log(value, "Não reconhecido")
            break;
    }
}
async function closeRoomApi() {
    try {
        const req = await fetch(`/app/room/close/${roomName}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'X-CSRF-Token': csrf,
            },
        })
        if (req.status != 200) {
            throw "Erro de requisicao"
        }
        const res = await req.text()
        console.log({res})
    } catch (error) {
        console.log(error)
    }
}
if (closeRoom)
    closeRoom.addEventListener("click", () => closeRoomApi())

eventSource.onerror = function (error) {
    console.log("Erro na conexão SSE:", error);
    redirecionarUsuario("Sala "+ window.room_name +" não está mais transmitindo!");
};

// Verificar periodicamente se a conexão foi fechada
const verificarConexao = setInterval(() => {
    if (eventSource.readyState === EventSource.CLOSED) {
        console.warn("Conexão SSE fechada.");
        redirecionarUsuario("Erro inexperado!");
    }
}, 3000); // Verifica a cada 3 segundos

function redirecionarUsuario(msg = "") {
    clearInterval(verificarConexao); // Para a verificação periódica
    eventSource.close(); // Garante que a conexão está fechada
    msg = msg != "" ? "?error=" + msg : "" ;
    window.location.href = `/app/user${msg}`;
}
window.sse = eventSource
