const AUTH_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

function saveTokens(accessToken) {
    localStorage.setItem(AUTH_TOKEN_KEY, accessToken);
    updateNavbar();
}

function getAccessToken() {
    return localStorage.getItem(AUTH_TOKEN_KEY);
}

function clearTokens() {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    updateNavbar();
}

function isAuthenticated() {
    return !!getAccessToken();
}

function getRoleFromToken() {
    const token = getAccessToken();
    if (!token) return null;

    try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        return payload.role;
    } catch (e) {
        return null;
    }
}

function updateNavbar() {
    const navLogin = document.getElementById('nav-login');
    if (!navLogin) return;

    const navRegister = document.getElementById('nav-register');
    const navProfile = document.getElementById('nav-profile');
    const navAdmin = document.getElementById('nav-admin');
    const navLogout = document.getElementById('nav-logout');

    if (isAuthenticated()) {
        navLogin.classList.add('d-none');
        navRegister.classList.add('d-none');
        navProfile.classList.remove('d-none');
        navLogout.classList.remove('d-none');

        if (getRoleFromToken() === 'admin') {
            navAdmin?.classList.remove('d-none');
        } else {
            navAdmin?.classList.add('d-none');
        }
    } else {
        navLogin.classList.remove('d-none');
        navRegister.classList.remove('d-none');
        navProfile.classList.add('d-none');
        navAdmin?.classList.add('d-none');
        navLogout.classList.add('d-none');
    }
}

document.addEventListener('DOMContentLoaded', updateNavbar);

document.addEventListener('click', (e) => {
    const logoutBtn = e.target.closest('#logoutBtn');
    if (logoutBtn) {
        e.preventDefault();
        clearTokens();
        window.location.href = '/';
    }
});