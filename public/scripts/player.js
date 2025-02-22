function PlayerManager() {
    const videoUrl = "";
    const state = "waiting";
    const time = 0;
    const user = {};
    const getUser = function () {
        this.user = window.user
        return this.user
    }
    return {videoUrl, state, time, user, getUser}
}
window.player_manager = new PlayerManager()
window.seekToPlayer = function(event) {
    const time = event.target.dataset.time
    window.sendRoomAction("player::change-time", {time})
}
window.createSlider = function (breakpoints) {
    const sliderContainer = document.getElementById('slider-container');
    sliderContainer.innerHTML = ''; // Limpa qualquer conteúdo anterior

    // Criar a trilha (linha do slider)
    const track = document.createElement('div');
    track.className = 'w-full h-2 rounded-full relative';
    sliderContainer.appendChild(track);

    // Criar os breakpoints (pequenos traços)
    for (let i = 0; i <= breakpoints; i++) {
        const marker = document.createElement('div');
        marker.className = 'absolute w-2 h-2 rounded-full bg-gray-700 transform -translate-x-1/2 -top-1 cursor-pointer';

        // Definir a posição de cada marcador
        const position = (i / breakpoints) * 100;
        marker.style.left = `${position}%`;

        // Criar tooltip
        const tooltip = document.createElement('div');
        tooltip.className = 'absolute -top-8 left-1/2 -translate-x-1/2 px-2 py-2 bg-gray-900 text-white text-xs rounded opacity-0 transition-opacity duration-200 pointer-events-none';
        tooltip.innerText = window.formatarTempo(i);
        marker.id = "tooltip-" + i
        marker.dataset.time = i
        marker.addEventListener("click", window.seekToPlayer)

        // Eventos para exibir/esconder tooltip
        marker.addEventListener('mouseenter', () => {
            tooltip.classList.remove('opacity-0');
        });
        marker.addEventListener('mouseleave', () => {
            tooltip.classList.add('opacity-0');
        });

        marker.appendChild(tooltip);
        track.appendChild(marker);
    }
}
window.dismarkToolTip = function(i) {
    const after = document.getElementById("tooltip-" + (i + 1))
    if(after && after.classList.contains("bg-purple-900")) {
        window.dismarkToolTip(i + 1)
    }
    const tooltip = document.getElementById("tooltip-" + i)
    if(!tooltip) return
    tooltip.classList.remove("bg-purple-900")
    tooltip.classList.add("bg-gray-700")

}
window.markToolTip = function(i) {
    const before = document.getElementById("tooltip-" + (i - 1))
    if(before && before.classList.contains("bg-gray-700")) {
        window.markToolTip(i - 1)
    }
    const tooltip = document.getElementById("tooltip-" + i)
    if(!tooltip) return
    tooltip.classList.remove("bg-gray-700")
    tooltip.classList.add("bg-purple-900")
}
