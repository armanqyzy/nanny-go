const API_URL = 'http://localhost:8080';

let authData = JSON.parse(localStorage.getItem('auth') || 'null');

if (!authData && localStorage.getItem('token') && localStorage.getItem('user')) {
    authData = {
        token: localStorage.getItem('token'),
        user: JSON.parse(localStorage.getItem('user'))
    };
}

if (!authData || !authData.token || !authData.user || authData.user.role !== 'owner') {
    window.location.href = 'login.html';
}

const token = authData.token;
const user = authData.user;

async function authFetch(url, options = {}) {
    const headers = options.headers || {};
    headers['Authorization'] = `Bearer ${token}`;
    if (!headers['Content-Type'] && options.body && !(options.body instanceof FormData)) {
        headers['Content-Type'] = 'application/json';
    }

    const fullUrl = url.startsWith('http') ? url : `${API_URL}${url}`;

    const res = await fetch(fullUrl, { ...options, headers });

    if (res.status === 401) {
        alert('Сессия истекла. Войдите снова.');
        logout();
        return;
    }

    return res;
}

document.getElementById('userEmail').textContent = user.email;

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
        case 'pets':
            loadPets();
            break;
        case 'bookings':
            loadBookings();
            break;
        case 'reviews':
            loadReviews();
            break;
    }
}

async function loadOverview() {
    try {
        const petsRes = await authFetch(`/api/owners/${user.id}/pets`);
        if (!petsRes) return;
        const pets = await petsRes.json();
        document.getElementById('petsCount').textContent = pets.length || 0;

        const bookingsRes = await authFetch(`/api/owners/${user.id}/bookings`);
        if (!bookingsRes) return;
        const bookings = await bookingsRes.json();
        document.getElementById('bookingsCount').textContent = bookings.length || 0;

        const recentDiv = document.getElementById('recentBookings');
        if (!bookings.length) {
            recentDiv.innerHTML = '<p class="empty-state">Пока нет бронирований</p>';
        } else {
            const recent = bookings.slice(0, 5);
            recentDiv.innerHTML = `
                <table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Питомец</th>
                            <th>Дата</th>
                            <th>Статус</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${recent.map(b => `
                            <tr>
                                <td>#${b.booking_id}</td>
                                <td>Питомец #${b.pet_id}</td>
                                <td>${new Date(b.start_time).toLocaleDateString('ru-RU')}</td>
                                <td><span class="badge badge-${b.status}">${b.status}</span></td>
                            </tr>
                        `).join('')}
                    </tbody>
                </table>
            `;
        }
    } catch (err) {
        console.error('Ошибка загрузки обзора:', err);
    }
}

