function displaySelectedBook() {
    const selectedBook = JSON.parse(sessionStorage.getItem('selectedBook'));

    if (!selectedBook) {
        document.getElementById('book-container').innerHTML = '<p>Книга не найдена.</p>';
        return;
    }

    document.getElementById('book-title').textContent = selectedBook.title;
    document.getElementById('book-author').textContent = selectedBook.author;
    document.getElementById('book-publisher').textContent = selectedBook.publisher || 'Нет данных';
    document.getElementById('book-copies-number').textContent = selectedBook.copies_number || 'Нет данных';
    document.getElementById('book-rarity').textContent = selectedBook.rarity || 'Нет данных';
    document.getElementById('book-genre').textContent = selectedBook.genre || 'Нет данных';
    document.getElementById('book-publishing-year').textContent = selectedBook.publishing_year || 'Нет данных';
    document.getElementById('book-language').textContent = selectedBook.language || 'Нет данных';
    document.getElementById('book-age-limit').textContent = selectedBook.age_limit || 'Нет данных';

}

function addButtonsIfAuthenticated() {
    // Предположим, что эта переменная устанавливается в зависимости от состояния авторизации пользователя
    const isAuthenticated = sessionStorage.getItem("isAuthenticated"); // Замените на реальную проверку авторизации пользователя

    // Находим контейнер кнопок по ID
    const btnContainer = document.getElementById('book-btn');

    if (isAuthenticated) {
        // Создаем кнопку "Забронировать"
        const bookButton = document.createElement('a');
        bookButton.href = '#'; // Замените на реальную ссылку для бронирования
        bookButton.className = 'btn btn-primary mt-3';
        bookButton.innerHTML = '<i class="fas fa-calendar-check"></i> Забронировать';

        // Создаем кнопку "Добавить в избранное"
        const favoriteButton = document.createElement('a');
        favoriteButton.href = '#'; // Замените на реальную ссылку для добавления в избранное
        favoriteButton.className = 'btn btn-secondary mt-3';
        favoriteButton.innerHTML = '<i class="fas fa-heart"></i> Добавить в избранное';

        // Убедитесь, что кнопки добавлены с отступами
        btnContainer.appendChild(bookButton);
        btnContainer.appendChild(favoriteButton);

        // Применение классов для обеспечения отступов
        btnContainer.classList.add('btn-container');
    }
}

// Вызываем функцию при загрузке страницы
document.addEventListener('DOMContentLoaded', addButtonsIfAuthenticated);
document.addEventListener('DOMContentLoaded', displaySelectedBook);
