const API_URL = 'http://localhost:8080';

let authData = null;

try {
    const raw = localStorage.getItem('auth');
    authData = raw ? JSON.parse(raw) : null;
} catch (_) {
    authData = null;
}

if (!authData || !authData.token || !authData.user || authData.user.role !== 'admin') {
    window.location.href = 'login.html';
}

const token = authData.token;
const user = authData.user;

async function authFetch(url, options = {}) {
    const headers = options.headers || {};
    headers['Authorization'] = `Bearer ${token}`;
    if (!headers['Content-Type'] && options.method && options.method !== 'GET') {
        headers['Content-Type'] = 'application/json';
    }

    const fullUrl = url.startsWith('http') ? url : `${API_URL}${url}`;

    const res = await fetch(fullUrl, { ...options, headers });

    if (res.status === 401) {
        alert('Сессия истекла. Войдите снова.');
        logout();
        return Promise.reject(new Error('Unauthorized'));
    }

    return res;
}

document.querySelectorAll('.sidebar-menu a').forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        const tab = e.target.dataset.tab;

        document.querySelectorAll('.sidebar-menu a').forEach(l => l.classList.remove('active'));
        e.target.classList.add('active');

        document.querySelectorAll('.tab-content').forEach(t => t.style.display = 'none');
        document.getElementById(tab + '-tab').style.display = 'block';

        loadTabData(tab);
    });
});

function loadTabData(tab) {
    switch (tab) {
        case 'overview': loadOverview(); break;
        case 'pending': loadPendingSitters(); break;
        case 'users':   loadUsers(); break;
        case 'sitters': loadSitters(); break;
    }
}

