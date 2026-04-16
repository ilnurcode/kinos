import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { cartApi } from "@/api/cart";
import { ordersApi } from "@/api/orders";

function normalizeCart(rawCart) {
  if (!rawCart) {
    return { items: [], total: 0 };
  }

  const items = (rawCart.items || []).map((item) => ({
    productId: item.product_id ?? item.productId,
    productName: item.product_name ?? item.productName,
    quantity: item.quantity ?? 1,
    price: item.price ?? 0,
    subtotal: item.subtotal ?? (item.price ?? 0) * (item.quantity ?? 1),
  }));

  return {
    items,
    total: rawCart.total ?? 0,
  };
}

export const useCartStore = defineStore("cart", () => {
  const items = ref([]);
  const total = ref(0);
  const loading = ref(false);
  const ordering = ref(false);
  const error = ref(null);
  const lastOrder = ref(null);
  const count = ref(0);

  const hasItems = computed(() => items.value.length > 0);
  const itemsCount = computed(() =>
    items.value.reduce((sum, item) => sum + Number(item.quantity || 0), 0)
  );

  function applyCart(payload) {
    const normalized = normalizeCart(payload?.cart ?? payload);
    items.value = normalized.items;
    total.value = payload?.total ?? normalized.total;
    count.value = itemsCount.value;
  }

  async function fetchCart() {
    loading.value = true;
    error.value = null;

    try {
      const response = await cartApi.getCart();
      applyCart(response);
      return response;
    } catch (err) {
      error.value = err?.response?.data?.error || err.message || "Ошибка загрузки корзины";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchItemsCount() {
    try {
      const response = await cartApi.getItemsCount();
      count.value = Number(response?.count ?? 0);
      return count.value;
    } catch (err) {
      count.value = 0;
      return 0;
    }
  }

  async function addItem(productId, quantity = 1) {
    loading.value = true;
    error.value = null;

    try {
      const response = await cartApi.addItem(productId, quantity);
      applyCart(response);
      return response;
    } catch (err) {
      error.value = err?.response?.data?.error || err.message || "Ошибка добавления в корзину";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updateItem(productId, quantity) {
    loading.value = true;
    error.value = null;

    try {
      const response = await cartApi.updateItem(productId, quantity);
      applyCart(response);
      return response;
    } catch (err) {
      error.value = err?.response?.data?.error || err.message || "Ошибка обновления корзины";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function removeItem(productId) {
    loading.value = true;
    error.value = null;

    try {
      const response = await cartApi.removeItem(productId);
      applyCart(response);
      return response;
    } catch (err) {
      error.value = err?.response?.data?.error || err.message || "Ошибка удаления из корзины";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function clearCart() {
    loading.value = true;
    error.value = null;

    try {
      await cartApi.clearCart();
      items.value = [];
      total.value = 0;
      count.value = 0;
    } catch (err) {
      error.value = err?.response?.data?.error || err.message || "Ошибка очистки корзины";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createOrder(payload) {
    ordering.value = true;
    error.value = null;

    try {
      const response = await ordersApi.createOrder(payload);
      lastOrder.value = response?.order ?? null;
      items.value = [];
      total.value = 0;
      count.value = 0;
      return response;
    } catch (err) {
      error.value = err?.response?.data?.error || err.message || "Ошибка оформления заказа";
      throw err;
    } finally {
      ordering.value = false;
    }
  }

  return {
    items,
    total,
    count,
    loading,
    ordering,
    error,
    lastOrder,
    hasItems,
    itemsCount,
    fetchCart,
    fetchItemsCount,
    addItem,
    updateItem,
    removeItem,
    clearCart,
    createOrder,
  };
});
