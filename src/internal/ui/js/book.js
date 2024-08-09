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

document.addEventListener('DOMContentLoaded', displaySelectedBook);
