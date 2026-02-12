(async function () {
  if (!window.RapidTech.requireAuthOrRedirect()) return;
  if (window.RapidTech.getRoleFromToken() !== "admin") {
    document.body.innerHTML = '<div class="empty">Forbidden. Not an admin.</div>';
    return;
  }

  const productsTable = document.getElementById("productsTable");
  const adminMsg = document.getElementById("adminMsg");
  const reloadBtn = document.getElementById("reloadBtn");
  const createForm = document.getElementById("createForm");
  const createMsg = document.getElementById("createMsg");

  function laptopToRow(p) {
    const safe = (v) => (v == null ? "" : String(v));
    return `
      <div class="admin-row" data-id="${p.id}">
        <div class="admin-main">
          <div class="admin-title">
            <strong>${safe(p.model_name)}</strong>
            <span class="muted">(${safe(p.brand_id)} • ${safe(p.category_id)})</span>
          </div>
          <div class="admin-id muted">${p.id}</div>
        </div>

        <div class="admin-fields">
          <label>Price <input type="number" class="field-price" value="${Number(p.price)||0}"></label>
          <label>Stock <input type="number" class="field-stock" value="${Number(p.stock)||0}"></label>
          <label class="check"><input type="checkbox" class="field-active" ${p.is_active ? "checked" : ""}> Active</label>
        </div>

        <details class="admin-details">
          <summary>Specs</summary>
          <div class="spec-grid">
            <label>CPU <input class="s-cpu" value="${safe(p.specs?.cpu)}"></label>
            <label>RAM <input class="s-ram" value="${safe(p.specs?.ram)}"></label>
            <label>GPU <input class="s-gpu" value="${safe(p.specs?.gpu)}"></label>
            <label>Storage <input class="s-storage" value="${safe(p.specs?.storage)}"></label>
            <label>Storage type <input class="s-storage-type" value="${safe(p.specs?.storage_type)}"></label>
            <label>Screen size <input class="s-screen-size" value="${safe(p.specs?.screen_size)}"></label>
            <label>Resolution <input class="s-screen-res" value="${safe(p.specs?.screen_resolution)}"></label>
          </div>
        </details>

        <div class="admin-actions">
          <button class="btn primary btn-save" type="button">Save</button>
          <button class="btn danger btn-delete" type="button">Delete</button>
        </div>

        <div class="admin-row-msg muted"></div>
      </div>
    `;
  }

  async function loadProducts() {
    adminMsg.textContent = "";
    productsTable.innerHTML = "";
    try {
      const products = await fetch("/api/laptops?include_inactive=true").then(r => r.json());
      productsTable.innerHTML = products.map(laptopToRow).join("");
    } catch (e) {
      adminMsg.textContent = e.message || "Failed to load products";
    }
  }

  productsTable.addEventListener("click", async (e) => {
    const row = e.target.closest(".admin-row");
    if (!row) return;
    const id = row.dataset.id;
    const msg = row.querySelector(".admin-row-msg");

    if (e.target.closest(".btn-save")) {
      msg.textContent = "";
      const payload = {
        model_name: row.querySelector("strong")?.textContent || "",
        brand_id: row.querySelector(".admin-title .muted")?.textContent?.split("•")[0]?.replace(/[()]/g,"").trim() || "",
        category_id: row.querySelector(".admin-title .muted")?.textContent?.split("•")[1]?.replace(/[()]/g,"").trim() || "",
        price: Number(row.querySelector(".field-price").value || 0),
        stock: Number(row.querySelector(".field-stock").value || 0),
        is_active: row.querySelector(".field-active").checked,
        specs: {
          cpu: row.querySelector(".s-cpu").value,
          ram: row.querySelector(".s-ram").value,
          gpu: row.querySelector(".s-gpu").value,
          storage: row.querySelector(".s-storage").value,
          storage_type: row.querySelector(".s-storage-type").value,
          screen_size: row.querySelector(".s-screen-size").value,
          screen_resolution: row.querySelector(".s-screen-res").value,
        }
      };

      try {
        await window.RapidTech.apiFetch(`/api/laptops/${id}`, {
          method: "PUT",
          body: JSON.stringify(payload),
        });
        msg.textContent = "Saved.";
      } catch (err) {
        msg.textContent = err.message || "Save failed";
      }
    }

    if (e.target.closest(".btn-delete")) {
      msg.textContent = "";
      if (!confirm("Delete this laptop?")) return;
      try {
        await window.RapidTech.apiFetch(`/api/laptops/${id}`, { method: "DELETE" });
        row.remove();
      } catch (err) {
        msg.textContent = err.message || "Delete failed";
      }
    }
  });

  reloadBtn.addEventListener("click", loadProducts);

  createForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    createMsg.textContent = "";
    const fd = new FormData(createForm);
    const payload = {
      model_name: fd.get("model_name"),
      brand_id: fd.get("brand_id"),
      category_id: fd.get("category_id"),
      price: Number(fd.get("price") || 0),
      stock: Number(fd.get("stock") || 0),
      is_active: true,
      specs: {
        cpu: fd.get("cpu") || "",
        ram: fd.get("ram") || "",
        gpu: fd.get("gpu") || "",
        storage: fd.get("storage") || "",
        storage_type: fd.get("storage_type") || "",
        screen_size: fd.get("screen_size") || "",
        screen_resolution: fd.get("screen_resolution") || "",
      },
    };

    try {
      await window.RapidTech.apiFetch("/api/laptops", { method: "POST", body: JSON.stringify(payload) });
      createMsg.textContent = "Created.";
      createForm.reset();
      await loadProducts();
    } catch (err) {
      createMsg.textContent = err.message || "Create failed";
    }
  });

  await loadProducts();
})();