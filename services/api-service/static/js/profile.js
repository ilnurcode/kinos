document.addEventListener('DOMContentLoaded', async () => {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return;
    }

    try {
        const user = await apiRequest('/api/profile', { method: 'GET' });

        document.getElementById('username').textContent = user.username || '-';
        document.getElementById('email').textContent = user.email || '-';
        document.getElementById('phone').textContent = user.phone || '-';

        const role = getRoleFromToken();
        const actionsDiv = document.getElementById('profile-actions');
        if (actionsDiv) {
            // Исправлена ссылка на редактирование профиля
            actionsDiv.innerHTML = '<a href="/profile/edit" class="btn btn-primary me-2">Редактировать профиль</a>';
            if (role === 'admin') {
                actionsDiv.innerHTML += '<a href="/admin/users" class="btn btn-warning">Управление пользователями</a>';
            }
        }
    } catch (err) {
        const errorDiv = document.getElementById('error-message');
        if (errorDiv) {
            errorDiv.textContent = err.message;
            errorDiv.style.display = 'block';
        }
    }

    document.getElementById('logoutBtn')?.addEventListener('click', () => {
        clearTokens();
        window.location.href = '/';
    });
});