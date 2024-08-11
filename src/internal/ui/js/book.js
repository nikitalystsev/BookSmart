import {fetchWithAuth} from "./tokens.js";

function displaySelectedBook() {
    const selectedBook = JSON.parse(sessionStorage.getItem('selectedBook'));

    if (!selectedBook) {
        document.getElementById('book-container').innerHTML = '<p>Книга не найдена.</p>';
        return;
    }

    const {copies_number, publisher, age_limit, rarity, title, author, genre, language, publishing_year} = selectedBook;
    document.getElementById('book-title').textContent = title;
    document.getElementById('book-author').textContent = author;
    document.getElementById('book-publisher').textContent = publisher || 'Нет данных';
    document.getElementById('book-copies-number').textContent = copies_number || 'Нет данных';
    document.getElementById('book-rarity').textContent = rarity || 'Нет данных';
    document.getElementById('book-genre').textContent = genre || 'Нет данных';
    document.getElementById('book-publishing-year').textContent = publishing_year || 'Нет данных';
    document.getElementById('book-language').textContent = language || 'Нет данных';
    document.getElementById('book-age-limit').textContent = age_limit || 'Нет данных';

}

async function reserveSelectedBook(event) {
    event.preventDefault()

    const selectedBook = JSON.parse(sessionStorage.getItem('selectedBook'));

    if (!selectedBook) {
        document.getElementById('book-container').innerHTML = '<p>Книга не найдена.</p>';
        return;
    }

    try {
        let response = await reserveBookOnStorage(selectedBook.id);
        if (!response.ok) {
            console.error(`HTTP error! Status: ${response.status}`);
            return response.text();
        }

        return null;
    } catch (error) {
        return `Error: ${error.message}`;
    }
}

async function reserveBookOnStorage(bookID) {
    return await fetchWithAuth("http://localhost:8000/api/reservations/", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(bookID)
    });
}

async function reserveSelectedBookWithMessage(event) {
    event.preventDefault();

    const message = await reserveSelectedBook(event)

    const messageElement = document.getElementById('message');
    if (message === null) {
        messageElement.className = 'alert alert-success';
        messageElement.textContent = 'Книга была успешно забронирована!';
    } else {
        messageElement.className = 'alert alert-danger';
        messageElement.textContent = message;
    }

    messageElement.classList.remove('d-none');
}

function addButtonsIfAuthenticated() {
    const isAuthenticated = sessionStorage.getItem("isAuthenticated");

    const btnContainer = document.getElementById('book-btn');

    if (!isAuthenticated) return

    const reserveBookBtn = document.createElement('a');
    reserveBookBtn.href = '#';
    reserveBookBtn.id = 'reserve-book-btn'
    reserveBookBtn.className = 'btn btn-primary mt-3';
    reserveBookBtn.innerHTML = '<i class="fas fa-calendar-check"></i> Забронировать';

    const addBookToFavoriteBtn = document.createElement('a');
    addBookToFavoriteBtn.href = '#';
    addBookToFavoriteBtn.className = 'btn btn-secondary mt-3';
    addBookToFavoriteBtn.innerHTML = '<i class="fas fa-heart"></i> Добавить в избранное';

    btnContainer.appendChild(reserveBookBtn);
    btnContainer.appendChild(addBookToFavoriteBtn);

    reserveBookBtn.addEventListener("click", reserveSelectedBookWithMessage)
}


// Вызываем функцию при загрузке страницы
document.addEventListener('DOMContentLoaded', addButtonsIfAuthenticated);
document.addEventListener('DOMContentLoaded', displaySelectedBook);
