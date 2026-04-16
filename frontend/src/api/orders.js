import apiClient from "@/api/client";

export const ordersApi = {
  async createOrder(payload) {
    const response = await apiClient.post("/orders", payload);
    return response.data;
  },
};
