window.popNotification = function (message, type) {
    alert(message)
}
document.addEventListener("DOMContentLoaded", () => {
    const header = document.getElementsByTagName('header')[0];
    const headerHeight = header.offsetHeight;
    const notificacoesContainer = document.getElementById("notifications");
    notificacoesContainer.style.paddingTop = `${headerHeight + 5}px`
})
function adicionarNotificacao(mensagem, tipo = "info") {
    const notificacoesContainer = document.getElementById("notifications");
    const notificacao = document.createElement("div");

    let bgColor;
    switch (tipo) {
        case "erro":
            bgColor = "bg-red-600";
            break;
        case "info":
            bgColor = "bg-purple-600";
            break;
        default:
            bgColor = "bg-blue-600";
            break;
    }

    notificacao.className = `${bgColor} text-white p-3 rounded-lg shadow-md transform translate-x-0 transition-opacity duration-500 opacity-100 relative flex justify-between items-center`;

    const texto = document.createElement("span");
    texto.innerText = mensagem;

    const botaoFechar = document.createElement("button");
    botaoFechar.innerText = "×";
    botaoFechar.className = "ml-4 text-white text-lg font-bold cursor-pointer";
    botaoFechar.addEventListener("click", (e) => e.target.parentNode.remove());

    notificacao.appendChild(texto);
    notificacao.appendChild(botaoFechar);
    notificacoesContainer.appendChild(notificacao);

    setTimeout(() => {
        notificacao.classList.add("opacity-0");
        setTimeout(() => notificacao.remove(), 500);
    }, 5000);
}
window.adicionarNotificacao = adicionarNotificacao
