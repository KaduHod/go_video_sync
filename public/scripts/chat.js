window.appendSystemMessage = function(text) {
    const messagesContainer = document.getElementById("messages-container");
    const messageDiv = document.createElement("div");
    messageDiv.classList.add("bg-gray-60","border-b", "border-gray-400", "text-white", "w-full", "jsutify-center", "p-1", "text-xs", "text-center", "font-semibold");
    messageDiv.textContent = text + " " + new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    messagesContainer.appendChild(messageDiv);
}
function scrollToBottom() {
    const messagesContainer = document.getElementById("messages-container");
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}
window.appendMessage = function(name, text, time = null) {
    const messagesContainer = document.getElementById("messages-container");
    const messageTime = time || new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    const messageDiv = document.createElement("div");
    messageDiv.classList.add(name === window.user_name ? "self-end" : "self-start", "text-sm");
    messageDiv.classList.add(name === window.user_name ? "bg-purple-600" : "bg-red-500", "text-white", "p-2", "rounded-lg", "max-w-xs");
    const messageText = document.createElement("div");
    messageText.textContent = text;
    const footerDiv = document.createElement("div");
    footerDiv.classList.add("text-xs", "text-gray-300", "mt-1", "flex", "justify-between");
    const timeSpan = document.createElement("span");
    timeSpan.textContent = messageTime;
    timeSpan.classList.add("mr-2");
    const nameSpan = document.createElement("span");
    nameSpan.textContent = name;
    footerDiv.appendChild(timeSpan);
    footerDiv.appendChild(nameSpan);
    messageDiv.appendChild(messageText);
    messageDiv.appendChild(footerDiv);
    messagesContainer.appendChild(messageDiv);
}
const sendMessageButton = document.getElementById("send-message");
const sendMessage = () => {
    const messageInput = document.getElementById("chat-input");
    if(messageInput.value == "") return
    window.sendRoomAction("chat::message", {message: messageInput.value})
    messageInput.value = ""
}
sendMessageButton.addEventListener("click", sendMessage);
window.handleChatEvent = function (data) {
    const { value } = data;
    switch (value) {
        case "chat::message":
            if(data.meta.message == "") return
            window.appendMessage(data.sender.name, data.meta.message);
            scrollToBottom()
            break;
        default:
            console.log(value, "Não reconhecido");
            break;
    }
};
