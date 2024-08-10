async function getBooks(event) {
    event.preventDefault();

    const searchParams = parseParams();

    const books = await getBooksFromStorage(searchParams);
    if (books != null) {
        displayBooks(books);
        sessionStorage.setItem("searchParams", JSON.stringify(searchParams))
        sessionStorage.setItem('currPageNum', "1");
        sessionStorage.setItem("1", JSON.stringify(books));
    }
}

function parseParams() {
    const form = document.getElementById('paramsForm');
    let searchParams = {}
    const title = form.elements.title,
        author = form.elements.author,
        publisher = form.elements.publisher,
        copiesNumber = form.elements.copies_number,
        rarity = form.elements.rarity,
        genre = form.elements.genre,
        publishingYear = form.elements.publishing_year,
        language = form.elements.language,
        ageLimit = form.elements.age_limit;

    if (title) searchParams.title = title.value;
    if (author) searchParams.author = author.value;
    if (publisher) searchParams.publisher = publisher.value;
    if (copiesNumber) searchParams.copies_number = parseInt(copiesNumber.value);
    if (genre) searchParams.genre = genre.value;
    if (publishingYear) searchParams.publishing_year = parseInt(publishingYear.value);
    if (language) searchParams.language = language.value;
    if (ageLimit) searchParams.age_limit = parseInt(ageLimit.value);
    if (rarity) {
        if (rarity.value === 'Обычная') searchParams.rarity = 'Common';
        if (rarity.value === 'Редкая') searchParams.rarity = 'Rare';
        if (rarity.value === 'Уникальная') searchParams.rarity = 'Unique';
    }

    searchParams.limit = 10;
    searchParams.offset = 0;

    return searchParams;
}

async function nextPageBooks(event) {
    console.log('Next page');
    event.preventDefault();

    if (!sessionStorage.getItem('currPageNum')) {
        return
    }
    const currPageNum = parseInt(sessionStorage.getItem('currPageNum'));
    const newPageNum = currPageNum + 1;
    await getPageBooks(newPageNum);
    updatePagination();
}

async function prevPageBooks(event) {
    event.preventDefault();

    if (!sessionStorage.getItem('currPageNum')) {
        return
    }
    const currPageNum = parseInt(sessionStorage.getItem('currPageNum'));
    if (currPageNum === 1) {
        return
    }
    const newPageNum = currPageNum - 1;
    await getPageBooks(newPageNum);
    updatePagination();
}

async function getPageBooks(newPageNum) {
    if (sessionStorage.getItem(newPageNum.toString())) {
        const books = JSON.parse(sessionStorage.getItem(newPageNum.toString()));
        sessionStorage.setItem('currPageNum', newPageNum.toString());
        displayBooks(books);
        return;
    }

    let searchParams;
    if (!sessionStorage.getItem("searchParams")) {
        searchParams = {'limit': 10, 'offset': 0}
    } else {
        searchParams = JSON.parse(sessionStorage.getItem("searchParams"));
    }
    searchParams['offset'] = (newPageNum - 1) * 10;

    const books = await getBooksFromStorage(searchParams);
    if (books != null) {
        displayBooks(books);
        sessionStorage.setItem('currPageNum', newPageNum.toString());
        sessionStorage.setItem(newPageNum.toString(), JSON.stringify(books));
        sessionStorage.setItem("searchParams", JSON.stringify(searchParams));
    }
}

async function getBooksFromStorage(searchParams) {
    const response = await fetch("http://localhost:8000/general/books", {
        method: 'POST', headers: {
            'Content-Type': 'application/json'
        }, body: JSON.stringify(searchParams)
    });

    if (!response.ok) {
        console.error(`HTTP error! Status: ${response.status}`);
        return null;
    }

    return await response.json();
}

