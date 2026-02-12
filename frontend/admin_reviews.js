(async function () {
  if (!window.RapidTech.requireAuthOrRedirect()) return;
  if (window.RapidTech.getRoleFromToken() !== "admin") {
    document.body.innerHTML = '<div class="empty">Forbidden. Not an admin.</div>';
    return;
  }

  const listEl = document.getElementById("pendingReviews");
  const msgEl = document.getElementById("reviewsMsg");
  const reloadBtn = document.getElementById("reloadReviews");

  function render(reviews, productsById) {
    if (!reviews.length) {
      listEl.innerHTML = '<div class="empty">No pending reviews.</div>';
      return;
    }
    listEl.innerHTML = reviews.map(r => {
      const p = productsById.get(r.laptop_id);
      const title = p ? p.model_name : r.laptop_id;
      return `
        <div class="review-admin-card" data-id="${r.id}">
          <div class="review-admin-head">
            <div>
              <strong>${title}</strong>
              <div class="muted">Rating: ${r.rating} • Status: ${r.status}</div>
            </div>
            <div class="review-admin-actions">
              <button class="btn primary btn-approve" type="button">Approve</button>
              <button class="btn danger btn-reject" type="button">Reject</button>
            </div>
          </div>
          <div class="review-admin-body">${(r.comment || "").replace(/</g,"&lt;")}</div>
          <div class="muted small">Review ID: ${r.id}</div>
          <div class="review-admin-msg muted"></div>
        </div>
      `;
    }).join("");
  }

  async function loadPending() {
    msgEl.textContent = "";
    listEl.innerHTML = "";
    try {
      const products = await fetch("/api/laptops?include_inactive=true").then(r => r.json());
      const byId = new Map(products.map(p => [p.id, p]));

      const allReviews = [];
      for (const p of products) {
        const rs = await window.RapidTech.apiFetch(`/api/reviews?product_id=${encodeURIComponent(p.id)}&all=true`);
        if (Array.isArray(rs)) allReviews.push(...rs);
      }
      const pending = allReviews.filter(r => r.status === "pending");
      render(pending, byId);
    } catch (e) {
      msgEl.textContent = e.message || "Failed to load pending reviews";
      listEl.innerHTML = '<div class="empty">Failed.</div>';
    }
  }

  listEl.addEventListener("click", async (e) => {
    const card = e.target.closest(".review-admin-card");
    if (!card) return;
    const id = card.dataset.id;
    const msg = card.querySelector(".review-admin-msg");

    if (e.target.closest(".btn-approve")) {
      try {
        await window.RapidTech.apiFetch(`/api/reviews/${id}/approve`, { method: "PUT", body: JSON.stringify({ status: "approved" }) });
        card.remove();
      } catch (err) {
        msg.textContent = err.message || "Approve failed";
      }
    }

    if (e.target.closest(".btn-reject")) {
      try {
        await window.RapidTech.apiFetch(`/api/reviews/${id}/approve`, { method: "PUT", body: JSON.stringify({ status: "rejected" }) });
        card.remove();
      } catch (err) {
        msg.textContent = err.message || "Reject failed";
      }
    }
  });

  reloadBtn.addEventListener("click", loadPending);

  await loadPending();
})();