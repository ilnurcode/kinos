import { storage } from "@/utils/storage";

const configuredBaseURL = import.meta.env.VITE_API_BASE_URL || "/api";

function resolveBaseURL() {
  if (configuredBaseURL.startsWith("http://") || configuredBaseURL.startsWith("https://")) {
    return configuredBaseURL.replace(/\/+$/, "");
  }

  if (typeof window !== "undefined") {
    return new URL(configuredBaseURL, window.location.origin).toString().replace(/\/+$/, "");
  }

  return configuredBaseURL.replace(/\/+$/, "");
}

const baseURL = resolveBaseURL();

function normalizePath(path) {
  if (typeof path !== "string") {
    return "";
  }

  return path.replace(/^\/+/, "");
}

let refreshPromise = null;

function buildUrl(path, params = {}) {
  const url = new URL(normalizePath(path), `${baseURL}/`);

  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") {
      return;
    }

    url.searchParams.set(key, String(value));
  });

  return url.toString();
}

function createHttpError(message, response, config) {
  const error = new Error(message);
  error.name = "HttpError";
  error.response = response;
  error.config = config;
  return error;
}

async function parseResponseBody(response) {
  const contentType = response.headers.get("content-type") || "";

  if (contentType.includes("application/json")) {
    return response.json();
  }

  const text = await response.text();
  return text ? { message: text } : null;
}

async function refreshAccessToken() {
  if (!refreshPromise) {
    refreshPromise = request("/users/refresh", {
      method: "POST",
      skipAuth: true,
      skipRetry: true,
    })
      .then((response) => {
        const newToken = response.data?.access_token;
        if (!newToken) {
          throw new Error("Missing access token in refresh response");
        }

        storage.set("access_token", newToken);
        return newToken;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }

  return refreshPromise;
}

async function request(path, options = {}) {
  const {
    method = "GET",
    params,
    data,
    headers = {},
    skipAuth = false,
    skipRetry = false,
    _retry = false,
  } = options;

  const requestConfig = {
    method,
    params,
    data,
    headers,
    skipAuth,
    skipRetry,
    _retry,
    path,
  };

  const finalHeaders = new Headers(headers);
  const token = storage.get("access_token");

  if (!skipAuth && token) {
    finalHeaders.set("Authorization", `Bearer ${token}`);
  }

  if (data !== undefined && !(data instanceof FormData) && !finalHeaders.has("Content-Type")) {
    finalHeaders.set("Content-Type", "application/json");
  }

  const fetchOptions = {
    method,
    headers: finalHeaders,
    credentials: "include",
  };

  if (data !== undefined) {
    fetchOptions.body = data instanceof FormData ? data : JSON.stringify(data);
  }

  let response;

  try {
    response = await fetch(buildUrl(path, params), fetchOptions);
  } catch (error) {
    throw createHttpError(error.message || "Network error", undefined, requestConfig);
  }

  const responseData = await parseResponseBody(response);
  const responseShape = {
    status: response.status,
    data: responseData,
    headers: response.headers,
  };

  if (response.status === 401 && !skipRetry && !_retry) {
    try {
      const newToken = await refreshAccessToken();
      return request(path, {
        ...options,
        _retry: true,
        headers: {
          ...headers,
          Authorization: `Bearer ${newToken}`,
        },
      });
    } catch (refreshError) {
      storage.remove("access_token");
      window.location.href = "/login";
      throw refreshError;
    }
  }

  if (!response.ok) {
    const message = responseData?.error || responseData?.message || `Request failed with status ${response.status}`;
    throw createHttpError(message, responseShape, requestConfig);
  }

  return responseShape;
}

const apiClient = {
  get(path, options = {}) {
    return request(path, { ...options, method: "GET" });
  },

  post(path, data, options = {}) {
    return request(path, { ...options, method: "POST", data });
  },

  put(path, data, options = {}) {
    return request(path, { ...options, method: "PUT", data });
  },

  delete(path, options = {}) {
    return request(path, { ...options, method: "DELETE" });
  },
};

export default apiClient;
