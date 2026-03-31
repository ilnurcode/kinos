export function formatPrice(price) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0
  }).format(price)
}

export function formatDate(date) {
  if (!date) return ''
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  }).format(new Date(date))
}

export function formatDateTime(date) {
  if (!date) return ''
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(new Date(date))
}

export function truncate(str, length = 50) {
  if (!str) return ''
  if (str.length <= length) return str
  return str.slice(0, length) + '...'
}

export function formatPhone(phone) {
  if (!phone) return ''
  const digits = phone.replace(/\D/g, '')
  if (digits.length === 11 && digits.startsWith('7')) {
    return '+7 (' + digits.substring(1, 4) + ') ' + 
           digits.substring(4, 7) + '-' + 
           digits.substring(7, 9) + '-' + 
           digits.substring(9, 11)
  }
  return phone
}
