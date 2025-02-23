window.appendMessage = function(name, text, time = null) {
    const messagesContainer = document.getElementById("messages-container");
    const messageTime = time || new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    const messageDiv = document.createElement("div");
    messageDiv.classList.add("bg-gray-500", "text-white", "p-2", "rounded-lg", "max-w-xs");
    const messageText = document.createElement("div");
    messageText.textContent = text;
    const footerDiv = document.createElement("div");
    footerDiv.classList.add("text-sm", "text-gray-300", "mt-1", "flex", "justify-between");
    const timeSpan = document.createElement("span");
    timeSpan.textContent = messageTime;
    const nameSpan = document.createElement("span");
    nameSpan.textContent = name;
    footerDiv.appendChild(timeSpan);
    footerDiv.appendChild(nameSpan);
    messageDiv.appendChild(messageText);
    messageDiv.appendChild(footerDiv);
    messagesContainer.appendChild(messageDiv);
}
const sendMessageButton = document.getElementById("send-message");
const sendMessage = (e) => {
    const messageInput = document.getElementById("chat-input");
    window.sendRoomAction("chat::message", {message: messageInput.value})
}
sendMessageButton.addEventListener("click", sendMessage);
window.handleChatEvent = function (data) {
    const { value } = data;
    switch (value) {
        case "chat::message":
            window.appendMessage(data.sender.name, data.meta.message);
            break;
        default:
            console.log(value, "Não reconhecido");
            break;
    }
};
