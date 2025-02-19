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
    console.log("Mensagem recebida:", event.data);
};

eventSource.onerror = function(error) {
    console.error("Erro na conexão SSE:", error);
};
async function closeRoomApi() {
    const formData = new URLSearchParams();
    formData.append('csrf', csrf);

    try {
        const req = await fetch(`/app/room/close/${roomName}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'X-CSRF-Token': csrf,
            },
            body: formData.toString(),
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
closeRoom.addEventListener("click", () => closeRoomApi())
