export const storage = {
  get(key) {
    try {
      return localStorage.getItem(key)
    } catch (e) {
      return null
    }
  },

  set(key, value) {
    try {
      localStorage.setItem(key, value)
    } catch (e) {
      console.error('Storage error:', e)
    }
  },

  remove(key) {
    try {
      localStorage.removeItem(key)
    } catch (e) {
      console.error('Storage error:', e)
    }
  }
}
