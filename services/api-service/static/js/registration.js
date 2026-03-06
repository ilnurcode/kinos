document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('registerForm');
    const message = document.getElementById('message');
    const phoneInput = document.getElementById('phone');

    // Маска телефона: +7 XXX XXX XX XX
    phoneInput.addEventListener('input', (e) => {
        let value = e.target.value.replace(/\D/g, '');

        if (value.length > 11) {
            value = value.substring(0, 11);
        }

        if (value.length > 1) {
            value = '+7 ' + value.substring(1, 4) +
                    (value.length > 4 ? ' ' + value.substring(4, 7) : '') +
                    (value.length > 7 ? ' ' + value.substring(7, 9) : '') +
                    (value.length > 9 ? ' ' + value.substring(9, 11) : '');
        }

        e.target.value = value;
    });

    // Перед отправкой удаляем пробелы из номера
    form.addEventListener('submit', (e) => {
        phoneInput.value = phoneInput.value.replace(/\D/g, '');
    });

    if (!form) {
        console.error('Форма registration не найдена');
        return;
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const btn = form.querySelector('button');
        btn.disabled = true;
        message.style.display = 'none';

        const data = {
            username: form.username.value.trim(),
            email: form.email.value.trim(),
            password: form.password.value.trim(),
            phone: '+7' + form.phone.value.trim()
        };

        try {
            const response = await apiRequest('/api/users/register', {
                method: 'POST',
                body: JSON.stringify(data)
            });

            if (!response.access_token) {
                throw new Error('Сервер не вернул токен');
            }

            saveTokens(response.access_token);
            showMessage('success', 'Регистрация успешна!');
            setTimeout(() => window.location.href = '/', 1500);
        } catch (err) {
            let errorMessage = 'Ошибка при регистрации';
            if (err.status === 400) {
                errorMessage = 'Некорректные данные';
            } else if (err.status === 409) {
                errorMessage = 'Пользователь с таким email уже существует';
            } else if (err.status === 500) {
                errorMessage = 'Ошибка сервера';
            } else if (err.name === 'ApiError') {
                errorMessage = err.message;
            }
            showMessage('danger', errorMessage);
            btn.disabled = false;
        }
    });

    function showMessage(type, text) {
        message.className = `alert alert-${type}`;
        message.textContent = text;
        message.style.display = 'block';
    }
});