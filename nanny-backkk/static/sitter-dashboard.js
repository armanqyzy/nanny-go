// ======================
// JWT + Авторизация
// ======================

const authData = JSON.parse(localStorage.getItem('auth'));

if (!authData || !authData.token || authData.user.role !== 'sitter') {
    window.location.href = '/login';
}

const token = authData.token;
const user = authData.user;

document.getElementById('userEmail').textContent = user.email;

// Универсальная функция для API запросов с JWT
async function authFetch(url, options = {}) {
    const headers = options.headers || {};
    headers['Authorization'] = `Bearer ${token}`;
    headers['Content-Type'] = headers['Content-Type'] || 'application/json';

    const response = await fetch(url, {
        ...options,
        headers
    });

    // если токен просрочен → выходим
    if (response.status === 401) {
        alert("Сессия истекла. Войдите снова.");
        logout();
        return;
    }

    return response;
}

// ======================
// Переключение табов
// ======================

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
    switch(tab) {
        case 'overview': loadOverview(); break;
        case 'bookings': loadBookings(); break;
        case 'services': loadServices(); break;
        case 'reviews': loadReviews(); break;
        case 'profile': loadProfile(); break;
    }
}

// ======================
// ОБЗОР
// ======================

async function loadOverview() {
    try {
        const servicesRes = await authFetch(`/api/sitters/${user.id}/services`);
        const services = await servicesRes.json();
        document.getElementById('servicesCount').textContent = services.length || 0;

        const bookingsRes = await authFetch(`/api/sitters/${user.id}/bookings`);
        const bookings = await bookingsRes.json();
        document.getElementById('bookingsCount').textContent = bookings.length || 0;

        const ratingRes = await authFetch(`/api/sitters/${user.id}/rating`);
        const rating = await ratingRes.json();
        document.getElementById('ratingValue').textContent = rating.average_rating.toFixed(1);
        document.getElementById('reviewsCount').textContent = rating.review_count;

        const pending = bookings.filter(b => b.status === 'pending');
        const activeDiv = document.getElementById('activeBookings');

        if (pending.length === 0) {
            activeDiv.innerHTML = '<p class="empty-state">Нет новых заявок</p>';
        } else {
            activeDiv.innerHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Владелец</th>
                            <th>Питомец</th>
                            <th>Дата</th>
                            <th>Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${pending.map(b => `
                            <tr>
                                <td>#${b.booking_id}</td>
                                <td>ID: ${b.owner_id}</td>
                                <td>ID: ${b.pet_id}</td>
                                <td>${new Date(b.start_time).toLocaleString('ru-RU')}</td>
                                <td>
                                    <button class="btn btn-success btn-sm" onclick="confirmBooking(${b.booking_id})">Принять</button>
                                    <button class="btn btn-danger btn-sm" onclick="rejectBooking(${b.booking_id})">Отклонить</button>
                                </td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        }

        checkAccountStatus();
    } catch (err) {
        console.error('Ошибка загрузки обзора:', err);
    }
}

// ======================
// Проверка статуса няни
// ======================

async function checkAccountStatus() {
    try {
        const res = await authFetch(`/api/admin/sitters/${user.id}`);
        const details = await res.json();

        const badge = document.getElementById('statusBadge');

        if (details.status === 'pending') {
            badge.innerHTML = '<span class="badge badge-pending">⏳ На модерации</span>';
        } else if (details.status === 'approved') {
            badge.innerHTML = '<span class="badge badge-approved">✅ Одобрен</span>';
        } else if (details.status === 'rejected') {
            badge.innerHTML = '<span class="badge badge-rejected">❌ Отклонён</span>';
        }

    } catch (err) {
        console.error('Ошибка проверки статуса:', err);
    }
}

// ======================
// БРОНИРОВАНИЯ
// ======================

async function loadBookings() {
    try {
        const res = await authFetch(`/api/sitters/${user.id}/bookings`);
        const bookings = await res.json();

        // Pending
        const pending = bookings.filter(b => b.status === 'pending');
        const pendingDiv = document.getElementById('pendingBookings');

        pendingDiv.innerHTML = pending.length === 0
            ? '<p class="empty-state">Нет новых заявок</p>'
            : renderPendingBookings(pending);

        // All
        const allDiv = document.getElementById('allBookings');
        allDiv.innerHTML = bookings.length === 0
            ? '<p class="empty-state">Нет бронирований</p>'
            : renderAllBookings(bookings);

    } catch (err) {
        console.error('Ошибка загрузки бронирований:', err);
    }
}

function renderPendingBookings(list) {
    return `
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Владелец</th>
                    <th>Питомец</th>
                    <th>Начало</th>
                    <th>Конец</th>
                    <th>Действия</th>
                </tr>
            </thead>
            <tbody>
                ${list.map(b => `
                    <tr>
                        <td>#${b.booking_id}</td>
                        <td>ID: ${b.owner_id}</td>
                        <td>ID: ${b.pet_id}</td>
                        <td>${new Date(b.start_time).toLocaleString('ru-RU')}</td>
                        <td>${new Date(b.end_time).toLocaleString('ru-RU')}</td>
                        <td>
                            <button class="btn btn-success" onclick="confirmBooking(${b.booking_id})">Принять</button>
                            <button class="btn btn-danger" onclick="rejectBooking(${b.booking_id})">Отклонить</button>
                        </td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

function renderAllBookings(list) {
    return `
        <table>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Владелец</th>
                    <th>Питомец</th>
                    <th>Дата</th>
                    <th>Статус</th>
                    <th>Действия</th>
                </tr>
            </thead>
            <tbody>
                ${list.map(b => `
                    <tr>
                        <td>#${b.booking_id}</td>
                        <td>ID: ${b.owner_id}</td>
                        <td>ID: ${b.pet_id}</td>
                        <td>${new Date(b.start_time).toLocaleString('ru-RU')}</td>
                        <td><span class="badge badge-${b.status}">${b.status}</span></td>
                        <td>
                            ${b.status === 'confirmed'
        ? `<button class="btn btn-success btn-sm" onclick="completeBooking(${b.booking_id})">Завершить</button>`
        : '-'}
                        </td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

// ======================
// Услуги
// ======================

async function loadServices() {
    try {
        const res = await authFetch(`/api/sitters/${user.id}/services`);
        const services = await res.json();

        const div = document.getElementById('servicesList');

        if (services.length === 0) {
            div.innerHTML = '<div class="empty-state"><h3>У вас пока нет услуг</h3></div>';
            return;
        }

        div.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Тип</th>
                        <th>Цена</th>
                        <th>Описание</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${services.map(s => `
                        <tr>
                            <td>${getServiceTypeName(s.type)}</td>
                            <td>${s.price_per_hour} ₸</td>
                            <td>${s.description || '-'}</td>
                            <td>
                                <button class="btn btn-danger btn-sm" onclick="deleteService(${s.service_id})">Удалить</button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (err) {
        console.error('Ошибка загрузки услуг:', err);
    }
}

// Добавление услуги
document.getElementById('addServiceForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const data = Object.fromEntries(new FormData(e.target));
    data.sitter_id = user.id;
    data.price_per_hour = Number(data.price_per_hour);

    try {
        const res = await authFetch('/api/services', {
            method: 'POST',
            body: JSON.stringify(data)
        });

        if (res.ok) {
            alert('Услуга добавлена!');
            closeModal('addServiceModal');
            e.target.reset();
            loadServices();
            loadOverview();
        } else {
            const err = await res.json();
            alert('Ошибка: ' + err.error);
        }
    } catch {
        alert('Ошибка соединения');
    }
});

// ======================
// Отзывы
// ======================

async function loadReviews() {
    try {
        const [reviewsRes, ratingRes] = await Promise.all([
            authFetch(`/api/sitters/${user.id}/reviews`),
            authFetch(`/api/sitters/${user.id}/rating`)
        ]);

        const reviews = await reviewsRes.json();
        const rating = await ratingRes.json();

        document.getElementById('avgRating').textContent = rating.average_rating.toFixed(1);
        document.getElementById('totalReviews').textContent = rating.review_count;

        const div = document.getElementById('reviewsList');

        if (reviews.length === 0) {
            div.innerHTML = '<p class="empty-state">Нет отзывов</p>';
            return;
        }

        div.innerHTML = reviews.map(r => `
            <div class="card">
                <div class="rating">${renderStars(r.rating)}</div>
                <p><strong>Владелец #${r.owner_id}</strong></p>
                <p>${r.comment}</p>
                <p class="date">${new Date(r.created_at).toLocaleDateString('ru-RU')}</p>
            </div>
        `).join('');

    } catch (err) {
        console.error('Ошибка загрузки отзывов:', err);
    }
}

// ======================
// Профиль няни
// ======================

async function loadProfile() {
    try {
        const res = await authFetch(`/api/admin/sitters/${user.id}`);
        const d = await res.json();

        document.getElementById('profileInfo').innerHTML = `
            <p><strong>Имя:</strong> ${d.full_name}</p>
            <p><strong>Email:</strong> ${d.email}</p>
            <p><strong>Телефон:</strong> ${d.phone}</p>
            <p><strong>Опыт:</strong> ${d.experience_years} лет</p>
            <p><strong>Сертификаты:</strong> ${d.certificates || '-'}</p>
            <p><strong>Предпочтения:</strong> ${d.preferences || '-'}</p>
            <p><strong>Локация:</strong> ${d.location || '-'}</p>
            <p><strong>Статус:</strong> <span class="badge badge-${d.status}">${d.status}</span></p>
        `;
    } catch (err) {
        console.error('Ошибка загрузки профиля:', err);
    }
}

// ======================
// Действия бронирований
// ======================

async function confirmBooking(id) {
    await authFetch(`/api/bookings/${id}/confirm`, { method: 'POST' });
    loadBookings();
    loadOverview();
}

async function rejectBooking(id) {
    if (!confirm('Отклонить?')) return;
    await authFetch(`/api/bookings/${id}/cancel`, { method: 'POST' });
    loadBookings();
    loadOverview();
}

async function completeBooking(id) {
    if (!confirm('Завершить?')) return;
    await authFetch(`/api/bookings/${id}/complete`, { method: 'POST' });
    loadBookings();
}

// ======================
// Удалить услугу
// ======================

async function deleteService(id) {
    if (!confirm('Удалить услугу?')) return;
    await authFetch(`/api/services/${id}`, { method: 'DELETE' });
    loadServices();
    loadOverview();
}

// ======================
// Utils
// ======================

function getServiceTypeName(type) {
    return {
        walking: 'Выгул',
        boarding: 'Передержка',
        'home-care': 'Уход на дому'
    }[type] || type;
}

function renderStars(r) {
    return '⭐'.repeat(Math.floor(r)) + '☆'.repeat(5 - Math.floor(r));
}

// ======================
// Logout
// ======================

function logout() {
    localStorage.removeItem('auth');
    window.location.href = '/login';
}

// Загружаем dashboard
loadOverview();
