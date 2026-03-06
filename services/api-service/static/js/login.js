document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('loginForm');
    const message = document.getElementById('message');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const btn = form.querySelector('button');
        btn.disabled = true;
        message.style.display = 'none';

        try {
            const response = await apiRequest('/api/users/login', {
                method: 'POST',
                body: JSON.stringify({
                    email: form.email.value.trim(),
                    password: form.password.value.trim()
                })
            });

            if (!response.access_token) {
                throw new Error('Сервер не вернул токен');
            }

            saveTokens(response.access_token);
            showMessage('success', 'Вход выполнен!');
            setTimeout(() => window.location.href = '/profile', 1500);
        } catch (err) {
            let errorMessage = 'Неверный email или пароль';
            if (err.status === 500) {
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