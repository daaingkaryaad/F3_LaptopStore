async function loadLaptops() {
  const res = await fetch("/api/laptops");
  const data = await res.json();

  const grid = document.getElementById("laptopGrid");
  grid.innerHTML = "";

  data.forEach((l) => {
    const out = l.stock === 0;
    const card = document.createElement("article");
    card.className = `laptop-card ${out ? "out" : ""}`;

    card.innerHTML = `
      <div class="laptop-name">${l.model_name}</div>
      <div class="laptop-meta">${l.brand_id} • ${l.category_id}</div>

      <ul class="spec-list">
        <li><span>CPU</span><span>${l.specs.cpu}</span></li>
        <li><span>RAM</span><span>${l.specs.ram}</span></li>
        <li><span>GPU</span><span>${l.specs.gpu}</span></li>
        <li><span>Storage</span><span>${l.specs.storage} (${l.specs.storage_type})</span></li>
        <li><span>Screen</span><span>${l.specs.screen_size} (${l.specs.screen_resolution})</span></li>
      </ul>

      <div class="laptop-price">₸${l.price.toLocaleString()}</div>
      <div class="laptop-meta">${out ? "Out of stock" : `In stock: ${l.stock}`}</div>

      <div class="card-actions">
        <button class="btn btn-primary add-to-cart" ${out ? "disabled" : ""} data-id="${l.id}">
          Add to Compare Cart
        </button>
        <a class="btn" href="reviews.html?product_id=${l.id}">Reviews</a>
      </div>
    `;
    grid.appendChild(card);
  });
}

document.addEventListener("click", async (e) => {
  const btn = e.target.closest(".add-to-cart");
  if (!btn) return;

  if (!window.RapidTech.requireAuthOrRedirect()) return;

  const laptop_id = btn.dataset.id;
  btn.disabled = true;
  const original = btn.textContent;
  btn.textContent = "Adding...";

  try {
    await window.RapidTech.apiFetch("/api/cart/items", {
      method: "POST",
      body: JSON.stringify({ laptop_id, quantity: 1 }),
    });
    btn.textContent = "Added";
    setTimeout(() => {
      btn.textContent = original;
      btn.disabled = false;
    }, 900);
  } catch (err) {
    btn.textContent = "Failed";
    btn.disabled = false;
    alert(err.message || "Failed to add to cart");
  }
});

loadLaptops();