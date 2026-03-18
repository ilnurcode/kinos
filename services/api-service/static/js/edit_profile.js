document.addEventListener('DOMContentLoaded', async () => {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return;
    }

    const form = document.getElementById('editProfileForm');
    const message = document.getElementById('message');
    const phoneInput = document.getElementById('phone');

    if (!form) {
        console.error('Форма редактирования профиля не найдена');
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

    try {
        const user = await apiRequest('/api/profile', { method: 'GET' });
        form.username.value = user.username || '';
        form.email.value = user.email || '';

        // Удаляем +7 из номера для отображения
        let phone = user.phone || '';
        if (phone.startsWith('+7')) {
            phone = phone.substring(2);
        }
        form.phone.value = phone;
    } catch (err) {
        showMessage('danger', err.message);
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const btn = form.querySelector('button');
        btn.disabled = true;

        // Получаем только цифры из телефона
        const phoneDigits = phoneInput.value.replace(/\D/g, '');

        // Формируем телефон в формате E.164
        let phoneE164 = '+' + phoneDigits;

        try {
            const response = await apiRequest('/api/profile', {
                method: 'PUT',
                body: JSON.stringify({
                    username: form.username.value.trim(),
                    email: form.email.value.trim(),
                    phone: phoneE164
                })
            });

            if (response.success) {
                showMessage('success', 'Профиль обновлён!');
                setTimeout(() => window.location.href = '/profile', 1500);
            } else {
                throw new Error('Ошибка обновления');
            }
        } catch (err) {
            showMessage('danger', err.message);
            btn.disabled = false;
        }
    });

    function showMessage(type, text) {
        message.className = `alert alert-${type}`;
        message.textContent = text;
        message.style.display = 'block';
    }
});
