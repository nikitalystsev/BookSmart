async function registerUser(event) {
    event.preventDefault();

    const user = parseRegistration();

    let res = await saveUserToStorage(user);
    if (res != null) {
    }
}

function parseRegistration() {
    const form = document.getElementById('loginForm');
    let userData = {}
    const phoneNumber = form.elements.phone_number,
        password = form.elements.password;

    if (phoneNumber) userData.phone_number = phoneNumber.value;
    if (password) userData.password = password.value;

    return userData;
}

async function saveUserToStorage(userData) {
    const response = await fetch("http://localhost:8000/auth/sign-in", {
        method: 'POST', headers: {
            'Content-Type': 'application/json'
        }, body: JSON.stringify(userData)
    });

    if (!response.ok) {
        console.error(`HTTP error! Status: ${response.status}`);
        return null;
    }

    return await response.json();
}

document.getElementById('loginForm').addEventListener('submit', registerUser);