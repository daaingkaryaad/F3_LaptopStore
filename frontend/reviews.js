const laptopSelect = document.getElementById("laptopSelect");
const filterLaptop = document.getElementById("filterLaptop");
const reviewList = document.getElementById("reviewList");
const commentInput = document.getElementById("commentInput");
const ratingSelect = document.getElementById("ratingSelect");
const formMsg = document.getElementById("formMsg");

function getQueryParam(name) {
  const u = new URL(window.location.href);
  return u.searchParams.get(name) || "";
}

async function loadLaptops() {
  const laptops = await fetch("/api/laptops").then((r) => r.json());

  laptopSelect.innerHTML = "";
  filterLaptop.innerHTML = "";

  laptops.forEach((l) => {
    const opt1 = document.createElement("option");
    opt1.value = l.id;
    opt1.textContent = l.model_name;
    laptopSelect.appendChild(opt1);

    const opt2 = document.createElement("option");
    opt2.value = l.id;
    opt2.textContent = l.model_name;
    filterLaptop.appendChild(opt2);
  });

  const preselect = getQueryParam("product_id");
  const first = preselect && laptops.some((x) => x.id === preselect) ? preselect : (laptops[0] && laptops[0].id) || "";
  if (first) {
    filterLaptop.value = first;
    laptopSelect.value = first;
    loadReviews(first);
  }
}

async function loadReviews(productId) {
  reviewList.innerHTML = `<p class="empty">Loading reviews...</p>`;

  const token = window.RapidTech.getToken();
  const headers = token ? { Authorization: `Bearer ${token}` } : {};

  const res = await fetch(`/api/reviews?product_id=${encodeURIComponent(productId)}`, { headers });
  const data = await res.json();

  if (!Array.isArray(data) || data.length === 0) {
    reviewList.innerHTML = `<p class="empty">No reviews yet.</p>`;
    return;
  }

  reviewList.innerHTML = "";
  data.forEach((r) => {
    const card = document.createElement("div");
    card.className = "review-card";
    card.innerHTML = `
      <div class="meta">User: ${r.user_id}</div>
      <div class="rating">Rating: ${r.rating}/5</div>
      <div class="comment">${r.comment}</div>
      <div class="meta">Status: ${r.status}</div>
    `;
    reviewList.appendChild(card);
  });
}

filterLaptop.addEventListener("change", (e) => {
  loadReviews(e.target.value);
});

document.getElementById("submitReview").addEventListener("click", async () => {
  if (!window.RapidTech.requireAuthOrRedirect()) return;

  formMsg.textContent = "";
  formMsg.style.color = "";

  const payload = {
    laptop_id: laptopSelect.value,
    rating: Number(ratingSelect.value),
    comment: commentInput.value.trim(),
  };

  try {
    await window.RapidTech.apiFetch("/api/reviews", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    formMsg.style.color = "#1db954";
    formMsg.textContent = "Review submitted (pending approval).";
    commentInput.value = "";
    loadReviews(filterLaptop.value);
  } catch (err) {
    formMsg.textContent = err.message || "Failed to submit review.";
  }
});

loadLaptops();
