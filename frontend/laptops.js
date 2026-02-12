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
    `;
    grid.appendChild(card);
  });
}

loadLaptops();