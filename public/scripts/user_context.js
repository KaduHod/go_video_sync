function user() {
    const type = window.user_type.toLowerCase() == "admin" ? "admin" : "guest"
    const player = "waiting"
    const sse = "waiting"
    const user_name = window.user_name
    return {type, player, sse, user_name}
}
window.user = new user()
