import {isBadRequest, isConflict, isInternalServerError, isNotFound} from "./errors.js";
import {fetchWithAuth} from "./tokens.js";

function displaySelectedReservation() {
    const selectedReservation = JSON.parse(sessionStorage.getItem('selectedReservation'));

    if (!selectedReservation) {
        document.getElementById('book-container').innerHTML = '<p>Книга не найдена.</p>';
        return;
    }

    const reservationStates = {
        "Issued": "Выдана",
        "Extended": "Продлена",
        "Expired": "Просрочена",
        "Closed": "Закрыта"
    };

    const {bookInfo, reservation} = selectedReservation;
    document.getElementById('book-title').textContent = bookInfo;
    document.getElementById('reservation-state').textContent = reservationStates[reservation.state];
    document.getElementById('reservation-state').textContent = reservationStates[reservation.state];
    document.getElementById('reservation-issue-date').textContent = reservation.issue_date.slice(0, 10);
    document.getElementById('reservation-return-date').textContent = reservation.return_date.slice(0, 10);
}

async function updateSelectedReservation(event) {
    event.preventDefault()

    const selectedReservation = JSON.parse(sessionStorage.getItem('selectedReservation'));
    const {bookInfo, reservation} = selectedReservation;
    if (!reservation) {
        document.getElementById('book-container').innerHTML = '<p>Книга не найдена.</p>';
        return;
    }

    try {
        let response = await updateSelectedReservationOnStorage(reservation.id);

        if (isBadRequest(response)) return "Ошибка запроса"
        if (isConflict(response)) return response.text()
        if (isNotFound(response)) return response.text()
        if (isInternalServerError(response)) return "Внутренняя ошибка сервера"

        return null;
    } catch (error) {
        return `Error: ${error.message}`;
    }
}

async function updateSelectedReservationOnStorage(reservationID) {
    return await fetchWithAuth(`http://localhost:8000/api/reservations/${reservationID}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(reservationID)
    });
}

async function updateSelectedReservationWithMessage(event) {
    console.log("а ну как нахуй, че за залупа")
    event.preventDefault();

    const message = await updateSelectedReservation(event)

    const messageElement = document.getElementById('message');
    if (message === null) {
        console.log("успех, ага бля")
        messageElement.className = 'alert alert-success';
        messageElement.textContent = 'Бронирование было успешно продлено!';
    } else {
        console.log("нет, мы ебемся в жопу")
        messageElement.className = 'alert alert-danger';
        messageElement.textContent = message;
    }

    messageElement.classList.remove('d-none');
}

document.addEventListener('DOMContentLoaded', displaySelectedReservation);
document.getElementById("update-reservation").addEventListener("click", updateSelectedReservationWithMessage)