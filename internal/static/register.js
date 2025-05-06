const form = document.getElementById('registrationForm');
const photoInput = document.getElementById('photo');
const previewContainer = document.getElementById('photoPreviewContainer'); // контейнер для превью

let createdUserId = null;

// Функция для показа превью нескольких фото
photoInput.addEventListener('change', () => {
    previewContainer.innerHTML = ''; // очищаем контейнер
  
    const files = photoInput.files;
    for (const file of files) {
      if (!file.type.startsWith('image/')) continue;
  
      const img = document.createElement('img');
      img.classList.add('photo-preview');
      img.style.width = '120px';
      img.style.height = '120px';
      img.style.objectFit = 'cover';
      img.style.borderRadius = '8px';
      img.style.margin = '5px';
  
      // --- Важно! Используем FileReader для получения DataURL ---
      const reader = new FileReader();
      reader.onload = (e) => {
        img.src = e.target.result; // вот тут правильный src
      };
      reader.readAsDataURL(file);
  
      previewContainer.appendChild(img);
    }
  });
// Получение координат через Promise
function getCoordinates() {
    return new Promise((resolve) => {
      if (!navigator.geolocation) {
        // Если геолокация не поддерживается, сразу возвращаем дефолт
        resolve({ latitude: 55.7558, longitude: 37.6176 });
        return;
      }
      navigator.geolocation.getCurrentPosition(
        (pos) => resolve({ latitude: pos.coords.latitude, longitude: pos.coords.longitude }),
        (err) => {
          // При ошибке возвращаем дефолтные координаты
          resolve({ latitude: 55.7558, longitude: 37.6176 });
        },
        { timeout: 10000 }
      );
    });
  }
form.addEventListener('submit', async (e) => {
  e.preventDefault();

  // Получаем координаты пользователя
  const coords = await getCoordinates();
  console.log(coords)
  const formData = new FormData(form);

  // Формируем объект данных для создания пользователя
  const data = {
    name: formData.get('name'),
    description: formData.get('description'),
    city: formData.get('city'),
    age: parseInt(formData.get('age')),
    email: formData.get('email'),
    password: formData.get('password'),
    coordinates: coords ? `(${coords.latitude.toFixed(6)}, ${coords.longitude.toFixed(6)})` : null,
  };

  try {
    // Создаём пользователя (без фото)
    const res = await fetch('http://localhost:8080/api/user/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    });
    
    const user_add = await res.json();
    
    if (!res.ok) {
      throw new Error('Ошибка создания пользователя: ' + user_add["error"]);
    }
    
    createdUserId = user_add.user_id;

    // Загружаем все выбранные фото по отдельным запросам
    const files = photoInput.files;
    for (const file of files) {
      const photoForm = new FormData();
      photoForm.append('files', file);
      const uploadRes = await fetch(`http://localhost:8080/upload/${createdUserId}`, {
        method: 'POST',
        body: photoForm
      });
      if (!uploadRes.ok) throw new Error('Ошибка загрузки фото');
    }

    alert('Пользователь создан! Заполните предпочтения.');
    window.location.href = `preferences/${createdUserId}`;

  } catch (err) {
    alert('Ошибка: ' + err.message);
  }
});
