document.addEventListener('DOMContentLoaded', async () => {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return;
    }

    const form = document.getElementById('editProfileForm');
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

    // Перед отправкой удаляем пробелы
    form.addEventListener('submit', (e) => {
        phoneInput.value = phoneInput.value.replace(/\D/g, '');
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

        try {
            const response = await apiRequest('/api/profile', {
                method: 'PUT',
                body: JSON.stringify({
                    username: form.username.value.trim(),
                    email: form.email.value.trim(),
                    phone: '+7' + form.phone.value.trim()
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