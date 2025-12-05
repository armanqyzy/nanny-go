// Проверка авторизации
const user = JSON.parse(localStorage.getItem('user'));
if (!user || user.role !== 'sitter') {
    window.location.href = '/login';
}

document.getElementById('userEmail').textContent = user.email;

// Переключение табов
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
        case 'overview':
            loadOverview();
            break;
        case 'bookings':
            loadBookings();
            break;
        case 'services':
            loadServices();
            break;
        case 'reviews':
            loadReviews();
            break;
        case 'profile':
            loadProfile();
            break;
    }
}

// Загрузка обзора
async function loadOverview() {
    try {
        // Загружаем услуги
        const servicesRes = await fetch(`/api/sitters/${user.id}/services`);
        const services = await servicesRes.json();
        document.getElementById('servicesCount').textContent = services.length || 0;

        // Загружаем бронирования
        const bookingsRes = await fetch(`/api/sitters/${user.id}/bookings`);
        const bookings = await bookingsRes.json();
        document.getElementById('bookingsCount').textContent = bookings.length || 0;

        // Загружаем рейтинг
        const ratingRes = await fetch(`/api/sitters/${user.id}/rating`);
        const rating = await ratingRes.json();
        document.getElementById('ratingValue').textContent = rating.average_rating.toFixed(1);
        document.getElementById('reviewsCount').textContent = rating.review_count;

        // Активные заявки (pending)
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

        // Проверяем статус аккаунта
        checkAccountStatus();
    } catch (err) {
        console.error('Ошибка загрузки обзора:', err);
    }
}

