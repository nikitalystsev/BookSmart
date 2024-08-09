function setActiveNavHome() {
    fetch('navbar.html')
        .then(response => response.text())
        .then(data => {
            document.getElementById('navbar-container').innerHTML = data;
            setActiveNavLink('nav-home');
        });
}

document.addEventListener("DOMContentLoaded", setActiveNavHome);