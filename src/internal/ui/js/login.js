async function loginUser(event) {
    event.preventDefault();

    const user = parseLogin();

    sessionStorage.setItem('phone_number', user.phone_number);

    try {
        let response = await loginUserOnStorage(user);
        if (!response.ok) {
            console.error(`HTTP error! Status: ${response.status}`);
            return response.text();
        }

        const tokens = await response.json();
        sessionStorage.setItem('tokens', JSON.stringify(tokens));
        sessionStorage.setItem('isAuthenticated', "true");

        return null;
    } catch (error) {
        return `Error: ${error.message}`;
    }
}

function parseLogin() {
    const form = document.getElementById('loginForm');
    let userData = {}
    const phoneNumber = form.elements.phone_number,
        password = form.elements.password;

    if (phoneNumber) userData.phone_number = phoneNumber.value;
    if (password) userData.password = password.value;

    return userData;
}

async function loginUserOnStorage(userData) {
    return await fetch("http://localhost:8000/auth/sign-in", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData)
    });
}

async function loginUserWithMessage(event) {
    event.preventDefault(); // Предотвращаем стандартное поведение отправки формы

    const message = await loginUser(event)

    const messageElement = document.getElementById('message');
    if (message === null) {
        messageElement.className = 'alert alert-success'; // Успех
        messageElement.textContent = 'Вход прошел успешно!';
        window.location.href = '../templates/index.html';
    } else {
        messageElement.className = 'alert alert-danger'; // Ошибка
        messageElement.textContent = message;
    }

    messageElement.classList.remove('d-none'); // Показываем сообщение
}


// document.getElementById('loginForm').addEventListener('submit', loginUser);
document.getElementById('loginForm').addEventListener('submit', loginUserWithMessage);