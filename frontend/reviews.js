const laptopSelect = document.getElementById("laptopSelect");
const filterLaptop = document.getElementById("filterLaptop");
const reviewList = document.getElementById("reviewList");
const tokenInput = document.getElementById("tokenInput");
const commentInput = document.getElementById("commentInput");
const ratingSelect = document.getElementById("ratingSelect");
const formMsg = document.getElementById("formMsg");

tokenInput.value = localStorage.getItem("jwt") || "";

async function loadLaptops() {
  const res = await fetch("/api/laptops");
  const laptops = await res.json();

  laptopSelect.innerHTML = "";
  filterLaptop.innerHTML = "";

  laptops.forEach(l => {
    const opt1 = document.createElement("option");
    opt1.value = l.id;
    opt1.textContent = l.model_name;
    laptopSelect.appendChild(opt1);

    const opt2 = document.createElement("option");
    opt2.value = l.id;
    opt2.textContent = l.model_name;
    filterLaptop.appendChild(opt2);
  });

  if (laptops.length > 0) {
    filterLaptop.value = laptops[0].id;
    loadReviews(laptops[0].id);
  }
}

async function loadReviews(productId) {
  reviewList.innerHTML = `<p class="empty">Loading reviews...</p>`;
  const res = await fetch(`/api/reviews?product_id=${productId}`, {
    headers: { "Authorization": `Bearer ${tokenInput.value.trim()}` }
  });
  const data = await res.json();

  if (!Array.isArray(data) || data.length === 0) {
    reviewList.innerHTML = `<p class="empty">No reviews yet.</p>`;
    return;
  }

  reviewList.innerHTML = "";
  data.forEach(r => {
    const card = document.createElement("div");
    card.className = "review-card";
    card.innerHTML = `
      <div class="meta">User: ${r.user_id}</div>
      <div class="rating">Rating: ${r.rating}/5</div>
      <div class="comment">${r.comment}</div>
    `;
    reviewList.appendChild(card);
  });
}

filterLaptop.addEventListener("change", (e) => {
  loadReviews(e.target.value);
});

document.getElementById("submitReview").addEventListener("click", async () => {
  const token = tokenInput.value.trim();
  if (!token) {
    formMsg.textContent = "Token required.";
    return;
  }

  localStorage.setItem("jwt", token);

  const payload = {
    laptop_id: laptopSelect.value,
    rating: Number(ratingSelect.value),
    comment: commentInput.value.trim()
  };

  const res = await fetch("/api/reviews", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) {
    const err = await res.json();
    formMsg.textContent = err.error || "Failed to submit review.";
    return;
  }

  formMsg.style.color = "#1db954";
  formMsg.textContent = "Review submitted (pending approval).";
  commentInput.value = "";
  loadReviews(filterLaptop.value);
});

loadLaptops();