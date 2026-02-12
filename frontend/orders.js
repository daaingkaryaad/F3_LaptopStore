(async function () {
  if (!window.RapidTech.requireAuthOrRedirect()) return;

  const selectedEl = document.getElementById("selectedItems");
  const totalEl = document.getElementById("orderTotal");
  const ordersListEl = document.getElementById("ordersList");
  const payForm = document.getElementById("payForm");
  const payMsg = document.getElementById("payMsg");

  function getCheckoutIDs() {
    try {
      return JSON.parse(localStorage.getItem("rapidtech_checkout_ids") || "[]");
    } catch {
      return [];
    }
  }

  async function loadProductsById() {
    const all = await fetch("/api/laptops").then(r => r.json());
    const byId = new Map(all.map(p => [p.id, p]));
    return byId;
  }

  async function loadCart() {
    return await window.RapidTech.apiFetch("/api/cart");
  }

  function calcTotal(items) {
    return items.reduce((sum, it) => sum + (Number(it.price)||0) * (it.qty||1), 0);
  }

  function renderSelected(items) {
    if (!items.length) {
      selectedEl.innerHTML = `<div class="empty">No selected items. Go back and select laptops in Compare.</div>`;
      totalEl.textContent = "₸0";
      return;
    }
    selectedEl.innerHTML = items.map(it => `
      <div class="order-item">
        <div>
          <div class="order-name">${it.model_name}</div>
          <div class="muted">${it.brand_id} • Qty: ${it.qty}</div>
        </div>
        <div class="order-price">${window.RapidTech.formatMoneyKZT((Number(it.price)||0) * it.qty)}</div>
      </div>
    `).join("");
    totalEl.textContent = window.RapidTech.formatMoneyKZT(calcTotal(items));
  }

  async function renderOrders() {
    try {
      const orders = await window.RapidTech.apiFetch("/api/orders");
      if (!Array.isArray(orders) || orders.length === 0) {
        ordersListEl.innerHTML = `<div class="empty">No orders yet.</div>`;
        return;
      }
      ordersListEl.innerHTML = orders.map(o => `
        <div class="order-card">
          <div class="order-card-head">
            <div><strong>Order</strong> <span class="muted">${o.id}</span></div>
            <div class="muted">${new Date(o.created_at).toLocaleString()}</div>
          </div>
          <div class="order-card-body">
            <div class="muted">Status: <strong>${o.status}</strong></div>
            <div class="muted">Total: <strong>${window.RapidTech.formatMoneyKZT(o.total)}</strong></div>
            <div class="order-lines">
              ${(o.items||[]).map(it => `<div class="order-line">• ${it.laptop_id} × ${it.quantity} @ ${window.RapidTech.formatMoneyKZT(it.price)}</div>`).join("")}
            </div>
          </div>
        </div>
      `).join("");
    } catch (e) {
      ordersListEl.innerHTML = `<div class="empty">${e.message || "Failed to load orders"}</div>`;
    }
  }

  let selected = [];
  try {
    const ids = getCheckoutIDs();
    const [cart, byId] = await Promise.all([loadCart(), loadProductsById()]);
    const cartItems = Array.isArray(cart.items) ? cart.items : [];
    selected = cartItems
      .filter(ci => ids.includes(ci.laptop_id))
      .map(ci => {
        const p = byId.get(ci.laptop_id);
        if (!p) return null;
        return { ...p, qty: ci.quantity || 1 };
      })
      .filter(Boolean);
  } catch (e) {
    selectedEl.innerHTML = `<div class="empty">${e.message || "Failed to load checkout items"}</div>`;
  }

  renderSelected(selected);
  await renderOrders();

  payForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    payMsg.textContent = "";
    if (selected.length === 0) {
      payMsg.textContent = "Nothing to order.";
      return;
    }

    const ids = selected.map(x => x.id);
    try {
      const order = await window.RapidTech.apiFetch("/api/orders", {
        method: "POST",
        body: JSON.stringify({ item_ids: ids }),
      });
      localStorage.removeItem("rapidtech_checkout_ids");
      payMsg.textContent = `Order created: ${order.id}. (No real money was harmed.)`;
      selectedEl.innerHTML = `<div class="empty">Order placed. Go back to Compare to pick more.</div>`;
      totalEl.textContent = "₸0";
      selected = [];
      await renderOrders();
    } catch (err) {
      payMsg.textContent = err.message || "Failed to place order";
    }
  });
})();