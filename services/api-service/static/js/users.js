document.addEventListener('DOMContentLoaded', async () => {
    if (!isAuthenticated() || getRoleFromToken() !== 'admin') {
        alert('Доступ запрещён');
        window.location.href = '/';
        return;
    }

    const tbody = document.querySelector('#usersTable tbody');
    const loading = document.getElementById('loading');

    try {
        const data = await apiRequest('/api/admin/users', { method: 'GET' });

        if (data.users && data.users.length > 0) {
            tbody.innerHTML = '';
            data.users.forEach(user => {
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${user.id}</td>
                    <td>${user.username}</td>
                    <td>${user.email}</td>
                    <td>${user.phone || '-'}</td>
                    <td>${user.role}</td>
                `;
                tbody.appendChild(row);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="5" class="text-center">Нет пользователей</td></tr>';
        }
    } catch (err) {
        console.error('Ошибка загрузки пользователей:', err);
        tbody.innerHTML = '<tr><td colspan="5" class="text-center text-danger">Ошибка загрузки</td></tr>';
    } finally {
        if (loading) loading.style.display = 'none';
        document.getElementById('usersTable').style.display = 'table';
    }

    // Форма изменения роли
    const form = document.getElementById('roleForm');
    const message = document.getElementById('message');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const userId = form.userId.value.trim();
        const newRole = form.newRole.value;

        if (!userId) {
            showMessage('danger', 'Введите ID пользователя');
            return;
        }

        const userIdNum = parseInt(userId);
        if (isNaN(userIdNum) || userIdNum <= 0) {
            showMessage('danger', 'ID должен быть положительным числом');
            return;
        }

        try {
            const resp = await apiRequest('/api/admin/users/role', {
                method: 'PUT',
                body: JSON.stringify({ user_id: userIdNum, role: newRole })
            });

            if (resp.success) {
                showMessage('success', 'Роль обновлена');
                form.userId.value = '';
                setTimeout(() => message.style.display = 'none', 3000);
            } else {
                showMessage('danger', 'Ошибка при обновлении');
            }
        } catch (err) {
            let msg = err.message;
            if (err.status === 403) {
                msg = 'Нельзя изменить роль самому себе';
            } else if (err.status === 404) {
                msg = 'Пользователь не найден';
            }
            showMessage('danger', msg);
        }
    });

    function showMessage(type, text) {
        message.className = `alert alert-${type}`;
        message.textContent = text;
        message.style.display = 'block';
    }
});