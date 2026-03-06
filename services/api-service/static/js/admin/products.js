document.addEventListener('DOMContentLoaded', async () => {
    // Проверка: пользователь должен быть admin
    if (!isAuthenticated() || getRoleFromToken() !== 'admin') {
        alert('Доступ запрещён. Требуется роль администратора.');
        window.location.href = '/';
        return;
    }

    // Элементы
    const tbody = document.querySelector('#productsTable tbody');
    const loading = document.getElementById('loading');
    const createForm = document.getElementById('createProductForm');
    const prodCategorySelect = document.getElementById('prodCategory');
    const prodManufacturerSelect = document.getElementById('prodManufacturer');
    const filterCategory = document.getElementById('filterCategory');
    const filterManufacturer = document.getElementById('filterManufacturer');
    const filterPriceMin = document.getElementById('filterPriceMin');
    const filterPriceMax = document.getElementById('filterPriceMax');
    const filterName = document.getElementById('filterName');
    const applyFiltersBtn = document.getElementById('applyFilters');
    const editModal = new bootstrap.Modal(document.getElementById('editProductModal'));
    const editProdCategory = document.getElementById('editProdCategory');
    const editProdManufacturer = document.getElementById('editProdManufacturer');
    const saveBtn = document.getElementById('saveProductBtn');

    let categories = [];
    let manufacturers = [];
    let currentEditId = null;

    // Загрузка справочников
    async function loadDictionaries() {
        const [catData, manData] = await Promise.all([
            apiRequest('/api/catalog/categories?limit=100&offset=0', { method: 'GET' }),
            apiRequest('/api/catalog/manufacturers?limit=100&offset=0', { method: 'GET' })
        ]);
        categories = catData.category || [];
        manufacturers = manData.manufacturer || [];

        // Заполнить селекты
        const optionsCat = categories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
        prodCategorySelect.innerHTML = '<option value="">Выберите категорию</option>' + optionsCat;
        filterCategory.innerHTML = '<option value="">Все</option>' + optionsCat;
        editProdCategory.innerHTML = '<option value="">Выберите категорию</option>' + optionsCat;

        const optionsMan = manufacturers.map(m => `<option value="${m.id}">${m.name}</option>`).join('');
        prodManufacturerSelect.innerHTML = '<option value="">Выберите производителя</option>' + optionsMan;
        filterManufacturer.innerHTML = '<option value="">Все</option>' + optionsMan;
        editProdManufacturer.innerHTML = '<option value="">Выберите производителя</option>' + optionsMan;
    }

    // Загрузка товаров с фильтрами
    async function loadProducts() {
        const params = new URLSearchParams({
            limit: 100,
            offset: 0
        });
        if (filterCategory.value) params.append('category_id', filterCategory.value);
        if (filterManufacturer.value) params.append('manufacturer_id', filterManufacturer.value);
        if (filterPriceMin.value) params.append('price_min', filterPriceMin.value);
        if (filterPriceMax.value) params.append('price_max', filterPriceMax.value);
        if (filterName.value) params.append('name_contains', filterName.value);

        try {
            const data = await apiRequest(`/api/catalog/products?${params.toString()}`, { method: 'GET' });
            renderProducts(data.product || []);
        } catch (err) {
            console.error('Failed to load products:', err);
            tbody.innerHTML = '<tr><td colspan="6" class="text-center text-danger">Ошибка загрузки</td></tr>';
        } finally {
            loading.style.display = 'none';
            // Показываем таблицу
            document.getElementById('productsTable').style.display = 'table';
        }
    }

    function renderProducts(products) {
        if (products.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center">Нет товаров</td></tr>';
            return;
        }
        // Для отображения названий категории и производителя нужны справочники
        const catMap = Object.fromEntries(categories.map(c => [c.id, c.name]));
        const manMap = Object.fromEntries(manufacturers.map(m => [m.id, m.name]));

        tbody.innerHTML = products.map(p => `
            <tr>
                <td>${p.id}</td>
                <td>${p.name}</td>
                <td>${catMap[p.category_id] || p.category_id}</td>
                <td>${manMap[p.manufacturer_id] || p.manufacturer_id}</td>
                <td>${p.price.toFixed(2)}</td>
                <td>
                    <button class="btn btn-sm btn-warning edit-btn" data-id="${p.id}" data-name="${p.name}" data-cat="${p.category_id}" data-man="${p.manufacturer_id}" data-price="${p.price}">✏️</button>
                    <button class="btn btn-sm btn-danger delete-btn" data-id="${p.id}">🗑️</button>
                </td>
            </tr>
        `).join('');
    }

    // Создание товара
    createForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('prodName').value.trim();
        const categoryId = prodCategorySelect.value;
        const manufacturerId = prodManufacturerSelect.value;
        const price = parseFloat(document.getElementById('prodPrice').value);
        if (!name || !categoryId || !manufacturerId || isNaN(price)) return;

        try {
            await apiRequest('/api/admin/catalog/products', {
                method: 'POST',
                body: JSON.stringify({
                    name,
                    category_id: parseInt(categoryId),
                    manufacturer_id: parseInt(manufacturerId),
                    price
                })
            });
            createForm.reset();
            await loadProducts();
        } catch (err) {
            const msg = err.status === 400 ? 'Проверьте правильность заполнения' : 
                        err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка создания: ' + msg);
        }
    });

    // Открыть редактирование
    tbody.addEventListener('click', (e) => {
        const editBtn = e.target.closest('.edit-btn');
        if (editBtn) {
            currentEditId = editBtn.dataset.id;
            document.getElementById('editProdId').value = currentEditId;
            document.getElementById('editProdName').value = editBtn.dataset.name;
            document.getElementById('editProdCategory').value = editBtn.dataset.cat;
            document.getElementById('editProdManufacturer').value = editBtn.dataset.man;
            document.getElementById('editProdPrice').value = editBtn.dataset.price;
            editModal.show();
        }
    });

    // Сохранить изменения
    saveBtn.addEventListener('click', async () => {
        const name = document.getElementById('editProdName').value.trim();
        const categoryId = editProdCategory.value;
        const manufacturerId = editProdManufacturer.value;
        const price = parseFloat(document.getElementById('editProdPrice').value);
        if (!name || !categoryId || !manufacturerId || isNaN(price) || !currentEditId) return;

        try {
            await apiRequest(`/api/admin/catalog/products/${currentEditId}`, {
                method: 'PUT',
                body: JSON.stringify({
                    name,
                    category_id: parseInt(categoryId),
                    manufacturer_id: parseInt(manufacturerId),
                    price
                })
            });
            editModal.hide();
            await loadProducts();
        } catch (err) {
            const msg = err.status === 400 ? 'Проверьте правильность заполнения' : 
                        err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка обновления: ' + msg);
        }
    });

    // Удаление
    tbody.addEventListener('click', async (e) => {
        const deleteBtn = e.target.closest('.delete-btn');
        if (!deleteBtn) return;
        const id = deleteBtn.dataset.id;
        if (!confirm('Удалить товар?')) return;

        try {
            await apiRequest(`/api/admin/catalog/products/${id}`, {
                method: 'DELETE'
            });
            await loadProducts();
        } catch (err) {
            const msg = err.status === 500 ? 'Ошибка сервера' : err.message;
            alert('Ошибка удаления: ' + msg);
        }
    });

    // Применить фильтры
    applyFiltersBtn.addEventListener('click', () => {
        loading.style.display = 'block';
        loadProducts();
    });

    // Сброс фильтров
    const resetFiltersBtn = document.getElementById('resetFilters');
    if (resetFiltersBtn) {
        resetFiltersBtn.addEventListener('click', () => {
            filterCategory.value = '';
            filterManufacturer.value = '';
            filterPriceMin.value = '';
            filterPriceMax.value = '';
            filterName.value = '';
            loading.style.display = 'block';
            loadProducts();
        });
    }

    // Инициализация
    await loadDictionaries();
    await loadProducts();
});