let inventoryData = [];
let productsData = [];
let warehousesData = [];

document.addEventListener('DOMContentLoaded', () => {
    loadInventory();
    loadProducts();
    loadWarehouses();
});

async function loadProducts() {
    try {
        const response = await apiRequest('/api/catalog/products?limit=100&offset=0', { method: 'GET' });
        productsData = response.product || [];

        // Заполняем select товаров в форме добавления
        const productSelect = document.getElementById('productId');
        if (productSelect) {
            productSelect.innerHTML = '<option value="">Выберите товар</option>';
            productsData.forEach(product => {
                const option = document.createElement('option');
                option.value = product.id;
                option.textContent = product.name;
                productSelect.appendChild(option);
            });
        }

        // Заполняем select товаров в фильтре
        const filterProductSelect = document.getElementById('filterProductId');
        if (filterProductSelect) {
            filterProductSelect.innerHTML = '<option value="">Все товары</option>';
            productsData.forEach(product => {
                const option = document.createElement('option');
                option.value = product.id;
                option.textContent = product.name;
                filterProductSelect.appendChild(option);
            });
        }

        // Перерисовываем таблицу если товары уже загружены
        if (inventoryData.length > 0) {
            renderTable();
        }
    } catch (err) {
        console.error('Ошибка загрузки товаров:', err);
    }
}

async function loadWarehouses() {
    try {
        const response = await apiRequest('/api/inventory/warehouses/list?limit=100&offset=0', { method: 'GET' });
        warehousesData = response.warehouses || [];

        // Заполняем select складов в форме добавления
        const warehouseSelect = document.getElementById('warehouseId');
        if (warehouseSelect) {
            warehouseSelect.innerHTML = '<option value="">Выберите склад</option>';
            warehousesData.forEach(warehouse => {
                const option = document.createElement('option');
                option.value = warehouse.id;
                option.textContent = `${warehouse.name} (${warehouse.city}, ${warehouse.street})`;
                warehouseSelect.appendChild(option);
            });
        }

        // Заполняем select складов в фильтре
        const filterWarehouseSelect = document.getElementById('filterWarehouseId');
        if (filterWarehouseSelect) {
            filterWarehouseSelect.innerHTML = '<option value="">Все склады</option>';
            warehousesData.forEach(warehouse => {
                const option = document.createElement('option');
                option.value = warehouse.id;
                option.textContent = `${warehouse.name} (${warehouse.city}, ${warehouse.street})`;
                filterWarehouseSelect.appendChild(option);
            });
        }
    } catch (err) {
        console.error('Ошибка загрузки складов:', err);
    }
}

async function loadInventory() {
    const table = document.getElementById('inventoryTable');
    table.innerHTML = '<tr><td colspan="7" class="text-center">Загрузка...</td></tr>';

    try {
        const response = await apiRequest('/api/inventory/list?limit=100&offset=0', { method: 'GET' });
        inventoryData = response.inventory || [];
        applyFilters();
    } catch (err) {
        table.innerHTML = `<tr><td colspan="7" class="text-center text-danger">Ошибка: ${err.message}</td></tr>`;
    }
}

function applyFilters() {
    const productId = document.getElementById('filterProductId').value;
    const warehouseId = document.getElementById('filterWarehouseId').value;
    const minQty = document.getElementById('filterMinQty').value;
    const maxQty = document.getElementById('filterMaxQty').value;

    let filtered = inventoryData;

    // Фильтр по товару
    if (productId) {
        filtered = filtered.filter(inv => inv.product_id == productId);
    }

    // Фильтр по складу
    if (warehouseId) {
        filtered = filtered.filter(inv => {
            const warehouse = warehousesData.find(w => w.id == warehouseId);
            return warehouse && inv.warehouse_location === warehouse.name;
        });
    }

    // Фильтр по минимальному количеству (available_quantity >= min)
    if (minQty !== '') {
        const min = parseInt(minQty);
        filtered = filtered.filter(inv => (inv.available_quantity || 0) >= min);
    }

    // Фильтр по максимальному количеству (available_quantity <= max)
    if (maxQty !== '') {
        const max = parseInt(maxQty);
        filtered = filtered.filter(inv => (inv.available_quantity || 0) <= max);
    }

    renderTable(filtered);
}

function renderTable(data = inventoryData) {
    const table = document.getElementById('inventoryTable');

    if (data.length === 0) {
        table.innerHTML = '<tr><td colspan="7" class="text-center">Нет данных</td></tr>';
        return;
    }

    table.innerHTML = data.map(inv => {
        // Находим название товара по ID
        const product = productsData.find(p => p.id === inv.product_id);
        const productName = product ? product.name : `ID: ${inv.product_id}`;

        return `
        <tr>
            <td>${inv.id}</td>
            <td>${productName}</td>
            <td>${inv.quantity || 0}</td>
            <td>${inv.reserved_quantity || 0}</td>
            <td>
                <span class="badge ${inv.available_quantity > 10 ? 'bg-success' : inv.available_quantity > 0 ? 'bg-warning' : 'bg-danger'}">
                    ${inv.available_quantity || 0}
                </span>
            </td>
            <td>${inv.warehouse_location}</td>
            <td>${formatDate(inv.updated_at)}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" onclick="openEditModal(${inv.id})">✏️</button>
                <button class="btn btn-sm btn-outline-danger" onclick="deleteInventory(${inv.id})">🗑️</button>
            </td>
        </tr>
        `;
    }).join('');
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function openCreateModal() {
    document.getElementById('modalTitle').textContent = 'Добавить запас';
    document.getElementById('inventoryId').value = '';
    document.getElementById('productId').value = '';
    document.getElementById('quantity').value = '0';
    document.getElementById('warehouseId').value = '';
}

function openEditModal(id) {
    const inventory = inventoryData.find(inv => inv.id === id);
    if (!inventory) return;

    document.getElementById('modalTitle').textContent = 'Редактировать запас';
    document.getElementById('inventoryId').value = inventory.id;
    document.getElementById('productId').value = inventory.product_id;
    document.getElementById('quantity').value = inventory.quantity;
    document.getElementById('warehouseId').value = inventory.warehouse_location;

    const modal = new bootstrap.Modal(document.getElementById('inventoryModal'));
    modal.show();
}

async function saveInventory() {
    const id = document.getElementById('inventoryId').value;
    const productId = document.getElementById('productId').value;
    const quantity = document.getElementById('quantity').value;
    const warehouseId = document.getElementById('warehouseId').value;

    if (!productId || !quantity || !warehouseId) {
        alert('Выберите товар, укажите количество и выберите склад');
        return;
    }

    // Находим название склада по ID
    const warehouse = warehousesData.find(w => w.id == warehouseId);
    const warehouseLocation = warehouse ? warehouse.name : 'Unknown';

    const data = {
        product_id: parseInt(productId),
        quantity: parseInt(quantity),
        warehouse_location: warehouseLocation
    };

    try {
        const url = id ? `/api/inventory/${id}` : '/api/inventory';
        const method = id ? 'PUT' : 'POST';

        await apiRequest(url, {
            method: method,
            body: JSON.stringify(data)
        });

        bootstrap.Modal.getInstance(document.getElementById('inventoryModal')).hide();
        loadInventory();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

async function deleteInventory(id) {
    if (!confirm('Удалить эту запись о запасе?')) return;

    try {
        await apiRequest(`/api/inventory/${id}`, { method: 'DELETE' });
        loadInventory();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}
