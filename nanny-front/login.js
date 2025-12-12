document.addEventListener('DOMContentLoaded', () => {
    const BASE_URL = "http://localhost:8080";

    const form = document.getElementById('loginForm');
    const errorBox = document.getElementById('error-box');

    if (!form) return;

    if (errorBox) {
        errorBox.style.display = 'none';
        errorBox.innerText = '';
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (errorBox) {
            errorBox.style.display = 'none';
            errorBox.innerText = '';
        }

        const formData = new FormData(form);
        const data = Object.fromEntries(formData);

        try {
            const res = await fetch(`${BASE_URL}/api/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            const result = await res.json();

            if (!res.ok) {
                if (errorBox) {
                    errorBox.style.display = 'block';
                    errorBox.innerText = result.error || 'Ошибка входа';
                }
                return;
            }

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

            if (result.role === 'admin') window.location.href = 'admin-dashboard.html';
            else if (result.role === 'sitter') window.location.href = 'sitter-dashboard.html';
            else window.location.href = 'dashboard.html';


        } catch (err) {
            if (errorBox) {
                errorBox.style.display = 'block';
                errorBox.innerText = 'Ошибка соединения с сервером';
            }
        }
    });
});
