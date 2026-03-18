document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('registerForm');
    const message = document.getElementById('message');
    const phoneInput = document.getElementById('phone');

    if (!form) {
        console.error('Форма registration не найдена');
        return;
    }

    // Маска телефона: +7 (___) ___-__-__
    phoneInput.addEventListener('input', (e) => {
        // Удаляем все нецифровые символы
        let digits = e.target.value.replace(/\D/g, '');

        // Если пользователь начинает вводить с 8, заменяем на 7
        if (digits.startsWith('8')) {
            digits = '7' + digits.substring(1);
        }

        // Добавляем 7 в начало если нет
        if (!digits.startsWith('7')) {
            digits = '7' + digits;
        }

        // Ограничиваем 11 цифрами (7 + 10 цифр номера)
        if (digits.length > 11) {
            digits = digits.substring(0, 11);
        }

        // Форматируем: +7 (XXX) XXX-XX-XX
        let formatted = '+7';
        if (digits.length > 1) {
            formatted += ' (' + digits.substring(1, 4);
        }
        if (digits.length > 4) {
            formatted += ') ' + digits.substring(4, 7);
        }
        if (digits.length > 7) {
            formatted += '-' + digits.substring(7, 9);
        }
        if (digits.length > 9) {
            formatted += '-' + digits.substring(9, 11);
        }

        e.target.value = formatted;
    });

    // Разрешаем удалять символы
    phoneInput.addEventListener('keydown', (e) => {
        if (e.key === 'Backspace' || e.key === 'Delete') {
            setTimeout(() => {
                let value = phoneInput.value;
                if (value === '+7') {
                    phoneInput.value = '+7 (';
                }
            }, 10);
        }
    });

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const btn = form.querySelector('button');
        btn.disabled = true;
        message.style.display = 'none';

        // Получаем только цифры из телефона
        const phoneDigits = phoneInput.value.replace(/\D/g, '');

        // Формируем телефон в формате E.164
        let phoneE164 = '+' + phoneDigits;

        const data = {
            username: form.username.value.trim(),
            email: form.email.value.trim(),
            password: form.password.value.trim(),
            phone: phoneE164
        };

        console.log('📤 Отправка данных:', data);
        console.log('📱 Телефон:', data.phone, 'длина:', data.phone.length);

        try {
            const response = await apiRequest('/api/users/register', {
                method: 'POST',
                body: JSON.stringify(data)
            });

            console.log('📥 Ответ сервера:', response);

            if (!response.access_token) {
                throw new Error('Сервер не вернул токен');
            }

            saveTokens(response.access_token);
            showMessage('success', 'Регистрация успешна!');
            setTimeout(() => window.location.href = '/', 1500);
        } catch (err) {
            console.error('❌ Ошибка регистрации:', err);
            let errorMessage = 'Ошибка при регистрации';
            if (err.status === 400) {
                errorMessage = 'Некорректные данные: ' + (err.message || '');
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