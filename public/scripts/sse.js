const pauseButton = document.getElementById('pauseButton');
const resumeButton = document.getElementById('resumeButton');
const roomName = window.room_name;
const userName = window.user_name;
const csrf = document.getElementById('csrf').value;
resumeButton.addEventListener("click", () => {
    console.log("Aqui")
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
        const res = await req.json();
        console.log({res})
    } catch (error) {
        console.log(error)
    }

}
const eventSource = new EventSource('/sse/room/' + roomName);
