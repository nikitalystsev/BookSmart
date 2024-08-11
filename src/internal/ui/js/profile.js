import {fetchWithAuth} from "./tokens.js";

async function getUser(event) {
    event.preventDefault();

    const phoneNumber = sessionStorage.getItem('phone_number');
    if (!phoneNumber) {
        return;
    }

    try {
        const response = await getUserFromStorage(phoneNumber);
        if (!response.ok) {
            console.error(`HTTP error! Status: ${response.status}`);
            return response.text();
        }

        const user = await response.json();

        displayUser(user)

        return null;
    } catch (error) {
        return `Error: ${error.message}`;
    }
}

async function getUserFromStorage(phoneNumber) {
    return await fetchWithAuth(`http://localhost:8000/api/readers/${phoneNumber}`, {
        method: 'GET',
    })
}

function displayUser(user) {
    document.getElementById('fio').textContent = user.fio;
    document.getElementById('phone_number').textContent = user.phone_number;
    document.getElementById('age').textContent = user.age;
}

async function getUserWithMessage(event) {
    event.preventDefault();

    const message = await getUser(event)

    const messageElement = document.getElementById('message');

    if (message) {
        messageElement.className = 'alert alert-danger'; // Ошибка
        messageElement.textContent = message;
        messageElement.classList.remove('d-none');
    } else messageElement.classList.add('d-none');

}

document.addEventListener('DOMContentLoaded', getUserWithMessage);
