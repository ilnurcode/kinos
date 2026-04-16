import apiClient from "./client";

export const productsApi = {
  async getList(params = {}) {
    const response = await apiClient.get("/catalog/products", { params });
    return response.data;
  },

  async getCategories() {
    const response = await apiClient.get("/catalog/categories");
    return response.data;
  },

  async getManufacturers() {
    const response = await apiClient.get("/catalog/manufacturers");
    return response.data;
  },
};
