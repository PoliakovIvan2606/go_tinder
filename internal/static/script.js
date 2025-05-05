document.addEventListener('DOMContentLoaded', function() {
    // Получаем id из URL
    function getUserIdFromUrl() {
        const pathParts = window.location.pathname.split('/');
        return pathParts[pathParts.length - 1];
    }
    const userId = getUserIdFromUrl();

    async function fetchUserData() {
        try {
            const response = await fetch(`http://127.0.0.1:8080/api/user/${userId}`);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            updateProfile(data.user);
        } catch (error) {
            console.error('Error fetching user data:', error);
            document.getElementById('user-name-age').textContent = "Ошибка загрузки данных";
            document.getElementById('user-description').textContent = "Попробуйте обновить страницу";
        }
    }

    // Функция для обновления профиля на основе полученных данных
    function updateProfile(userData) {
        // Обновляем имя и возраст
        document.getElementById('user-name-age').textContent = `${userData.name}, ${userData.age}`;
        
        // Обновляем описание
        document.getElementById('user-description').textContent = userData.description;
        document.getElementById('user-bio').textContent = userData.description;
        
        // Обновляем город в личной информации
        const personalInfo = document.getElementById('user-personal-info');
        personalInfo.innerHTML = `<strong>${userData.city}</strong> • ${personalInfo.innerHTML.split('•').slice(1).join('•')}`;
        
        // Обновляем фотографии в карусели
        updateCarousel(userData.photos);
        
        // Обновляем статус отношений (мужской/женский вариант)
        const statusElement = document.getElementById('user-status');
        if (userData.name.endsWith('а') || userData.name.endsWith('я')) {
            statusElement.textContent = "Свободна";
        } else {
            statusElement.textContent = "Свободен";
        }
    }

    // Функция для обновления карусели с фотографиями
    function updateCarousel(photos) {
        const carouselInner = document.getElementById('carousel-inner');
        carouselInner.innerHTML = '';
        
        if (!photos || photos.length === 0) {
            // Если фотографий нет, добавляем заглушку
            carouselInner.innerHTML = `
                <div class="carousel-item active">
                    <img src="https://via.placeholder.com/600x800?text=No+photo" alt="Нет фото" />
                </div>
            `;
            return;
        }
        
        // Добавляем фотографии в карусель
        photos.forEach((photo, index) => {
            const carouselItem = document.createElement('div');
            carouselItem.className = `carousel-item ${index === 0 ? 'active' : ''}`;
            carouselItem.innerHTML = `
                <img src="${photo}" alt="Фото пользователя ${index + 1}" />
            `;
            carouselInner.appendChild(carouselItem);
        });
    }

    // Вызываем функцию загрузки данных при загрузке страницы
    fetchUserData();
});