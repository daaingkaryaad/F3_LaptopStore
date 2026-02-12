(async function () {
  const compareBtn = document.getElementById("compareBtn");
  const checkoutBtn = document.getElementById("checkoutBtn");
  const winnerBox = document.getElementById("winnerBox");
  const productsEl = document.getElementById("products");
  

  function parseFirstNumber(text) {
    const m = String(text || "").match(/([0-9]+(\.[0-9]+)?)/);
    return m ? Number(m[1]) : 0;
  }

  function parseRAMGB(ramText) {
    const n = parseFirstNumber(ramText);
    return Number.isFinite(n) ? n : 0;
  }

  function parseStorageGB(storageText) {
    const raw = String(storageText || "").toLowerCase();
    const n = parseFirstNumber(raw);
    if (!Number.isFinite(n)) return 0;
    if (raw.includes("tb")) return n * 1024;
    return n;
  }

  function parseScreenSizeInches(sizeText) {
    const n = parseFirstNumber(sizeText);
    return Number.isFinite(n) ? n : 0;
  }

  function calculateValueScore(laptop) {
    const price = Number(laptop.price) || 0;
    const ram = parseRAMGB(laptop.specs?.ram);
    const storage = parseStorageGB(laptop.specs?.storage);
    const screen = parseScreenSizeInches(laptop.specs?.screen_size);

    const utility = ram * 2.2 + (storage / 256) * 2.0 + screen * 0.35;
    if (price <= 0) return 0;
    return (utility / (price / 1000)) * 10;
  }

  function renderEmpty(message) {
    productsEl.innerHTML = `<div class="empty">${message}</div>`;
    winnerBox.textContent = "No comparison yet";
  }

  async function fetchCartAndProducts() {
    if (!window.RapidTech.requireAuthOrRedirect()) return [];

    const cart = await window.RapidTech.apiFetch("/api/cart");
    const items = Array.isArray(cart.items) ? cart.items : [];
    if (items.length === 0) return [];

    const products = await fetch("/api/laptops").then((r) => r.json());
    const byId = new Map(products.map((p) => [p.id, p]));

    return items
      .map((it) => {
        const p = byId.get(it.laptop_id);
        if (!p) return null;
        return { ...p, _qty: it.quantity || 1 };
      })
      .filter(Boolean);
  }

  function renderProducts(items) {
    productsEl.innerHTML = "";
    items.forEach((l) => {
      const card = document.createElement("article");
      card.className = "product-card";
      card.dataset.id = l.id;

      card.innerHTML = `
        <h2>${l.model_name}</h2>
        <div class="muted">${l.brand_id} • ${l.category_id}</div>
        <ul>
          <li>Price: ${window.RapidTech.formatMoneyKZT(l.price)}</li>
          <li>CPU: ${l.specs?.cpu || "-"}</li>
          <li>RAM: ${l.specs?.ram || "-"}</li>
          <li>GPU: ${l.specs?.gpu || "-"}</li>
          <li>Storage: ${l.specs?.storage || "-"} (${l.specs?.storage_type || "-"})</li>
          <li>Screen: ${l.specs?.screen_size || "-"} (${l.specs?.screen_resolution || "-"})</li>
        </ul>
        <div class="card-row"><label class="check"><input type="checkbox" class="select-item" checked> Select</label><span class="qty">Qty: ${l._qty}</span></div>
        <p class="score">Value score: <span>-</span></p>
      `;
      productsEl.appendChild(card);
    });
  }

  let currentItems = [];

  try {
    currentItems = await fetchCartAndProducts();
  } catch (err) {
    renderEmpty(err.message || "Failed to load cart");
    return;
  }

  if (currentItems.length === 0) {
    renderEmpty("Your compare cart is empty. Add laptops from the All Laptops page.");
    return;
  }

  renderProducts(currentItems);

  compareBtn.addEventListener("click", () => {
    const cards = Array.from(document.querySelectorAll(".product-card"));
    if (cards.length < 2) {
      winnerBox.textContent = "Add at least 2 laptops to compare.";
      return;
    }

    let bestId = "";
    let bestScore = -Infinity;

    cards.forEach((card) => {
      card.classList.remove("winner");

      const id = card.dataset.id;
      const laptop = currentItems.find((x) => x.id === id);
      const score = calculateValueScore(laptop);
      card.querySelector(".score span").textContent = score.toFixed(2);

      if (score > bestScore) {
        bestScore = score;
        bestId = id;
      }
    });

    const bestCard = cards.find((c) => c.dataset.id === bestId);
    const bestLaptop = currentItems.find((x) => x.id === bestId);
    if (bestCard && bestLaptop) {
      bestCard.classList.add("winner");
      winnerBox.textContent = `Best value: ${bestLaptop.model_name} (value score ${bestScore.toFixed(2)})`;
    }
  });

  checkoutBtn.addEventListener("click", () => {
    const cards = Array.from(document.querySelectorAll(".product-card"));
    const selected = cards
      .filter(c => c.querySelector(".select-item")?.checked)
      .map(c => c.dataset.id);

    if (selected.length === 0) {
      winnerBox.textContent = "Select at least 1 laptop to order.";
      return;
    }
    localStorage.setItem("rapidtech_checkout_ids", JSON.stringify(selected));
    window.location.href = "orders.html";
  });

  document.getElementById("checkoutBtn")?.addEventListener("click", () => {
  const selected = getSelectedProductIds(); 
  if (!selected.length) {
    alert("Select at least one laptop to order");
    return;
  }

  localStorage.setItem("order_selection", JSON.stringify(selected));
  window.location.href = "/orders.html";
});


})();
