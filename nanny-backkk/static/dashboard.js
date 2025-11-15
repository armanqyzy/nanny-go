// Проверка авторизации
const user = JSON.parse(localStorage.getItem('user'));
if (!user || user.role !== 'owner') {
    window.location.href = '/login';
}

document.getElementById('userEmail').textContent = user.email;

// Переключение табов
document.querySelectorAll('.sidebar-menu a').forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        const tab = e.target.dataset.tab;
        
        // Убираем active класс со всех ссылок
        document.querySelectorAll('.sidebar-menu a').forEach(l => l.classList.remove('active'));
        e.target.classList.add('active');
        
        // Скрываем все табы
        document.querySelectorAll('.tab-content').forEach(t => t.style.display = 'none');
        
        // Показываем нужный таб
        document.getElementById(tab + '-tab').style.display = 'block';
        
        // Загружаем данные
        loadTabData(tab);
    });
});

// Загрузка данных таба
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

// Загрузка обзора
async function loadOverview() {
    try {
        // Загружаем количество питомцев
        const petsRes = await fetch(`/api/owners/${user.id}/pets`);
        const pets = await petsRes.json();
        document.getElementById('petsCount').textContent = pets.length || 0;
        
        // Загружаем бронирования
        const bookingsRes = await fetch(`/api/owners/${user.id}/bookings`);
        const bookings = await bookingsRes.json();
        document.getElementById('bookingsCount').textContent = bookings.length || 0;
        
        // Отображаем последние 5 бронирований
        const recentDiv = document.getElementById('recentBookings');
        if (bookings.length === 0) {
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

// Загрузка питомцев
async function loadPets() {
    try {
        const res = await fetch(`/api/owners/${user.id}/pets`);
        const pets = await petsRes.json();
        
        const petsDiv = document.getElementById('petsList');
        
        if (pets.length === 0) {
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
    }
}

// Загрузка бронирований
async function loadBookings() {
    try {
        const res = await fetch(`/api/owners/${user.id}/bookings`);
        const bookings = await res.json();
        
        const bookingsDiv = document.getElementById('bookingsList');
        
        if (bookings.length === 0) {
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

// Поиск услуг
async function searchServices() {
    const type = document.getElementById('serviceTypeFilter').value;
    const location = document.getElementById('locationFilter').value;
    
    try {
        const params = new URLSearchParams();
        if (type) params.append('type', type);
        if (location) params.append('location', location);
        
        const res = await fetch(`/api/services/search?${params}`);
        const services = await res.json();
        
        const resultsDiv = document.getElementById('searchResults');
        
        if (services.length === 0) {
            resultsDiv.innerHTML = '<p class="empty-state">Услуги не найдены</p>';
            return;
        }
        
        resultsDiv.innerHTML = services.map(s => `
            <div class="card" style="margin-bottom: 15px;">
                <h3>${s.sitter_name}</h3>
                <div class="rating">
                    ${renderStars(s.sitter_rating)}
                    <span>(${s.sitter_rating.toFixed(1)})</span>
                </div>
                <p><strong>Услуга:</strong> ${getServiceTypeName(s.type)}</p>
                <p><strong>Цена:</strong> ${s.price_per_hour} ₸/час</p>
                <p>${s.description || ''}</p>
                <button class="btn btn-primary" onclick="bookService(${s.sitter_id}, ${s.service_id})">Забронировать</button>
            </div>
        `).join('');
    } catch (err) {
        console.error('Ошибка поиска:', err);
    }
}

// Загрузка отзывов
async function loadReviews() {
    // TODO: Реализовать загрузку отзывов владельца
    document.getElementById('reviewsList').innerHTML = '<p>Функция в разработке</p>';
}

// Модальные окна
function showAddPetModal() {
    document.getElementById('addPetModal').classList.add('active');
}

function showCreateBookingModal() {
    // Загружаем список питомцев
    loadPetsForBooking();
    document.getElementById('createBookingModal').classList.add('active');
}

function closeModal(modalId) {
    document.getElementById(modalId).classList.remove('active');
}

async function loadPetsForBooking() {
    try {
        const res = await fetch(`/api/owners/${user.id}/pets`);
        const pets = await res.json();
        
        const select = document.getElementById('bookingPetSelect');
        select.innerHTML = pets.map(pet => 
            `<option value="${pet.pet_id}">${pet.name} (${pet.type})</option>`
        ).join('');
    } catch (err) {
        console.error('Ошибка загрузки питомцев:', err);
    }
}

// Обработчики форм
document.getElementById('addPetForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const data = Object.fromEntries(formData);
    data.owner_id = user.id;
    
    try {
        const res = await fetch('/api/pets', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        
        if (res.ok) {
            alert('✅ Питомец добавлен!');
            closeModal('addPetModal');
            e.target.reset();
            loadPets();
        } else {
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
    
    // Конвертируем время в ISO формат
    data.start_time = new Date(data.start_time).toISOString();
    data.end_time = new Date(data.end_time).toISOString();
    
    try {
        const res = await fetch('/api/bookings', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        
        if (res.ok) {
            alert('✅ Бронирование создано!');
            closeModal('createBookingModal');
            e.target.reset();
            loadBookings();
        } else {
            const err = await res.json();
            alert('❌ ' + err.error);
        }
    } catch (err) {
        alert('Ошибка соединения');
    }
});

// Удалить питомца
async function deletePet(petId) {
    if (!confirm('Удалить питомца?')) return;
    
    try {
        const res = await fetch(`/api/pets/${petId}`, { method: 'DELETE' });
        if (res.ok) {
            alert('✅ Питомец удалён');
            loadPets();
        }
    } catch (err) {
        alert('Ошибка удаления');
    }
}

// Отменить бронирование
async function cancelBooking(bookingId) {
    if (!confirm('Отменить бронирование?')) return;
    
    try {
        const res = await fetch(`/api/bookings/${bookingId}/cancel`, { method: 'POST' });
        if (res.ok) {
            alert('✅ Бронирование отменено');
            loadBookings();
        }
    } catch (err) {
        alert('Ошибка отмены');
    }
}

// Вспомогательные функции
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
    localStorage.removeItem('user');
    window.location.href = '/login';
}

// Загружаем обзор при старте
loadOverview();
