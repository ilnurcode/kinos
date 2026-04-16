<template>
  <div class="cart-view py-4">
    <div class="container">
      <h1 class="mb-4">
        Корзина
      </h1>

      <div
        v-if="storeError"
        class="alert alert-danger"
      >
        {{ storeError }}
      </div>

      <div
        v-if="successMessage"
        class="alert alert-success"
      >
        {{ successMessage }}
      </div>

      <div
        v-if="loading"
        class="text-center py-5"
      >
        <LoadingSpinner text="Загрузка корзины..." />
      </div>

      <div
        v-else-if="!hasItems"
        class="text-center py-5 bg-light rounded"
      >
        <h3 class="mb-3">
          Корзина пуста
        </h3>
        <router-link
          class="btn btn-primary"
          to="/catalog"
        >
          Перейти в каталог
        </router-link>
      </div>

      <div
        v-else
        class="row g-4"
      >
        <div class="col-lg-8">
          <div class="card shadow-sm">
            <div class="card-body">
              <div
                v-for="item in items"
                :key="item.productId"
                class="d-flex flex-column flex-sm-row justify-content-between align-items-sm-center gap-3 py-3 border-bottom"
              >
                <div class="flex-grow-1">
                  <h5 class="mb-1">
                    {{ item.productName }}
                  </h5>
                  <p class="mb-0 text-muted">
                    {{ formatPrice(item.price) }} за шт.
                  </p>
                </div>

                <div class="d-flex align-items-center gap-2">
                  <button
                    class="btn btn-outline-secondary btn-sm"
                    :disabled="loading || item.quantity <= 1"
                    @click="changeQuantity(item, item.quantity - 1)"
                  >
                    -
                  </button>
                  <span class="fw-semibold">{{ item.quantity }}</span>
                  <button
                    class="btn btn-outline-secondary btn-sm"
                    :disabled="loading"
                    @click="changeQuantity(item, item.quantity + 1)"
                  >
                    +
                  </button>
                </div>

                <div class="text-sm-end">
                  <div class="fw-bold mb-1">
                    {{ formatPrice(item.subtotal) }}
                  </div>
                  <button
                    class="btn btn-link btn-sm text-danger p-0"
                    :disabled="loading"
                    @click="removeItem(item.productId)"
                  >
                    Удалить
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="col-lg-4">
          <div class="card shadow-sm mb-3">
            <div class="card-body">
              <h4 class="mb-3">
                Итого
              </h4>
              <div class="d-flex justify-content-between mb-3">
                <span>Товаров:</span>
                <strong>{{ itemsCount }}</strong>
              </div>
              <div class="d-flex justify-content-between fs-5 mb-3">
                <span>Сумма:</span>
                <strong>{{ formatPrice(total) }}</strong>
              </div>
              <button
                class="btn btn-outline-danger w-100"
                :disabled="loading"
                @click="clearCart"
              >
                Очистить корзину
              </button>
            </div>
          </div>

          <div class="card shadow-sm">
            <div class="card-body">
              <h4 class="mb-3">
                Оформление заказа
              </h4>

              <form @submit.prevent="submitOrder">
                <div class="mb-3">
                  <label
                    for="deliveryAddress"
                    class="form-label"
                  >
                    Адрес доставки
                  </label>
                  <input
                    id="deliveryAddress"
                    v-model.trim="orderForm.delivery_address"
                    type="text"
                    class="form-control"
                    required
                  >
                </div>

                <div class="mb-3">
                  <label
                    for="phone"
                    class="form-label"
                  >
                    Телефон
                  </label>
                  <input
                    id="phone"
                    v-model.trim="orderForm.phone"
                    type="tel"
                    class="form-control"
                    required
                  >
                </div>

                <div class="mb-3">
                  <label
                    for="comment"
                    class="form-label"
                  >
                    Комментарий
                  </label>
                  <textarea
                    id="comment"
                    v-model.trim="orderForm.comment"
                    class="form-control"
                    rows="3"
                  />
                </div>

                <button
                  type="submit"
                  class="btn btn-success w-100"
                  :disabled="ordering || loading"
                >
                  {{ ordering ? "Оформляем..." : "Заказать" }}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import LoadingSpinner from "@/components/common/LoadingSpinner.vue";
import { useCartStore } from "@/stores/cart";
import { formatPrice } from "@/utils/format";

const cartStore = useCartStore();
const successMessage = ref("");

const orderForm = reactive({
  delivery_address: "",
  phone: "",
  comment: "",
});

const items = computed(() => cartStore.items);
const total = computed(() => cartStore.total);
const itemsCount = computed(() => cartStore.itemsCount);
const hasItems = computed(() => cartStore.hasItems);
const loading = computed(() => cartStore.loading);
const ordering = computed(() => cartStore.ordering);
const storeError = computed(() => cartStore.error);

const changeQuantity = async (item, quantity) => {
  if (quantity < 1) {
    return;
  }

  await cartStore.updateItem(item.productId, quantity);
};

const removeItem = async (productId) => {
  await cartStore.removeItem(productId);
};

const clearCart = async () => {
  await cartStore.clearCart();
};

const submitOrder = async () => {
  successMessage.value = "";

  await cartStore.createOrder({
    delivery_address: orderForm.delivery_address,
    phone: orderForm.phone,
    comment: orderForm.comment,
  });

  orderForm.delivery_address = "";
  orderForm.phone = "";
  orderForm.comment = "";
  successMessage.value = "Заказ успешно оформлен";
};

onMounted(async () => {
  await cartStore.fetchCart();
});
</script>
