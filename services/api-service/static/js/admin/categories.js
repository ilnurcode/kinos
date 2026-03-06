// static/js/admin/categories.js
document.addEventListener('DOMContentLoaded', async () => {
    // Проверка: пользователь должен быть admin
    if (!isAuthenticated() || getRoleFromToken() !== 'admin') {
        alert('Доступ запрещён. Требуется роль администратора.');
        window.location.href = '/';
        return;
    }

    const tbody = document.querySelector('#categoriesTable tbody');
    const loading = document.getElementById('loading');
    const createForm = document.getElementById('createCategoryForm');
    const editModal = new bootstrap.Modal(document.getElementById('editCategoryModal'));
    const saveBtn = document.getElementById('saveCategoryBtn');
    let currentEditId = null;

    async function loadCategories() {
        try {
            const data = await apiRequest('/api/catalog/categories?limit=100&offset=0', { method: 'GET' });
            renderCategories(data.category || []);
        } catch (err) {
            console.error('Failed to load categories:', err);
            tbody.innerHTML = '<tr><td colspan="3" class="text-center text-danger">Ошибка загрузки</td></tr>';
        } finally {
            loading.style.display = 'none';
            // Показываем таблицу
            document.getElementById('categoriesTable').style.display = 'table';
        }
    }

    function renderCategories(categories) {
        if (categories.length === 0) {
            tbody.innerHTML = '<tr><td colspan="3" class="text-center">Нет категорий</td></tr>';
            return;
        }
        tbody.innerHTML = categories.map(c => `
            <tr>
                <td>${c.id}</td>
                <td>${c.name}</td>
                <td>
                    <button class="btn btn-sm btn-warning edit-btn" data-id="${c.id}" data-name="${c.name}">✏️</button>
                    <button class="btn btn-sm btn-danger delete-btn" data-id="${c.id}">🗑️</button>
                </td>
            </tr>
        `).join('');
    }

    createForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('catName').value.trim();
        if (!name) return;

        try {
            await apiRequest('/api/admin/catalog/categories', {
                method: 'POST',
                body: JSON.stringify({ name })
            });
            document.getElementById('catName').value = '';
            await loadCategories();
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
            document.getElementById('editCatId').value = currentEditId;
            document.getElementById('editCatName').value = editBtn.dataset.name;
            editModal.show();
        }
    });

    saveBtn.addEventListener('click', async () => {
        const name = document.getElementById('editCatName').value.trim();
        if (!name || !currentEditId) return;

        try {
            await apiRequest(`/api/admin/catalog/categories/${currentEditId}`, {
                method: 'PUT',
                body: JSON.stringify({ name })
            });
            editModal.hide();
            await loadCategories();
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
        if (!confirm('Удалить категорию?')) return;

        try {
            await apiRequest(`/api/admin/catalog/categories/${id}`, {
                method: 'DELETE'
            });
            await loadCategories();
        } catch (err) {
            const msg = err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка удаления: ' + msg);
        }
    });

    await loadCategories();
});