async function loadPets() {
    const petsDiv = document.getElementById('petsList');

    try {
        const res = await authFetch(`/api/owners/${user.id}/pets`);
        if (!res) return;

        const pets = await res.json();

        if (!Array.isArray(pets) || pets.length === 0) {
            petsDiv.innerHTML = '<div class="empty-state"><h3>У вас пока нет питомцев</h3><p>Добавьте первого питомца!</p></div>';
            return;
        }

        petsDiv.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Имя</th>
                        <th>Тип</th>
                        <th>Возраст</th>
                        <th>Особенности</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${pets.map(pet => `
                        <tr>
                            <td><strong>${pet.name}</strong></td>
                            <td>${getPetTypeIcon(pet.type)} ${pet.type}</td>
                            <td>${pet.age} ${getPetAgeWord(pet.age)}</td>
                            <td>${pet.notes || '-'}</td>
                            <td>
                                <button class="btn btn-danger btn-sm" onclick="deletePet(${pet.pet_id})">Удалить</button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (err) {
        console.error('Ошибка загрузки питомцев:', err);
        petsDiv.innerHTML = '<div class="empty-state"><h3>Ошибка загрузки питомцев</h3></div>';
    }
}

async function loadBookings() {
    try {
        const res = await authFetch(`/api/owners/${user.id}/bookings`);
        if (!res) return;
        const bookings = await res.json();

        const bookingsDiv = document.getElementById('bookingsList');

        if (!bookings.length) {
            bookingsDiv.innerHTML = '<div class="empty-state"><h3>У вас пока нет бронирований</h3></div>';
            return;
        }

        bookingsDiv.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Питомец</th>
                        <th>Няня</th>
                        <th>Дата</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${bookings.map(b => `
                        <tr>
                            <td>#${b.booking_id}</td>
                            <td>Питомец #${b.pet_id}</td>
                            <td>Няня #${b.sitter_id}</td>
                            <td>${new Date(b.start_time).toLocaleString('ru-RU')}</td>
                            <td><span class="badge badge-${b.status}">${b.status}</span></td>
                            <td>
                                ${b.status === 'pending' ? `<button class="btn btn-danger btn-sm" onclick="cancelBooking(${b.booking_id})">Отменить</button>` : ''}
                                ${b.status === 'completed' ? `<button class="btn btn-primary btn-sm" onclick="showReviewModal(${b.booking_id}, ${b.sitter_id})">Оставить отзыв</button>` : ''}
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    } catch (err) {
        console.error('Ошибка загрузки бронирований:', err);
    }
}

async function searchServices() {
    const typeEl =
        document.getElementById('serviceTypeFilter') ||
        document.getElementById('serviceType');

    const locationEl =
        document.getElementById('locationFilter') ||
        document.getElementById('location');

    const type = typeEl ? typeEl.value : '';
    const location = locationEl ? locationEl.value.trim() : '';

    try {
        const params = new URLSearchParams();
        if (type && type !== 'all') params.append('type', type);
        if (location) params.append('location', location);

        const res = await authFetch(`/api/services/search?${params.toString()}`);
        if (!res) return;
        const services = await res.json();

        console.log('services search result:', services);

        const resultsDiv = document.getElementById('searchResults');

        if (!Array.isArray(services) || services.length === 0) {
            resultsDiv.innerHTML = '<p class="empty-state">Услуги не найдены</p>';
            return;
        }

        resultsDiv.innerHTML = services.map(s => `
            <div class="card" style="margin-bottom: 15px;">
                <h3>${s.sitter_name}</h3>
                <div class="rating">
                    ${renderStars(s.sitter_rating || 0)}
                    <span>(${(s.sitter_rating || 0).toFixed(1)})</span>
                </div>
                <p><strong>Услуга:</strong> ${getServiceTypeName(s.type)}</p>
                <p><strong>Цена:</strong> ${s.price_per_hour} ₸/час</p>
                <p>${s.description || ''}</p>
                <button class="btn btn-primary"
                        onclick="bookService(${s.sitter_id}, ${s.service_id})">
                    Забронировать
                </button>
            </div>
        `).join('');
    } catch (err) {
        console.error('Ошибка поиска услуг:', err);
        document.getElementById('searchResults').innerHTML =
            '<p class="empty-state">Ошибка поиска услуг</p>';
    }
}

async function loadReviews() {
    const user = authData?.user;
    const container = document.getElementById('reviewsList');

    if (!container) {
        console.error('loadReviews: элемент #reviewsList не найден');
        return;
    }

    if (!user || !user.id) {
        console.error('loadReviews: нет пользователя в authData:', authData);
        container.innerHTML = `
            <div class="card" style="text-align: center; padding: 40px; background: #fff3cd; border-color: #ffc107;">
                <p style="font-size: 48px; margin-bottom: 20px;">⚠️</p>
                <h3 style="margin-bottom: 10px; color: #856404;">Вы не авторизованы</h3>
                <p style="color: #856404;">Пожалуйста, войдите в систему, чтобы просмотреть отзывы.</p>
            </div>
        `;
        return;
    }

    try {
        const bookingsRes = await authFetch(`/api/owners/${user.id}/bookings`);

        if (!bookingsRes) {
            throw new Error('Сервер не вернул ответ по бронированиям (bookingsRes = null)');
        }

        console.log('loadReviews: статус /bookings =', bookingsRes.status);

        if (bookingsRes.status === 204) {
            console.log('loadReviews: 204 No Content — считаем, что бронирований нет');
            renderEmptyReviews(container);
            return;
        }

        if (!bookingsRes.ok) {
            const text = await bookingsRes.text().catch(() => '');
            console.error('loadReviews: ошибка ответа /bookings:', bookingsRes.status, text);
            throw new Error(`Ошибка загрузки бронирований (код ${bookingsRes.status})`);
        }

        let bookings;
        try {
            bookings = await bookingsRes.json();
        } catch (jsonErr) {
            console.error('loadReviews: не удалось распарсить JSON бронирований:', jsonErr);
            throw new Error('Некорректный JSON от сервера при загрузке бронирований');
        }

        console.log('loadReviews: bookings JSON =', bookings);

        let safeBookings;
        if (Array.isArray(bookings)) {
            safeBookings = bookings;
        } else if (bookings === null || typeof bookings === 'undefined') {
            console.warn('loadReviews: bookings = null/undefined, используем пустой массив');
            safeBookings = [];
        } else {
            console.warn('loadReviews: bookings не массив, но продолжаем как с пустым массивом. Значение:', bookings);
            safeBookings = [];
        }

        const completedBookings = safeBookings.filter(b => b && b.status === 'completed');
        console.log('loadReviews: completedBookings =', completedBookings);

        const results = await Promise.allSettled(
            completedBookings.map(async (booking) => {
                try {
                    const reviewRes = await authFetch(`/api/bookings/${booking.booking_id}/review`);

                    if (!reviewRes) {
                        console.warn('loadReviews: reviewRes = null для booking', booking.booking_id);
                        return null;
                    }

                    if (!reviewRes.ok) {
                        if (reviewRes.status !== 404) {
                            const txt = await reviewRes.text().catch(() => '');
                            console.warn(
                                'loadReviews: неожиданный статус при загрузке отзыва',
                                booking.booking_id,
                                reviewRes.status,
                                txt
                            );
                        }
                        return null;
                    }

                    const review = await reviewRes.json();
                    return { ...review, booking };
                } catch (e) {
                    console.error('loadReviews: ошибка при загрузке отзыва для booking', booking.booking_id, e);
                    return null;
                }
            })
        );

        const reviews = results
            .filter(r => r.status === 'fulfilled' && r.value)
            .map(r => r.value);

        console.log('loadReviews: итоговый список reviews =', reviews);

        if (reviews.length === 0) {
            renderEmptyReviews(container);
            return;
        }

        container.innerHTML = reviews.map(review => `
            <div class="card">
                <div style="display: flex; justify-content: space-between; align-items: start;">
                    <div style="flex: 1;">
                        <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 10px;">
                            <strong>Рейтинг:</strong>
                            <span style="color: #f39c12; font-size: 20px;">
                                ${'⭐'.repeat(review.rating)}${'☆'.repeat(5 - review.rating)}
                            </span>
                            <span style="color: #666;">(${review.rating}/5)</span>
                        </div>

                        <p><strong>Бронирование:</strong> #${review.booking_id}</p>
                        <p><strong>Няня:</strong> ID ${review.sitter_id}</p>
                        <p><strong>Дата отзыва:</strong> ${new Date(review.created_at).toLocaleDateString('ru-RU')}</p>

                        ${review.comment ? `
                            <div style="margin-top: 15px; padding: 15px; background: #f9f9f9; border-radius: 8px; border-left: 4px solid #667eea;">
                                <strong>Ваш комментарий:</strong>
                                <p style="margin-top: 8px; line-height: 1.6;">${review.comment}</p>
                            </div>
                        ` : '<p style="color: #999; font-style: italic;">Комментарий не оставлен</p>'}
                    </div>

                    <div style="display: flex; flex-direction: column; gap: 10px; min-width: 140px;">
                        <button onclick="editReview(${review.review_id}, ${review.rating}, \`${(review.comment || '').replace(/`/g, '\\`')}\`)"
                                class="btn btn-secondary">
                            ✏️ Редактировать
                        </button>

                        <button onclick="deleteReview(${review.review_id})"
                                class="btn btn-danger">
                            🗑️ Удалить
                        </button>
                    </div>
                </div>
            </div>
        `).join('');

    } catch (err) {
        console.error('Ошибка загрузки отзывов:', err);

        container.innerHTML = `
            <div class="card" style="text-align: center; padding: 40px; background: #fff3cd; border-color: #ffc107;">
                <p style="font-size: 48px; margin-bottom: 20px;">⚠️</p>
                <h3 style="margin-bottom: 10px; color: #856404;">Ошибка загрузки отзывов</h3>
                <p style="color: #856404;">${err.message}</p>
                <button onclick="loadReviews()" class="btn btn-primary" style="margin-top: 20px;">
                    🔄 Попробовать снова
                </button>
            </div>
        `;
    }
}

