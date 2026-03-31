import apiClient from './client'

export const adminApi = {
  // Пользователи
  async getUsers(params = {}) {
    const response = await apiClient.get('/admin/users', { params })
    return response.data
  },

  async updateUserRole(userId, role) {
    const response = await apiClient.put('/admin/users/role', {
      user_id: userId,
      role
    })
    return response.data
  },

  async deleteUser(userId) {
    const response = await apiClient.delete(`/admin/users/${userId}`)
    return response.data
  },

  // Категории
  async getCategories(params = {}) {
    const response = await apiClient.get('/catalog/categories', { params })
    return response.data
  },

  async createCategory(name) {
    const response = await apiClient.post('/admin/catalog/categories', { name })
    return response.data
  },

  async updateCategory(id, name) {
    const response = await apiClient.put(`/admin/catalog/categories/${id}`, { name })
    return response.data
  },

  async deleteCategory(id) {
    const response = await apiClient.delete(`/admin/catalog/categories/${id}`)
    return response.data
  },

  // Производители
  async getManufacturers(params = {}) {
    const response = await apiClient.get('/catalog/manufacturers', { params })
    return response.data
  },

  async createManufacturer(name) {
    const response = await apiClient.post('/admin/catalog/manufacturers', { name })
    return response.data
  },

  async updateManufacturer(id, name) {
    const response = await apiClient.put(`/admin/catalog/manufacturers/${id}`, { name })
    return response.data
  },

  async deleteManufacturer(id) {
    const response = await apiClient.delete(`/admin/catalog/manufacturers/${id}`)
    return response.data
  },

  // Товары
  async getProducts(params = {}) {
    const response = await apiClient.get('/catalog/products', { params })
    return response.data
  },

  async createProduct(product) {
    const response = await apiClient.post('/admin/catalog/products', product)
    return response.data
  },

  async updateProduct(id, product) {
    const response = await apiClient.put(`/admin/catalog/products/${id}`, product)
    return response.data
  },

  async deleteProduct(id) {
    const response = await apiClient.delete(`/admin/catalog/products/${id}`)
    return response.data
  },

  // Запасы
  async getInventory(params = {}) {
    const response = await apiClient.get('/inventory/list', { params })
    return response.data
  },

  async createInventory(inventory) {
    const response = await apiClient.post('/inventory', inventory)
    return response.data
  },

  async updateInventory(id, inventory) {
    const response = await apiClient.put(`/inventory/${id}`, inventory)
    return response.data
  },

  async deleteInventory(id) {
    const response = await apiClient.delete(`/inventory/${id}`)
    return response.data
  },

  // Склады
  async getWarehouses(params = {}) {
    const response = await apiClient.get('/inventory/warehouses/list', { params })
    return response.data
  },

  async createWarehouse(warehouse) {
    const response = await apiClient.post('/inventory/warehouses', warehouse)
    return response.data
  },

  async updateWarehouse(id, warehouse) {
    const response = await apiClient.put(`/inventory/warehouses/${id}`, warehouse)
    return response.data
  },

  async deleteWarehouse(id) {
    const response = await apiClient.delete(`/inventory/warehouses/${id}`)
    return response.data
  }
}
