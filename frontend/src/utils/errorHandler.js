/**
 * Форматирует HTTP ошибки в человекочитаемый вид
 * @param {Error} error - Объект ошибки
 * @param {string} defaultMessage - Сообщение по умолчанию
 * @returns {string} Человекочитаемое сообщение об ошибке
 */
export function formatError(error, defaultMessage = 'Произошла ошибка') {
  // Нет ответа от сервера
  if (!error.response) {
    return 'Нет соединения с сервером. Проверьте подключение к интернету.'
  }

  const status = error.response.status
  const data = error.response.data

  // Ошибки по статус кодам
  const errorMessages = {
    400: data?.error || 'Некорректные данные. Проверьте заполнение полей.',
    401: 'Сессия истекла. Выполните вход заново.',
    403: 'Недостаточно прав для выполнения этого действия.',
    404: 'Ресурс не найден.',
    409: data?.error || 'Конфликт данных. Возможно, такая запись уже существует.',
    500: 'Ошибка сервера. Попробуйте позже или обратитесь к администратору.',
    502: 'Сервер временно недоступен.',
    503: 'Сервис временно недоступен.'
  }

  return errorMessages[status] || data?.error || error.message || defaultMessage
}
