let warehouseData = [];

document.addEventListener('DOMContentLoaded', () => {
    loadWarehouses();
});

async function loadWarehouses() {
    const table = document.getElementById('warehouseTable');
    table.innerHTML = '<tr><td colspan="5" class="text-center">Загрузка...</td></tr>';

    try {
        // Используем inventory API для складов
        const response = await apiRequest('/api/inventory/warehouses/list?limit=100&offset=0', { method: 'GET' });
        warehouseData = response.warehouses || [];
        renderTable();
    } catch (err) {
        // Если API еще нет, показываем заглушку
        warehouseData = [];
        renderTable();
    }
}

function renderTable() {
    const table = document.getElementById('warehouseTable');

    if (warehouseData.length === 0) {
        table.innerHTML = '<tr><td colspan="5" class="text-center">Нет данных. Создайте первый склад!</td></tr>';
        return;
    }

    table.innerHTML = warehouseData.map(w => `
        <tr>
            <td>${w.id}</td>
            <td><strong>${w.name}</strong></td>
            <td>${w.city || ''} ${w.street ? ', ' + w.street : ''}${w.building ? ', ' + w.building : ''}</td>
            <td>${w.name}</td>
            <td>
                <button class="btn btn-sm btn-outline-primary" onclick="openEditModal(${w.id})">✏️</button>
                <button class="btn btn-sm btn-outline-danger" onclick="deleteWarehouse(${w.id})">🗑️</button>
            </td>
        </tr>
    `).join('');
}

function openCreateModal() {
    document.getElementById('modalTitle').textContent = 'Добавить склад';
    document.getElementById('warehouseId').value = '';
    document.getElementById('warehouseName').value = '';
    document.getElementById('warehouseCity').value = '';
    document.getElementById('warehouseStreet').value = '';
    document.getElementById('warehouseBuilding').value = '';
    document.getElementById('warehouseBuilding2').value = '';
}

function openEditModal(id) {
    const warehouse = warehouseData.find(w => w.id === id);
    if (!warehouse) return;

    document.getElementById('modalTitle').textContent = 'Редактировать склад';
    document.getElementById('warehouseId').value = warehouse.id;
    document.getElementById('warehouseName').value = warehouse.name;
    document.getElementById('warehouseCity').value = warehouse.city || '';
    document.getElementById('warehouseStreet').value = warehouse.street || '';
    document.getElementById('warehouseBuilding').value = warehouse.building || '';
    document.getElementById('warehouseBuilding2').value = warehouse.building2 || '';

    const modal = new bootstrap.Modal(document.getElementById('warehouseModal'));
    modal.show();
}

async function saveWarehouse() {
    const id = document.getElementById('warehouseId').value;
    const name = document.getElementById('warehouseName').value;
    const city = document.getElementById('warehouseCity').value;
    const street = document.getElementById('warehouseStreet').value;
    const building = document.getElementById('warehouseBuilding').value;
    const building2 = document.getElementById('warehouseBuilding2').value;

    if (!name || !city || !street) {
        alert('Название, город и улица обязательны');
        return;
    }

    const data = {
        name: name,
        city: city,
        street: street,
        building: building,
        building2: building2
    };

    try {
        const url = '/api/inventory/warehouses';
        await apiRequest(url, {
            method: 'POST',
            body: JSON.stringify(data)
        });

        bootstrap.Modal.getInstance(document.getElementById('warehouseModal')).hide();
        loadWarehouses();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

async function deleteWarehouse(id) {
    if (!confirm('Удалить этот склад?')) return;

    try {
        await apiRequest(`/api/inventory/warehouses/${id}`, { method: 'DELETE' });
        loadWarehouses();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}