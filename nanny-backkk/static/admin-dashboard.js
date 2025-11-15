// Проверка авторизации
const user = JSON.parse(localStorage.getItem('user'));
if (!user || user.role !== 'admin') {
    window.location.href = '/login';
}

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
        case 'pending':
            loadPendingSitters();
            break;
        case 'users':
            loadUsers();
            break;
        case 'sitters':
            loadSitters();
            break;
    }
}

// Загрузка обзора
async function loadOverview() {
    try {
        // Загружаем всех пользователей
        const usersRes = await fetch('/api/admin/users');
        const users = await usersRes.json();
        document.getElementById('totalUsers').textContent = users.length;
        
        // Фильтруем нянь
        const sitters = users.filter(u => u.role === 'sitter');
        document.getElementById('totalSitters').textContent = sitters.length;
        
        // Загружаем заявки нянь
        const pendingRes = await fetch('/api/admin/sitters/pending');
        const pending = await pendingRes.json();
        document.getElementById('pendingCount').textContent = pending.length;
        
        // Считаем одобренных
        // Для этого нам нужно получить всех нянь с полной информацией
        // В реальности это можно оптимизировать
        document.getElementById('approvedCount').textContent = sitters.length - pending.length;
        
        // Последние 5 пользователей
        const recentDiv = document.getElementById('recentUsers');
        const recent = users.slice(0, 5);
        
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
                            <td><span class="badge badge-${u.role === 'admin' ? 'approved' : 'pending'}">${u.role}</span></td>
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

// Загрузка заявок нянь
async function loadPendingSitters() {
    try {
        const res = await fetch('/api/admin/sitters/pending');
        const sitters = await res.json();
        
        const pendingDiv = document.getElementById('pendingSitters');
        
        if (sitters.length === 0) {
            pendingDiv.innerHTML = '<div class="empty-state"><h3>Нет заявок на модерацию</h3></div>';
            return;
        }
        
        // Загружаем полную информацию для каждой няни
        const sittersWithDetails = await Promise.all(
            sitters.map(async (s) => {
                try {
                    const detailsRes = await fetch(`/api/admin/sitters/${s.sitter_id}`);
                    return await detailsRes.json();
                } catch {
                    return s;
                }
            })
        );
        
        pendingDiv.innerHTML = sittersWithDetails.map(s => `
            <div class="card" style="margin-bottom: 20px;">
                <h3>${s.full_name || 'Няня #' + s.sitter_id}</h3>
                <p><strong>Email:</strong> ${s.email || '-'}</p>
                <p><strong>Телефон:</strong> ${s.phone || '-'}</p>
                <p><strong>Опыт:</strong> ${s.experience_years || 0} лет</p>
                <p><strong>Сертификаты:</strong> ${s.certificates || '-'}</p>
                <p><strong>Предпочтения:</strong> ${s.preferences || '-'}</p>
                <p><strong>Локация:</strong> ${s.location || '-'}</p>
                <div style="margin-top: 15px;">
                    <button class="btn btn-success" onclick="approveSitter(${s.sitter_id})">✅ Одобрить</button>
                    <button class="btn btn-danger" onclick="rejectSitter(${s.sitter_id})">❌ Отклонить</button>
                    <button class="btn btn-secondary" onclick="showSitterDetails(${s.sitter_id})">👁️ Подробнее</button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        console.error('Ошибка загрузки заявок:', err);
    }
}

// Загрузка всех пользователей
async function loadUsers() {
    try {
        const res = await fetch('/api/admin/users');
        const users = await res.json();
        
        const usersDiv = document.getElementById('usersList');
        
        usersDiv.innerHTML = `
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
                            <td><span class="badge badge-${u.role === 'admin' ? 'approved' : 'pending'}">${u.role}</span></td>
                            <td>${new Date(u.created_at).toLocaleDateString('ru-RU')}</td>
                            <td>
                                ${u.role !== 'admin' ? `<button class="btn btn-danger btn-sm" onclick="deleteUser(${u.user_id}, '${u.full_name}')">Удалить</button>` : '-'}
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

// Загрузка всех нянь
async function loadSitters() {
    try {
        const usersRes = await fetch('/api/admin/users');
        const users = await usersRes.json();
        
        const sitters = users.filter(u => u.role === 'sitter');
        
        const sittersDiv = document.getElementById('sittersList');
        
        if (sitters.length === 0) {
            sittersDiv.innerHTML = '<p class="empty-state">Нет зарегистрированных нянь</p>';
            return;
        }
        
        // Загружаем детали для каждой няни
        const sittersWithDetails = await Promise.all(
            sitters.map(async (s) => {
                try {
                    const detailsRes = await fetch(`/api/admin/sitters/${s.user_id}`);
                    const details = await detailsRes.json();
                    
                    const ratingRes = await fetch(`/api/sitters/${s.user_id}/rating`);
                    const rating = await ratingRes.json();
                    
                    return { ...details, rating: rating.average_rating, review_count: rating.review_count };
                } catch {
                    return s;
                }
            })
        );
        
        sittersDiv.innerHTML = `
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
                            <td>#${s.sitter_id || s.user_id}</td>
                            <td>${s.full_name}</td>
                            <td>${s.location || '-'}</td>
                            <td>${s.experience_years || 0} лет</td>
                            <td>${renderStars(s.rating || 0)} (${(s.rating || 0).toFixed(1)})</td>
                            <td><span class="badge badge-${s.status}">${s.status}</span></td>
                            <td>
                                <button class="btn btn-secondary btn-sm" onclick="showSitterDetails(${s.sitter_id || s.user_id})">Подробнее</button>
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

// Показать детали няни
async function showSitterDetails(sitterId) {
    try {
        const [detailsRes, reviewsRes, servicesRes] = await Promise.all([
            fetch(`/api/admin/sitters/${sitterId}`),
            fetch(`/api/sitters/${sitterId}/reviews`),
            fetch(`/api/sitters/${sitterId}/services`)
        ]);
        
        const details = await detailsRes.json();
        const reviews = await reviewsRes.json();
        const services = await servicesRes.json();
        
        const content = document.getElementById('sitterDetailsContent');
        content.innerHTML = `
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
                <p>${renderStars(details.rating)} (${details.rating.toFixed(1)}) - ${details.reviews} отзывов</p>
            </div>
            
            <h3 style="margin-top: 30px;">Услуги (${services.length})</h3>
            ${services.length > 0 ? `
                <ul>
                    ${services.map(s => `<li>${getServiceTypeName(s.type)} - ${s.price_per_hour} ₸/час</li>`).join('')}
                </ul>
            ` : '<p>Нет услуг</p>'}
            
            <h3 style="margin-top: 30px;">Отзывы (${reviews.length})</h3>
            ${reviews.length > 0 ? reviews.slice(0, 3).map(r => `
                <div style="margin-bottom: 15px; padding: 10px; background: #f5f5f5; border-radius: 8px;">
                    <div class="rating">${renderStars(r.rating)}</div>
                    <p style="margin-top: 5px;">${r.comment}</p>
                </div>
            `).join('') : '<p>Нет отзывов</p>'}
        `;
        
        document.getElementById('sitterDetailsModal').classList.add('active');
    } catch (err) {
        console.error('Ошибка загрузки деталей:', err);
        alert('Ошибка загрузки деталей няни');
    }
}

// Одобрить няню
async function approveSitter(sitterId) {
    if (!confirm('Одобрить эту няню?')) return;
    
    try {
        const res = await fetch(`/api/admin/sitters/${sitterId}/approve`, { method: 'POST' });
        if (res.ok) {
            alert('✅ Няня одобрена!');
            loadPendingSitters();
            loadOverview();
        } else {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка одобрения');
    }
}

// Отклонить няню
async function rejectSitter(sitterId) {
    if (!confirm('Отклонить эту няню?')) return;
    
    try {
        const res = await fetch(`/api/admin/sitters/${sitterId}/reject`, { method: 'POST' });
        if (res.ok) {
            alert('✅ Няня отклонена');
            loadPendingSitters();
            loadOverview();
        } else {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка отклонения');
    }
}

// Удалить пользователя
async function deleteUser(userId, userName) {
    if (!confirm(`Удалить пользователя "${userName}"?`)) return;
    
    try {
        const res = await fetch(`/api/admin/users/${userId}`, { method: 'DELETE' });
        if (res.ok) {
            alert('✅ Пользователь удалён');
            loadUsers();
            loadOverview();
        } else {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка удаления');
    }
}

// Модальные окна
function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
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
