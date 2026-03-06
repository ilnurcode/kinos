document.addEventListener('DOMContentLoaded', async () => {
    const productsGrid = document.getElementById('productsGrid');
    const loading = document.getElementById('loading');
    const emptyState = document.getElementById('emptyState');
    const totalProductsEl = document.getElementById('totalProducts');

    const filterCategory = document.getElementById('filterCategory');
    const filterManufacturer = document.getElementById('filterManufacturer');
    const filterPriceMin = document.getElementById('filterPriceMin');
    const filterPriceMax = document.getElementById('filterPriceMax');
    const filterName = document.getElementById('filterName');
    const applyFiltersBtn = document.getElementById('applyFilters');
    const resetFiltersBtn = document.getElementById('resetFilters');
    const sortOrder = document.getElementById('sortOrder');

    let categories = [];
    let manufacturers = [];
    let currentProducts = [];

    async function loadDictionaries() {
        try {
            const [catData, manData] = await Promise.all([
                apiRequest('/api/catalog/categories?limit=100&offset=0', { method: 'GET' }),
                apiRequest('/api/catalog/manufacturers?limit=100&offset=0', { method: 'GET' })
            ]);
            categories = catData.category || [];
            manufacturers = manData.manufacturer || [];

            filterCategory.innerHTML = '<option value="">Все</option>' +
                categories.map(c => `<option value="${c.id}">${c.name}</option>`).join('');

            filterManufacturer.innerHTML = '<option value="">Все</option>' +
                manufacturers.map(m => `<option value="${m.id}">${m.name}</option>`).join('');
        } catch (err) {
            console.error('Ошибка загрузки справочников:', err);
        }
    }

    async function loadProducts() {
        loading.style.display = 'block';
        productsGrid.style.display = 'none';
        emptyState.style.display = 'none';

        const params = new URLSearchParams({ limit: 100, offset: 0 });
        if (filterCategory.value) params.append('category_id', filterCategory.value);
        if (filterManufacturer.value) params.append('manufacturer_id', filterManufacturer.value);
        if (filterPriceMin.value) params.append('price_min', filterPriceMin.value);
        if (filterPriceMax.value) params.append('price_max', filterPriceMax.value);
        if (filterName.value) params.append('name_contains', filterName.value);

        try {
            const data = await apiRequest(`/api/catalog/products?${params.toString()}`, { method: 'GET' });
            currentProducts = data.product || [];
            
            // Применяем сортировку
            applySort();
            
            renderProducts(currentProducts);
            totalProductsEl.textContent = data.total || currentProducts.length;
        } catch (err) {
            console.error('Failed to load products:', err);
            productsGrid.innerHTML = '<div class="col-12 text-center text-danger">Ошибка загрузки товаров</div>';
            productsGrid.style.display = 'block';
        } finally {
            loading.style.display = 'none';
        }
    }

    // Применение сортировки
    function applySort() {
        const sortValue = sortOrder.value;
        switch (sortValue) {
            case 'price_asc':
                currentProducts.sort((a, b) => a.price - b.price);
                break;
            case 'price_desc':
                currentProducts.sort((a, b) => b.price - a.price);
                break;
            case 'name_asc':
                currentProducts.sort((a, b) => a.name.localeCompare(b.name));
                break;
            case 'name_desc':
                currentProducts.sort((a, b) => b.name.localeCompare(a.name));
                break;
            default:
                // По умолчанию - как пришло с сервера
                break;
        }
    }

    // Рендер товаров
    function renderProducts(products) {
        if (products.length === 0) {
            emptyState.style.display = 'block';
            productsGrid.style.display = 'none';
            return;
        }

        emptyState.style.display = 'none';
        productsGrid.style.display = 'flex';

        // Создаём мапу для быстрого поиска названий
        const catMap = Object.fromEntries(categories.map(c => [c.id, c.name]));
        const manMap = Object.fromEntries(manufacturers.map(m => [m.id, m.name]));

        productsGrid.innerHTML = products.map(p => `
            <div class="col-md-6 col-lg-4">
                <div class="card h-100 shadow-sm product-card">
                    <div class="card-body d-flex flex-column">
                        <h5 class="card-title">${escapeHtml(p.name)}</h5>
                        <div class="product-info flex-grow-1">
                            <p class="text-muted mb-1">
                                <small>Категория: ${catMap[p.category_id] || '—'}</small>
                            </p>
                            <p class="text-muted mb-1">
                                <small>Производитель: ${manMap[p.manufacturer_id] || '—'}</small>
                            </p>
                        </div>
                        <div class="d-flex justify-content-between align-items-center mt-3">
                            <span class="fs-5 fw-bold text-primary">${p.price.toFixed(2)} ₽</span>
                            <button class="btn btn-sm btn-outline-primary" onclick="openProductModal('${escapeHtml(p.name)}', ${p.price}, '${escapeHtml(catMap[p.category_id] || '—')}', '${escapeHtml(manMap[p.manufacturer_id] || '—')}')">
                                Подробнее
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `).join('');
    }

    // Экранирование HTML
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // Обработчики
    applyFiltersBtn.addEventListener('click', () => {
        loadProducts();
    });

    resetFiltersBtn.addEventListener('click', () => {
        filterCategory.value = '';
        filterManufacturer.value = '';
        filterPriceMin.value = '';
        filterPriceMax.value = '';
        filterName.value = '';
        sortOrder.value = 'default';
        loadProducts();
    });

    sortOrder.addEventListener('change', () => {
        applySort();
        renderProducts(currentProducts);
    });

    // Поиск по Enter
    filterName.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            loadProducts();
        }
    });

    // Открытие модального окна товара
    window.openProductModal = function(name, price, category, manufacturer) {
        const modalElement = document.getElementById('productModal');
        const modalContent = document.getElementById('productModalContent');
        const modalTitle = document.getElementById('productModalLabel');

        modalTitle.textContent = name;
        modalContent.innerHTML = `
            <div class="product-details">
                <div class="mb-3">
                    <h6 class="text-muted mb-1">Цена</h6>
                    <p class="fs-4 fw-bold text-primary mb-0">${price.toFixed(2)} ₽</p>
                </div>
                <div class="mb-3">
                    <h6 class="text-muted mb-1">Категория</h6>
                    <p class="mb-0">${category}</p>
                </div>
                <div class="mb-3">
                    <h6 class="text-muted mb-1">Производитель</h6>
                    <p class="mb-0">${manufacturer}</p>
                </div>
                <div class="mb-3">
                    <h6 class="text-muted mb-1">Описание</h6>
                    <p class="text-muted mb-0">Описание будет добавлено...</p>
                </div>
            </div>
        `;

        const modal = new bootstrap.Modal(modalElement);
        modal.show();
    }

    // Инициализация
    await loadDictionaries();
    await loadProducts();
});
