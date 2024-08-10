import {fetchWithAuth} from "./tokens.js";

async function getUser(event) {
    event.preventDefault();

    const phoneNumber = sessionStorage.getItem('phone_number');
    if (!phoneNumber) {
        return;
    }
    console.log(`http://localhost:8000/api/readers/${phoneNumber}`)

    try {
        const user = await getUserFromStorage(phoneNumber);
        if (user) {
            displayUser(user)
        }

        return null;
    } catch (error) {
        return `Error: ${error.message}`;
    }
}

function getUserFromStorage(phoneNumber) {
    return fetchWithAuth(`http://localhost:8000/api/readers/${phoneNumber}`, {
        method: 'GET',
    })
        .then((res) => {
            if (res.status === 200) {
                return res.json();
            }
            return null;
        });
}

function displayUser(user) {
    document.getElementById('fio').textContent = user.fio;
    document.getElementById('phone_number').textContent = user.phone_number;
    document.getElementById('age').textContent = user.age;
}

document.addEventListener('DOMContentLoaded', getUser);
