export const storage = {
  get(key) {
    try {
      return localStorage.getItem(key);
    } catch (e) {
      return null;
    }
  },

  set(key, value) {
    try {
      localStorage.setItem(key, value);
    } catch (e) {
      // Тихая ошибка записи
    }
  },

  remove(key) {
    try {
      localStorage.removeItem(key);
    } catch (e) {
      // Тихая ошибка удаления
    }
  },
};
