document.addEventListener('DOMContentLoaded', async () => {
    // Проверка: пользователь должен быть admin
    if (!isAuthenticated() || getRoleFromToken() !== 'admin') {
        alert('Доступ запрещён. Требуется роль администратора.');
        window.location.href = '/';
        return;
    }

    const tbody = document.querySelector('#manufacturersTable tbody');
    const loading = document.getElementById('loading');
    const createForm = document.getElementById('createManufacturerForm');
    const editModal = new bootstrap.Modal(document.getElementById('editManufacturerModal'));
    const editForm = document.getElementById('editManufacturerForm');
    const saveBtn = document.getElementById('saveManufacturerBtn');

    let currentEditId = null;

    async function loadManufacturers() {
        try {
            const data = await apiRequest('/api/catalog/manufacturers?limit=100&offset=0', { method: 'GET' });
            renderManufacturers(data.manufacturer || []);
        } catch (err) {
            console.error('Failed to load manufacturers:', err);
            tbody.innerHTML = '<tr><td colspan="3" class="text-center text-danger">Ошибка загрузки</td></tr>';
        } finally {
            loading.style.display = 'none';
            // Показываем таблицу
            document.getElementById('manufacturersTable').style.display = 'table';
        }
    }

    function renderManufacturers(manufacturers) {
        if (manufacturers.length === 0) {
            tbody.innerHTML = '<tr><td colspan="3" class="text-center">Нет производителей</td></tr>';
            return;
        }
        tbody.innerHTML = manufacturers.map(m => `
            <tr>
                <td>${m.id}</td>
                <td>${m.name}</td>
                <td>
                    <button class="btn btn-sm btn-warning edit-btn" data-id="${m.id}" data-name="${m.name}">✏️</button>
                    <button class="btn btn-sm btn-danger delete-btn" data-id="${m.id}">🗑️</button>
                </td>
            </tr>
        `).join('');
    }

    createForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('manName').value.trim();
        if (!name) return;

        try {
            await apiRequest('/api/admin/catalog/manufacturers', {
                method: 'POST',
                body: JSON.stringify({ name })
            });
            document.getElementById('manName').value = '';
            await loadManufacturers();
        } catch (err) {
            const msg = err.status === 400 ? 'Проверьте правильность заполнения' : 
                        err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка создания: ' + msg);
        }
    });

    tbody.addEventListener('click', (e) => {
        const editBtn = e.target.closest('.edit-btn');
        if (editBtn) {
            currentEditId = editBtn.dataset.id;
            document.getElementById('editManId').value = currentEditId;
            document.getElementById('editManName').value = editBtn.dataset.name;
            editModal.show();
        }
    });

    saveBtn.addEventListener('click', async () => {
        const name = document.getElementById('editManName').value.trim();
        if (!name || !currentEditId) return;

        try {
            await apiRequest(`/api/admin/catalog/manufacturers/${currentEditId}`, {
                method: 'PUT',
                body: JSON.stringify({ name })
            });
            editModal.hide();
            await loadManufacturers();
        } catch (err) {
            const msg = err.status === 400 ? 'Проверьте правильность заполнения' : 
                        err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка обновления: ' + msg);
        }
    });

    tbody.addEventListener('click', async (e) => {
        const deleteBtn = e.target.closest('.delete-btn');
        if (!deleteBtn) return;
        const id = deleteBtn.dataset.id;
        if (!confirm('Удалить производителя?')) return;

        try {
            await apiRequest(`/api/admin/catalog/manufacturers/${id}`, {
                method: 'DELETE'
            });
            await loadManufacturers();
        } catch (err) {
            const msg = err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка удаления: ' + msg);
        }
    });

    await loadManufacturers();
});