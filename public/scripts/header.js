const header = document.getElementsByTagName('header')[0];
const headerHeight = header.offsetHeight;
const main = document.getElementsByTagName('main')[0];
console.log(headerHeight, main)
main.style.paddingTop = `${headerHeight}px`;
console.log("Aqui")