function renderEmptyReviews(container) {
    container.innerHTML = `
        <div class="card" style="text-align: center; padding: 40px;">
            <p style="font-size: 48px; margin-bottom: 20px;">⭐</p>
            <h3 style="margin-bottom: 10px;">У вас пока нет отзывов</h3>
            <p style="color: #666;">Отзывы появятся после завершения бронирований</p>
        </div>
    `;
}

async function editReview(reviewId, currentRating, currentComment) {
    const modal = document.createElement('div');
    modal.className = 'modal active';
    modal.innerHTML = `
        <div class="modal-content">
            <span class="close" onclick="this.parentElement.parentElement.remove()">&times;</span>
            <h2>✏️ Редактировать отзыв</h2>
            <form id="editReviewForm">
                <div class="form-group">
                    <label>Рейтинг (1-5):</label>
                    <select id="editReviewRating" required>
                        <option value="5" ${currentRating === 5 ? 'selected' : ''}>⭐⭐⭐⭐⭐ Отлично</option>
                        <option value="4" ${currentRating === 4 ? 'selected' : ''}>⭐⭐⭐⭐ Хорошо</option>
                        <option value="3" ${currentRating === 3 ? 'selected' : ''}>⭐⭐⭐ Средне</option>
                        <option value="2" ${currentRating === 2 ? 'selected' : ''}>⭐⭐ Плохо</option>
                        <option value="1" ${currentRating === 1 ? 'selected' : ''}>⭐ Ужасно</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Комментарий:</label>
                    <textarea id="editReviewComment" rows="4" placeholder="Расскажите о вашем опыте...">${currentComment}</textarea>
                </div>
                <div style="display: flex; gap: 10px; justify-content: flex-end;">
                    <button type="button" onclick="this.closest('.modal').remove()" class="btn btn-secondary">
                        Отмена
                    </button>
                    <button type="submit" class="btn btn-primary">
                        💾 Сохранить изменения
                    </button>
                </div>
            </form>
        </div>
    `;
    document.body.appendChild(modal);

    document.getElementById('editReviewForm').addEventListener('submit', async (e) => {
        e.preventDefault();

        const newRating = parseInt(document.getElementById('editReviewRating').value);
        const newComment = document.getElementById('editReviewComment').value.trim();

        try {
            const res = await authFetch(`/api/reviews/${reviewId}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    rating: newRating,
                    comment: newComment
                })
            });

            if (res.ok) {
                alert('✅ Отзыв успешно обновлён!');
                modal.remove();
                loadReviews();
            } else {
                const err = await res.json().catch(() => ({}));
                alert('❌ Ошибка: ' + (err.error || `код ${res.status}`));
            }
        } catch (err) {
            console.error('Ошибка обновления отзыва:', err);
            alert('❌ Ошибка соединения с сервером');
        }
    });
}

async function deleteReview(reviewId) {
    if (!confirm('🗑️ Вы уверены, что хотите удалить этот отзыв?\n\nЭто действие нельзя отменить.')) {
        return;
    }

    try {
        const res = await authFetch(`/api/reviews/${reviewId}`, {
            method: 'DELETE'
        });

        if (res.ok) {
            alert('✅ Отзыв успешно удалён');
            loadReviews();
        } else {
            const err = await res.json().catch(() => ({}));
            alert('❌ Ошибка удаления: ' + (err.error || `код ${res.status}`));
        }
    } catch (err) {
        console.error('Ошибка удаления отзыва:', err);
        alert('❌ Ошибка соединения с сервером');
    }
}

function showReviewModal(bookingId, sitterId) {
    const modal = document.createElement('div');
    modal.className = 'modal active';
    modal.innerHTML = `
        <div class="modal-content">
            <span class="close" onclick="this.parentElement.parentElement.remove()">&times;</span>
            <h2>⭐ Оставить отзыв</h2>
            <p style="color: #666; margin-bottom: 20px;">
                Расскажите о вашем опыте с няней
            </p>
            <form id="reviewForm">
                <div class="form-group">
                    <label>Рейтинг (1-5): <span style="color: red;">*</span></label>
                    <select id="reviewRating" required>
                        <option value="">-- Выберите рейтинг --</option>
                        <option value="5">⭐⭐⭐⭐⭐ Отлично</option>
                        <option value="4">⭐⭐⭐⭐ Хорошо</option>
                        <option value="3">⭐⭐⭐ Средне</option>
                        <option value="2">⭐⭐ Плохо</option>
                        <option value="1">⭐ Ужасно</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Комментарий:</label>
                    <textarea id="reviewComment" rows="5" 
                              placeholder="Что вам понравилось или не понравилось?&#10;Как няня обращалась с питомцем?&#10;Рекомендуете ли вы эту няню другим?"></textarea>
                    <small style="color: #666;">Комментарий необязателен, но будет полезен другим владельцам</small>
                </div>
                <div style="display: flex; gap: 10px; justify-content: flex-end;">
                    <button type="button" onclick="this.closest('.modal').remove()" class="btn btn-secondary">
                        Отмена
                    </button>
                    <button type="submit" class="btn btn-primary">
                        📤 Отправить отзыв
                    </button>
                </div>
            </form>
        </div>
    `;
    document.body.appendChild(modal);

    document.getElementById('reviewForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        await submitReview(bookingId, sitterId);
        modal.remove();
    });
}

async function submitReview(bookingId, sitterId) {
    const user = JSON.parse(localStorage.getItem('user'));
    const rating = parseInt(document.getElementById('reviewRating').value);
    const comment = document.getElementById('reviewComment').value.trim();

    if (!rating) {
        alert('❌ Пожалуйста, выберите рейтинг');
        return;
    }

    try {
        const res = await authFetch('/api/reviews', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                booking_id: bookingId,
                owner_id: user.id,
                sitter_id: sitterId,
                rating: rating,
                comment: comment
            })
        });

        if (res.ok) {
            const result = await res.json().catch(() => ({}));
            console.log('submitReview: результат создания отзыва =', result);
            alert('✅ Спасибо за ваш отзыв!\n\nОн поможет другим владельцам выбрать няню.');

            document.querySelector('[data-tab="reviews"]').click();
            loadReviews();
        } else {
            const err = await res.json().catch(() => ({}));
            alert('❌ Ошибка: ' + (err.error || `код ${res.status}`));
        }
    } catch (err) {
        console.error('Ошибка отправки отзыва:', err);
        alert('❌ Ошибка соединения с сервером');
    }
}

function showAddPetModal() {
    document.getElementById('addPetModal').classList.add('active');
}

function showCreateBookingModal() {
    loadPetsForBooking();
    document.getElementById('createBookingModal').classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

async function loadPetsForBooking() {
    try {
        const res = await authFetch(`/api/owners/${user.id}/pets`);
        if (!res) return;
        const pets = await res.json();

        const select = document.getElementById('bookingPetSelect');
        select.innerHTML = pets.map(pet =>
            `<option value="${pet.pet_id}">${pet.name} (${pet.type})</option>`
        ).join('');
    } catch (err) {
        console.error('Ошибка загрузки питомцев:', err);
    }
}

document.getElementById('addPetForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData);

    if (data.age) {
        data.age = Number(data.age);
    }
    data.owner_id = Number(user.id);

    try {
        const res = await authFetch('/api/pets', {
            method: 'POST',
            body: JSON.stringify(data)
        });

        if (res && res.ok) {
            alert('✅ Питомец добавлен!');
            closeModal('addPetModal');
            e.target.reset();
            loadPets();
            loadOverview();
        } else if (res) {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка соединения');
    }
});

document.getElementById('createBookingForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData);
    data.owner_id = user.id;

    data.start_time = new Date(data.start_time).toISOString();
    data.end_time = new Date(data.end_time).toISOString();

    try {
        const res = await authFetch('/api/bookings', {
            method: 'POST',
            body: JSON.stringify(data)
        });

        if (res && res.ok) {
            alert('✅ Бронирование создано!');
            closeModal('createBookingModal');
            e.target.reset();
            loadBookings();
            loadOverview();
        } else if (res) {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка соединения');
    }
});

async function deletePet(petId) {
    if (!confirm('Удалить питомца?')) return;

    try {
        const res = await authFetch(`/api/pets/${petId}`, { method: 'DELETE' });
        if (res && res.ok) {
            alert('✅ Питомец удалён');
            loadPets();
            loadOverview();
        }
    } catch (err) {
        alert('Ошибка удаления');
    }
}

async function cancelBooking(bookingId) {
    if (!confirm('Отменить бронирование?')) return;

    try {
        const res = await authFetch(`/api/bookings/${bookingId}/cancel`, { method: 'POST' });
        if (res && res.ok) {
            alert('✅ Бронирование отменено');
            loadBookings();
            loadOverview();
        }
    } catch (err) {
        alert('Ошибка отмены');
    }
}

function getPetTypeIcon(type) {
    const icons = { cat: '🐱', dog: '🐕', rodent: '🐹' };
    return icons[type] || '🐾';
}

function getPetAgeWord(age) {
    if (age === 1) return 'год';
    if (age >= 2 && age <= 4) return 'года';
    return 'лет';
}

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
    localStorage.removeItem('auth');
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    window.location.href = 'login.html';
}

loadOverview();
loadPets();
