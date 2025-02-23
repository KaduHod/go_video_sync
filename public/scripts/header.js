const header = document.getElementsByTagName('header')[0];
const headerHeight = header.offsetHeight;
const main = document.getElementsByTagName('main')[0];
window.header_height = headerHeight
main.style.paddingTop = `${headerHeight}px`;