async function loadOverview() {
    try {
        const usersRes = await authFetch('/api/admin/users');
        const users = await usersRes.json();
        document.getElementById('totalUsers').textContent = users.length;

        const sitters = users.filter(u => u.role === 'sitter');
        document.getElementById('totalSitters').textContent = sitters.length;

        const pendingRes = await authFetch('/api/admin/sitters/pending');
        const pending = await pendingRes.json();
        document.getElementById('pendingCount').textContent = pending.length;

        document.getElementById('approvedCount').textContent = sitters.length - pending.length;

        const recent = users.slice(0, 5);
        const recentDiv = document.getElementById('recentUsers');

        recentDiv.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Имя</th>
                        <th>Email</th>
                        <th>Роль</th>
                        <th>Дата регистрации</th>
                    </tr>
                </thead>
                <tbody>
                    ${recent.map(u => `
                        <tr>
                            <td>#${u.user_id}</td>
                            <td>${u.full_name}</td>
                            <td>${u.email}</td>
                            <td><span class="badge badge-${u.role}">${u.role}</span></td>
                            <td>${new Date(u.created_at).toLocaleDateString('ru-RU')}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (err) {
        console.error('Ошибка загрузки обзора:', err);
    }
}

async function loadPendingSitters() {
    try {
        const res = await authFetch('/api/admin/sitters/pending');
        const sitters = await res.json();

        const div = document.getElementById('pendingSitters');

        if (sitters.length === 0) {
            div.innerHTML = '<div class="empty-state"><h3>Нет заявок на модерацию</h3></div>';
            return;
        }

        const sittersWithDetails = await Promise.all(
            sitters.map(async (s) => {
                const dRes = await authFetch(`/api/admin/sitters/${s.sitter_id}`);
                return await dRes.json();
            })
        );

        div.innerHTML = sittersWithDetails.map(s => `
            <div class="card" style="margin-bottom: 20px;">
                <h3>${s.full_name}</h3>
                <p><strong>Email:</strong> ${s.email}</p>
                <p><strong>Телефон:</strong> ${s.phone}</p>
                <p><strong>Опыт:</strong> ${s.experience_years} лет</p>
                <p><strong>Локация:</strong> ${s.location}</p>

                <div class="actions" style="margin-top: 10px;">
                    <button class="btn btn-success" onclick="approveSitter(${s.sitter_id})">Одобрить</button>
                    <button class="btn btn-danger" onclick="rejectSitter(${s.sitter_id})">Отклонить</button>
                    <button class="btn btn-secondary" onclick="showSitterDetails(${s.sitter_id})">Подробнее</button>
                </div>
            </div>
        `).join('');

    } catch (e) {
        console.error('Ошибка загрузки pending sitters', e);
    }
}

async function loadUsers() {
    try {
        const res = await authFetch('/api/admin/users');
        const users = await res.json();

        const div = document.getElementById('usersList');

        div.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Имя</th>
                        <th>Email</th>
                        <th>Телефон</th>
                        <th>Роль</th>
                        <th>Дата регистрации</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${users.map(u => `
                        <tr>
                            <td>#${u.user_id}</td>
                            <td>${u.full_name}</td>
                            <td>${u.email}</td>
                            <td>${u.phone}</td>
                            <td>${u.role}</td>
                            <td>${new Date(u.created_at).toLocaleDateString('ru-RU')}</td>
                            <td>
                                ${u.role !== 'admin'
            ? `<button class="btn btn-danger btn-sm" onclick="deleteUser(${u.user_id}, '${u.full_name}')">Удалить</button>`
            : '-'}
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (err) {
        console.error('Ошибка загрузки пользователей:', err);
    }
}

async function loadSitters() {
    try {
        const usersRes = await authFetch('/api/admin/users');
        const users = await usersRes.json();

        const sitters = users.filter(u => u.role === 'sitter');

        const div = document.getElementById('sittersList');

        if (sitters.length === 0) {
            div.innerHTML = '<p class="empty-state">Нет нянь</p>';
            return;
        }

        const sittersWithDetails = await Promise.all(
            sitters.map(async (s) => {
                const detailsRes = await authFetch(`/api/admin/sitters/${s.user_id}`);
                const details = await detailsRes.json();

                const ratingRes = await authFetch(`/api/sitters/${s.user_id}/rating`);
                const rating = await ratingRes.json();

                return { ...details, rating: rating.average_rating };
            })
        );

        div.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Имя</th>
                        <th>Локация</th>
                        <th>Опыт</th>
                        <th>Рейтинг</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${sittersWithDetails.map(s => `
                        <tr>
                            <td>#${s.sitter_id}</td>
                            <td>${s.full_name}</td>
                            <td>${s.location || '-'}</td>
                            <td>${s.experience_years} лет</td>
                            <td>${renderStars(s.rating || 0)}</td>
                            <td><span class="badge badge-${s.status}">${s.status}</span></td>
                            <td>
                                <button onclick="showSitterDetails(${s.sitter_id})" class="btn btn-secondary btn-sm">Подробнее</button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

    } catch (err) {
        console.error('Ошибка загрузки нянь:', err);
    }
}

async function showSitterDetails(id) {
    try {
        const dRes = await authFetch(`/api/admin/sitters/${id}`);
        const details = await dRes.json();

        const reviewsRes = await authFetch(`/api/sitters/${id}/reviews`);
        const reviews = await reviewsRes.json();

        const servicesRes = await authFetch(`/api/sitters/${id}/services`);
        const services = await servicesRes.json();

        const content = document.getElementById('sitterDetailsContent');

        content.innerHTML = `
            <p><strong>Имя:</strong> ${details.full_name}</p>
            <p><strong>Email:</strong> ${details.email}</p>
            <p><strong>Телефон:</strong> ${details.phone}</p>
            <p><strong>Опыт:</strong> ${details.experience_years} лет</p>
            <p><strong>Статус:</strong> ${details.status}</p>

            <h3>Услуги</h3>
            ${services.length
            ? services.map(s => `<p>${getServiceTypeName(s.type)} — ${s.price_per_hour} ₸</p>`).join('')
            : 'Нет услуг'}

            <h3>Отзывы</h3>
            ${reviews.length
            ? reviews.slice(0, 3).map(r => `
                    <div class="card" style="margin-bottom: 10px;">
                        <div class="rating">${renderStars(r.rating)}</div>
                        <p>${r.comment}</p>
                    </div>
                `).join('')
            : 'Нет отзывов'}
        `;

        document.getElementById('sitterDetailsModal').classList.add('active');

    } catch (err) {
        console.error('Ошибка загрузки деталей', err);
    }
}

async function approveSitter(id) {
    if (!confirm('Одобрить няню?')) return;
    try {
        await authFetch(`/api/admin/sitters/${id}/approve`, { method: 'POST' });
        loadPendingSitters();
        loadOverview();
    } catch (err) {
        console.error('Ошибка одобрения няни', err);
    }
}

async function rejectSitter(id) {
    if (!confirm('Отклонить няню?')) return;
    try {
        await authFetch(`/api/admin/sitters/${id}/reject`, { method: 'POST' });
        loadPendingSitters();
        loadOverview();
    } catch (err) {
        console.error('Ошибка отклонения няни', err);
    }
}

async function deleteUser(id, name) {
    if (!confirm(`Удалить пользователя ${name}?`)) return;
    try {
        await authFetch(`/api/admin/users/${id}`, { method: 'DELETE' });
        loadUsers();
        loadOverview();
    } catch (err) {
        console.error('Ошибка удаления пользователя', err);
    }
}

function getServiceTypeName(t) {
    return {
        walking: 'Выгул',
        boarding: 'Передержка',
        'home-care': 'Уход на дому'
    }[t] || t;
}

function renderStars(r) {
    const f = Math.floor(r);
    return '⭐'.repeat(f) + '☆'.repeat(5 - f);
}

function closeModal(id) {
    document.getElementById(id).classList.remove('active');
}

function logout() {
    localStorage.removeItem('auth');
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    window.location.href = 'login.html';
}

loadOverview();
