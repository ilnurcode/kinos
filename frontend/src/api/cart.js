import apiClient from "@/api/client";

export const cartApi = {
  async getCart() {
    const response = await apiClient.get("/cart");
    return response.data;
  },

  async addItem(productId, quantity = 1) {
    const response = await apiClient.post("/cart/items", {
      product_id: productId,
      quantity,
    });
    return response.data;
  },

  async updateItem(productId, quantity) {
    const response = await apiClient.put(`/cart/items/${productId}`, { quantity });
    return response.data;
  },

  async removeItem(productId) {
    const response = await apiClient.delete(`/cart/items/${productId}`);
    return response.data;
  },

  async clearCart() {
    const response = await apiClient.post("/cart/clear");
    return response.data;
  },

  async getItemsCount() {
    const response = await apiClient.get("/cart/count");
    return response.data;
  },
};
