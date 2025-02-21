const roomName = window.room_name;
const userName = window.user_name;
const closeRoom = document.getElementById('close-room');
const csrf = document.getElementById('csrf').value;
const playPauseButton = document.getElementById('pausePlayButton')
let play = true;//"paused"
playPauseButton.addEventListener("click",({target}) => {
    play = !play
    if(play)
        sendRoomAction("player::pause")
    else
        sendRoomAction("player::resume")
})
const playPauseAnimation = play => {
    const bar1 = document.getElementById('play');
    const bar2 = document.getElementById('pause');
    if (play == "play") {
        bar1.classList.remove("hidden")
        bar2.classList.add("hidden")
    } else {
        bar1.classList.add("hidden");
        bar2.classList.remove("hidden");
    }
}
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
if(window.user.type == "admin") {
    const addVideoButton = document.getElementById('loadVideoButton')
    addVideoButton.addEventListener("click", (e) => {
        console.log({e})
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
            if (isGuest(data) && bothUsersReady()) {
                window.friend = data.meta
                deployInputToAddVideo()
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
const setUpCounter = () => {
    const counter = document.getElementById('counter-video')
    if (!window.player || !window.player.getDuration) return false
    let duration = window.player.getDuration();
    let tempoTotal = window.formatarTempo(duration);
    console.log({tempoTotal, duration})
    const incrementaCounter = () => {
        window.player_manager.time++
        if(tempoTotal == "00:00") {
            tempoTotal = window.formatarTempo(window.player.getDuration())
        }
        console.log(window.player_manager.time, tempoTotal)
        counter.innerText = window.formatarTempo(window.player_manager.time) + " / " + tempoTotal
    }
    return setInterval(incrementaCounter, 1000);
}
playPauseAnimation("pause")
window.setUpCounter = setUpCounter
const handlePlayerEvents = (data) => {
    const { value } = data;
    console.log(`Evento Player :: ${value} `, data);
    switch (value) {
        case "player::pause":
            window.player.pauseVideo();
            playPauseAnimation("pause");
            if (window.player_manager) {
                window.player_manager.state = "paused";
                clearInterval(window.player_manager.interval);
                window.player_manager.state = "waiting";
            }
            break;
        case "player::resume":
            window.player.playVideo();
            playPauseAnimation("play");
            if (window.player_manager) {
                if (window.player_manager.state === "waiting") {
                    window.player_manager.interval = setUpCounter();
                }
                window.player_manager.state = "playing";
            }
            break;
        case "player::change-video":
            window.player.loadVideoById(window.getYouTubeVideoId(data.meta.url), 0); // Altera o vídeo para o novo ID
            window.player.pauseVideo();
            if (window.player_manager) {
                window.player_manager.videoUrl = data.meta.url;

                window.player_manager.state = "waiting";
            }
            break;
        default:
            console.log(value, "Não reconhecido");
            break;
    }
};
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
