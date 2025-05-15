
document.addEventListener("DOMContentLoaded", () => {
    const nameAgeEl = document.getElementById("user-name-age");
    const descriptionEl = document.getElementById("user-description");
    const bioEl = document.getElementById("user-bio");
    const interestsContainer = document.getElementById("user-interests");
    const carouselInner = document.getElementById("carousel-inner");

    const likeBtn = document.querySelector(".like-btn");
    const dislikeBtn = document.querySelector(".dislike-btn");

    let currentUserId = null;

    function loadUser() {
        fetch("/api/user/preferences")
            .then(res => res.json())
            .then(data => {
                const profileContainer = document.querySelector(".profile-container");
                const noUsersMessage = document.getElementById("no-users-message");
    
                if (data.message === "No more preferences") {
                    profileContainer.style.display = "none";
                    noUsersMessage.style.display = "block";
                    currentUserId = null;
                    return;
                }
    
                profileContainer.style.display = "block";
                noUsersMessage.style.display = "none";
    
                const user = data.user;
                currentUserId = user.id;
    
                nameAgeEl.textContent = `${user.name}, ${user.age}`;
                descriptionEl.textContent = user.description || "";
                bioEl.textContent = user.city || "";
    
                interestsContainer.innerHTML = "";
                if (user.interests) {
                    user.interests.forEach(interest => {
                        const span = document.createElement("span");
                        span.className = "badge bg-primary me-1";
                        span.textContent = interest;
                        interestsContainer.appendChild(span);
                    });
                }
    
                carouselInner.innerHTML = "";
                if (user.photos && user.photos.length > 0) {
                    user.photos.forEach((url, index) => {
                        const item = document.createElement("div");
                        item.className = `carousel-item ${index === 0 ? "active" : ""}`;
                        item.innerHTML = `<img src="${url}" class="d-block w-100" alt="photo ${index + 1}" />`;
                        carouselInner.appendChild(item);
                    });
                }
            })
            .catch(err => {
                console.error("Ошибка при получении пользователя:", err);
            });
    }

    likeBtn.addEventListener("click", () => {
        if (!currentUserId) return;
        fetch(`/api/user/love/${currentUserId}`, {
            method: "POST"
        }).then(() => {
            loadUser();
        });
    });

    dislikeBtn.addEventListener("click", () => {
        loadUser();
    });

    loadUser();
});