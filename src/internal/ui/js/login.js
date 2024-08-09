document.getElementById('myForm').addEventListener('submit', async function handleSubmit(event) {
    event.preventDefault(); // Предотвращает отправку формы по умолчанию

    // Собираем данные формы
    const formData = new FormData(this);
    const data = {};

    formData.forEach((value, key) => {
        data[key] = value;
    });

    try {
        // Отправляем данные в JSON формате на указанный URL
        const response = await fetch("http://localhost:8000/auth/sign-in", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        // Обработка ответа сервера
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const result = await response.json();
        console.log('Success:', result);
    } catch (error) {
        console.error('Error:', error);
    }
});
