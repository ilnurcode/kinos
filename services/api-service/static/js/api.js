async function apiRequest(url, options = {}) {
    const accessToken = localStorage.getItem('access_token');

    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
    }

    options.headers = headers;
    options.credentials = options.credentials || 'include';
    options.method = options.method || 'GET';

    try {
        const response = await fetch(url, options);

        if (response.status === 401) {
            try {
                const refreshResp = await fetch('/api/users/refresh', {
                    method: 'POST',
                    credentials: 'include'
                });

                if (refreshResp.ok) {
                    const refreshData = await refreshResp.json();
                    if (refreshData.access_token) {
                        localStorage.setItem('access_token', refreshData.access_token);
                        options.headers['Authorization'] = `Bearer ${refreshData.access_token}`;
                        const retryResp = await fetch(url, options);
                        return handleResponse(retryResp);
                    }
                }
            } catch (e) {
                // Игнорируем ошибку refresh
            }

            clearTokens();
            window.location.href = '/login';
            return;
        }

        return handleResponse(response);
    } catch (error) {
        throw error;
    }
}

function handleResponse(response) {
    if (response.status === 204) {
        return {};
    }

    const contentType = response.headers.get('content-type');
    
    if (contentType && contentType.includes('application/json')) {
        return response.json().then(data => {
            if (!response.ok) {
                throw new ApiError(data?.error || data?.message || 'Ошибка запроса', response.status);
            }
            return data;
        });
    }

    return response.text().then(text => {
        if (response.ok) {
            return text ? { message: text } : {};
        }
        throw new ApiError('Ошибка сервера', response.status);
    });
}

class ApiError extends Error {
    constructor(message, status) {
        super(message);
        this.name = 'ApiError';
        this.status = status;
    }
}