// static/login.js

document.addEventListener('DOMContentLoaded', () => {
    console.log('login.js loaded');

    const form = document.getElementById('loginForm');
    const errorBox = document.getElementById('error-box');

    if (!form) {
        console.error('loginForm not found in DOM');
        return;
    }

    if (errorBox) {
        errorBox.style.display = 'none';
        errorBox.innerText = '';
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        console.log('login form submit');

        if (errorBox) {
            errorBox.style.display = 'none';
            errorBox.innerText = '';
        }

        const formData = new FormData(form);
        const data = Object.fromEntries(formData); // { email, password }
        console.log('sending login request', data.email);

        try {
            const res = await fetch('/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            const result = await res.json();
            console.log('login response', res.status, result);

            if (!res.ok) {
                const msg = result.error || 'Ошибка входа';
                if (errorBox) {
                    errorBox.style.display = 'block';
                    errorBox.innerText = msg;
                } else {
                    alert('❌ ' + msg);
                }
                return;
            }

            // ожидаем: { token, user_id, role, email, full_name, ... }
            const authData = {
                token: result.token,
                user: {
                    id: result.user_id,
                    role: result.role,
                    email: result.email,
                    full_name: result.full_name
                }
            };

            localStorage.setItem('auth', JSON.stringify(authData));
            localStorage.setItem('user', JSON.stringify(authData.user));

            // редирект по роли
            if (result.role === 'admin') {
                window.location.href = '/admin-dashboard.html';
            } else if (result.role === 'sitter') {
                window.location.href = '/sitter-dashboard.html';
            } else {
                window.location.href = '/dashboard.html';
            }

        } catch (err) {
            console.error('login error', err);
            if (errorBox) {
                errorBox.style.display = 'block';
                errorBox.innerText = 'Ошибка соединения с сервером';
            } else {
                alert('Ошибка соединения с сервером');
            }
        }
    });
});