// Проверка статуса аккаунта
async function checkAccountStatus() {
    try {
        const res = await fetch(`/api/admin/sitters/${user.id}`);
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

// Загрузка бронирований
async function loadBookings() {
    try {
        const res = await fetch(`/api/sitters/${user.id}/bookings`);
        const bookings = await res.json();

        // Новые заявки
        const pending = bookings.filter(b => b.status === 'pending');
        const pendingDiv = document.getElementById('pendingBookings');

        if (pending.length === 0) {
            pendingDiv.innerHTML = '<p class="empty-state">Нет новых заявок</p>';
        } else {
            pendingDiv.innerHTML = `
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
                        ${pending.map(b => `
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

        // Вся история
        const allDiv = document.getElementById('allBookings');
        if (bookings.length === 0) {
            allDiv.innerHTML = '<p class="empty-state">Нет бронирований</p>';
        } else {
            allDiv.innerHTML = `
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
                        ${bookings.map(b => `
                            <tr>
                                <td>#${b.booking_id}</td>
                                <td>ID: ${b.owner_id}</td>
                                <td>ID: ${b.pet_id}</td>
                                <td>${new Date(b.start_time).toLocaleString('ru-RU')}</td>
                                <td><span class="badge badge-${b.status}">${b.status}</span></td>
                                <td>
                                    ${b.status === 'confirmed' ? `<button class="btn btn-success btn-sm" onclick="completeBooking(${b.booking_id})">Завершить</button>` : '-'}
                                </td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        }
    } catch (err) {
        console.error('Ошибка загрузки бронирований:', err);
    }
}

// Загрузка услуг
async function loadServices() {
    try {
        const res = await fetch(`/api/sitters/${user.id}/services`);
        const services = await res.json();

        const servicesDiv = document.getElementById('servicesList');

        if (services.length === 0) {
            servicesDiv.innerHTML = '<div class="empty-state"><h3>У вас пока нет услуг</h3><p>Добавьте первую услугу!</p></div>';
            return;
        }

        servicesDiv.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Тип</th>
                        <th>Цена (₸/час)</th>
                        <th>Описание</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${services.map(s => `
                        <tr>
                            <td><strong>${getServiceTypeName(s.type)}</strong></td>
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

// Загрузка отзывов
async function loadReviews() {
    try {
        const [reviewsRes, ratingRes] = await Promise.all([
            fetch(`/api/sitters/${user.id}/reviews`),
            fetch(`/api/sitters/${user.id}/rating`)
        ]);

        const reviews = await reviewsRes.json();
        const rating = await ratingRes.json();

        document.getElementById('avgRating').textContent = rating.average_rating.toFixed(1);
        document.getElementById('totalReviews').textContent = rating.review_count;

        const reviewsDiv = document.getElementById('reviewsList');

        if (reviews.length === 0) {
            reviewsDiv.innerHTML = '<p class="empty-state">Пока нет отзывов</p>';
            return;
        }

        reviewsDiv.innerHTML = reviews.map(r => `
            <div class="card" style="margin-bottom: 15px;">
                <div class="rating">${renderStars(r.rating)}</div>
                <p style="margin-top: 10px;"><strong>Владелец #${r.owner_id}</strong></p>
                <p>${r.comment}</p>
                <p style="color: #999; font-size: 12px;">${new Date(r.created_at).toLocaleDateString('ru-RU')}</p>
            </div>
        `).join('');
    } catch (err) {
        console.error('Ошибка загрузки отзывов:', err);
    }
}

// Загрузка профиля
async function loadProfile() {
    try {
        const res = await fetch(`/api/admin/sitters/${user.id}`);
        const details = await res.json();

        document.getElementById('profileInfo').innerHTML = `
            <div class="form-group">
                <label>Имя:</label>
                <p>${details.full_name}</p>
            </div>
            <div class="form-group">
                <label>Email:</label>
                <p>${details.email}</p>
            </div>
            <div class="form-group">
                <label>Телефон:</label>
                <p>${details.phone}</p>
            </div>
            <div class="form-group">
                <label>Опыт работы:</label>
                <p>${details.experience_years} лет</p>
            </div>
            <div class="form-group">
                <label>Сертификаты:</label>
                <p>${details.certificates || '-'}</p>
            </div>
            <div class="form-group">
                <label>Предпочтения:</label>
                <p>${details.preferences || '-'}</p>
            </div>
            <div class="form-group">
                <label>Локация:</label>
                <p>${details.location || '-'}</p>
            </div>
            <div class="form-group">
                <label>Статус:</label>
                <p><span class="badge badge-${details.status}">${details.status}</span></p>
            </div>
            <div class="form-group">
                <label>Рейтинг:</label>
                <p>${renderStars(details.rating)} (${details.rating.toFixed(1)})</p>
            </div>
        `;
    } catch (err) {
        console.error('Ошибка загрузки профиля:', err);
    }
}

// Модальные окна
function showAddServiceModal() {
    document.getElementById('addServiceModal').classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

// Обработчик формы добавления услуги
document.getElementById('addServiceForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData);
    data.sitter_id = user.id;
    data.price_per_hour = parseFloat(data.price_per_hour);

    try {
        const res = await fetch('/api/services', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });

        if (res.ok) {
            alert('✅ Услуга добавлена!');
            closeModal('addServiceModal');
            e.target.reset();
            loadServices();
            loadOverview();
        } else {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка соединения');
    }
});

// Действия с бронированиями
async function confirmBooking(bookingId) {
    try {
        const res = await fetch(`/api/bookings/${bookingId}/confirm`, { method: 'POST' });
        if (res.ok) {
            alert('✅ Бронирование подтверждено');
            loadBookings();
            loadOverview();
        }
    } catch (err) {
        alert('Ошибка подтверждения');
    }
}

async function rejectBooking(bookingId) {
    if (!confirm('Отклонить заявку?')) return;

    try {
        const res = await fetch(`/api/bookings/${bookingId}/cancel`, { method: 'POST' });
        if (res.ok) {
            alert('✅ Заявка отклонена');
            loadBookings();
            loadOverview();
        }
    } catch (err) {
        alert('Ошибка отклонения');
    }
}

async function completeBooking(bookingId) {
    if (!confirm('Завершить бронирование?')) return;

    try {
        const res = await fetch(`/api/bookings/${bookingId}/complete`, { method: 'POST' });
        if (res.ok) {
            alert('✅ Бронирование завершено');
            loadBookings();
        }
    } catch (err) {
        alert('Ошибка завершения');
    }
}

// Удалить услугу
async function deleteService(serviceId) {
    if (!confirm('Удалить услугу?')) return;

    try {
        const res = await fetch(`/api/services/${serviceId}`, { method: 'DELETE' });
        if (res.ok) {
            alert('✅ Услуга удалена');
            loadServices();
            loadOverview();
        }
    } catch (err) {
        alert('Ошибка удаления');
    }
}

// Вспомогательные функции
function getServiceTypeName(type) {
    const names = {
        walking: 'Выгул',
        boarding: 'Передержка',
        'home-care': 'Уход на дому'
    };
    return names[type] || type;
}

function renderStars(rating) {
    const full = Math.floor(rating);
    const empty = 5 - full;
    return '⭐'.repeat(full) + '☆'.repeat(empty);
}

function logout() {
    localStorage.removeItem('user');
    window.location.href = '/login';
}

// Загружаем обзор при старте
loadOverview();
