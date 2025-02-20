const pauseButton = document.getElementById('pauseButton');
const resumeButton = document.getElementById('resumeButton');
const roomName = window.room_name;
const userName = window.user_name;
const closeRoom = document.getElementById('close-room');
const csrf = document.getElementById('csrf').value;
resumeButton.addEventListener("click", () => {
    sendRoomAction('resume')
})
async function sendRoomAction(action) {
    const formData = new URLSearchParams();
    formData.append('action', action);
    formData.append('user_name', userName);
    formData.append('csrf', csrf);

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
const eventSource = new EventSource('/sse/room/' + roomName);
eventSource.onmessage = function(event) {
    const msg = JSON.parse(event.data)
    console.log("Mensagem recebida >> ", {msg})
    if (msg.value == "close") {
        redirecionarUsuario("O stream da sala foi fechado pelo admin.");
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
    redirecionarUsuario("Erro ao acessar sala!");
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
