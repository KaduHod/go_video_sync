window.player;
function setUp() {
    var firstScriptTag = document.getElementsByTagName('script')[0];
    var tag = document.createElement('script');
    window.player_done = false
    tag.src = "https://www.youtube.com/iframe_api";
    firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);
    console.log("Player :: Preparando...")
}
window.setUpPlayer = setUp
function onYouTubePlayerAPIReady() {
    window.player = new YT.Player('ytplayer', {
        height: '360',
        width: '640',
        videoId: '',
        events: {
            "onReady": onPlayerReady,
            "onStateChange": onPlayerStateChange
        }
    });
}
function isValidYouTubeUrl(url) {
    try {
        // Cria um objeto URL
        const parsedUrl = new URL(url);

        // Verifica o domínio
        const validDomains = [
            "youtube.com",
            "www.youtube.com",
            "youtu.be"
        ];

        // Verifica se o domínio é válido
        const isYouTubeDomain = validDomains.includes(parsedUrl.hostname) ||
                                parsedUrl.hostname.endsWith(".youtube.com");

        if (!isYouTubeDomain) {
            return false; // Domínio não é do YouTube
        }

        // Verifica o caminho da URL
        if (parsedUrl.hostname === "youtu.be") {
            // URLs encurtadas do YouTube (youtu.be)
            return parsedUrl.pathname.length > 1; // Deve ter um ID após o /
        } else if (parsedUrl.pathname === "/watch") {
            // URLs de vídeos no YouTube (youtube.com/watch?v=...)
            return parsedUrl.searchParams.has("v"); // Deve ter o parâmetro "v"
        } else {
            // Outros caminhos válidos (canais, playlists, etc.)
            return true;
        }
    } catch (error) {
        return false; // URL inválida
    }
}
window.isValidYouTubeUrl = isValidYouTubeUrl
function getYouTubeVideoId(url) {
    try {
        const parsedUrl = new URL(url);

        if (parsedUrl.hostname === "youtu.be") {
            return parsedUrl.pathname.slice(1);
        }

        if (parsedUrl.hostname.includes("youtube.com")) {
            return parsedUrl.searchParams.get("v");
        }

        return null;
    } catch (error) {
        return null;
    }
}
window.getYouTubeVideoId = getYouTubeVideoId
function onPlayerReady(event) {
    window.user.player = "ready"
    console.log("Player :: Iniciado")

    //window.player.loadVideoById(getYouTubeVideoId("https://www.youtube.com/watch?v=fO5UGknokwc&ab_channel=AugustoGalego"), 0);
}
function changeVideo(url) {
    window.sendRoomAction("player::change-video", {...window.user, url})
    console.log("Player :: Trocando video para", url)
}
window.changeVideo = changeVideo
function onPlayerStateChange(event) {
    if (event.data == YT.PlayerState.PLAYING && !window.player_done) {
        setTimeout(window.player.stopVideo(), 6000);
        window.player_done = true;
    }
    const state = event.data

    switch (state) {
        case YT.PlayerState.ENDED:
            break;
        case YT.PlayerState.PAUSED:
            break;
        case YT.PlayerState.PLAYING:
            break;
        default:
            break;
    }
}