function displayBooks(books) {
    const bookCardsContainer = document.getElementById('book-cards');
    bookCardsContainer.innerHTML = '';

    if (books.length === 0) {
        bookCardsContainer.innerHTML = '<p>Книги не найдены.</p>';
        return;
    }

    books.forEach(book => {
        const card = document.createElement('div');
        card.className = 'col-md-4 mb-4';

        card.innerHTML = `
            <div class="card h-100" style="border-radius: 10px; box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);">
                <h5 class="card-header">${book.title}</h5>
                <div class="card-body">
                    <h5 class="card-title">${book.author}</h5>
                    <p class="card-text" style="max-width: 150px;">${book.description || 'Нет описания.'}</p>
                    <a href="#" class="btn btn-primary" onclick='choiceBook(${JSON.stringify(book)})'>Подробнее</a>
                </div>
            </div> `;

        bookCardsContainer.appendChild(card);
    });
}


function choiceBook(book) {
    sessionStorage.setItem('selectedBook', JSON.stringify(book));
    window.location.href = '../templates/book.html';
}

function displayPageBooks() {
    const pageNumber = sessionStorage.getItem('currPageNum');
    if (!pageNumber) {
        return;
    }

    const books = JSON.parse(sessionStorage.getItem(pageNumber));

    displayBooks(books);
}

function updatePagination() {
    const pagination = document.getElementById('pagination');
    let currentPage;
    if (!sessionStorage.getItem('currPageNum')) {
        currentPage = 1;
    } else {
        currentPage = parseInt(sessionStorage.getItem('currPageNum'));
    }
    const pageItems = Array.from(pagination.getElementsByClassName('page-item'));
    pageItems.slice(1, -1).forEach(item => item.remove());

    let startPage, endPage;
    // обработка начального случая
    if (currentPage === 1) {
        document.getElementById('prevPageBtn').disabled = true;
        startPage = 1;
        endPage = 3;
    } else {
        document.getElementById('prevPageBtn').disabled = false;
        startPage = currentPage - 1;
        endPage = currentPage + 1;
    }

    // Добавляем номера страниц
    for (let i = startPage; i <= endPage; i++) {
        const pageItem = document.createElement('li');
        pageItem.className = 'page-item' + (i === currentPage ? ' active' : '');

        const pageLink = document.createElement('a');
        pageLink.className = 'page-link';
        pageLink.href = '#';
        pageLink.textContent = i;

        pageItem.appendChild(pageLink);
        pagination.insertBefore(pageItem, pagination.children[pagination.children.length - 1]);
    }

}

function setActiveNavLink(activeLinkId) {
    const links = document.querySelectorAll('.nav-link');
    links.forEach(link => {
        link.classList.remove('active');
    });
    const activeLink = document.getElementById(activeLinkId);
    if (activeLink) {
        activeLink.classList.add('active');
    }
}

function setCatalogNavbar() {
    const isAuthenticated = sessionStorage.getItem('isAuthenticated') === 'true'; // Проверка статуса аутентификации

    document.getElementById('navbar-container').innerHTML = `
        <nav class="navbar navbar-expand-lg navbar-light bg-light shadow-sm">
                <a class="navbar-brand" href="#"><b>BookSmart</b></a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"
                        aria-controls="navbarNav" aria-expanded="false" aria-label="Переключить навигацию">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarNav">
                    <ul class="navbar-nav me-auto">
                        <li class="nav-item">
                            <a class="nav-link" id="nav-home" href="index.html">Главная</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link" id="nav-catalog" href="catalog.html">Каталог</a>
                        </li>
                    </ul>
                </div>
                <div id="auth-links">
                    ${isAuthenticated
        ? '<a href="../templates/profile.html" class="btn btn-outline-primary ms-2">Личный кабинет</a>'
        : '<a href="../templates/login.html" class="btn btn-outline-primary ms-2">Войти</a>'}
                </div>
        </nav> `
    ;

    // Устанавливаем активную ссылку
    setActiveNavLink('nav-catalog'); // Замените 'nav-catalog' на нужный ID, если требуется
}

// Вызываем функцию при загрузке страницы
window.onload = setCatalogNavbar;

// Инициализация пагинации при загрузке страницы
updatePagination();
displayPageBooks()

document.getElementById('paramsForm').addEventListener("submit", getBooks);
document.getElementById('nextPageBtn').addEventListener("click", nextPageBooks);
document.getElementById('prevPageBtn').addEventListener("click", prevPageBooks);
document.getElementById('prevPageBtn').addEventListener("click", prevPageBooks);