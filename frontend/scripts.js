const AUTH_TOKEN_KEY = "rapidtech_token";

function getToken() {
  return localStorage.getItem(AUTH_TOKEN_KEY) || "";
}

function setToken(token) {
  if (!token) {
    localStorage.removeItem(AUTH_TOKEN_KEY);
    return;
  }
  localStorage.setItem(AUTH_TOKEN_KEY, token);
}

function formatMoneyKZT(value) {
  const n = Number(value);
  if (!Number.isFinite(n)) return String(value);
  return "₸" + n.toLocaleString();
}

async function apiFetch(url, options = {}) {
  const token = getToken();
  const headers = new Headers(options.headers || {});
  if (!headers.has("Content-Type") && options.body && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }
  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const res = await fetch(url, { ...options, headers });

  let payload;
  const ct = res.headers.get("content-type") || "";
  if (ct.includes("application/json")) {
    payload = await res.json().catch(() => null);
  } else {
    payload = await res.text().catch(() => "");
  }

  if (!res.ok) {
    const msg = payload && payload.error ? payload.error : `Request failed (${res.status})`;
    const err = new Error(msg);
    err.status = res.status;
    throw err;
  }
  return payload;
}

function requireAuthOrRedirect() {
  const token = getToken();
  if (!token) {
    window.location.href = "login.html";
    return false;
  }
  return true;
}

window.RapidTech = {
  getToken,
  setToken,
  formatMoneyKZT,
  apiFetch,
  requireAuthOrRedirect,
};


function parseJwtPayload(token) {
  try {
    const part = token.split(".")[1];
    const json = atob(part.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(decodeURIComponent(escape(json)));
  } catch {
    return null;
  }
}

function getRoleFromToken() {
  const t = getToken();
  if (!t) return "";
  const p = parseJwtPayload(t);
  return (p && p.role) ? p.role : "";
}

function logout() {
  setToken("");
  window.location.href = "index.html";
}

function ensureHeaderNav() {
  const nav = document.querySelector(".header-nav ul");
  if (!nav) return;

  nav.querySelectorAll("[data-injected='true']").forEach(el => el.remove());

  const role = getRoleFromToken();
  const token = getToken();

  if (token) {
    const li = document.createElement("li");
    li.dataset.injected = "true";
    li.innerHTML = `<a href="orders.html">Orders</a>`;
    nav.appendChild(li);

    const liLogout = document.createElement("li");
    liLogout.dataset.injected = "true";
    liLogout.innerHTML = `<a class="nav-pill" href="#" id="navLogoutBtn">Logout</a>`;
    nav.appendChild(liLogout);
    liLogout.querySelector("#navLogoutBtn").addEventListener("click", (e) => {
      e.preventDefault();
      logout();
    });
  }

  if (token) {
    nav.querySelectorAll("a[href='login.html'], a[href='signup.html']").forEach(a => {
      const li = a.closest("li");
      if (li) li.style.display = "none";
    });
  }

  if (role === "admin") {
    const liA = document.createElement("li");
    liA.dataset.injected = "true";
    liA.innerHTML = `<a href="admin.html">Admin</a>`;
    nav.appendChild(liA);

    const liR = document.createElement("li");
    liR.dataset.injected = "true";
    liR.innerHTML = `<a href="admin_reviews.html">Review Moderation</a>`;
    nav.appendChild(liR);
  }
}

document.addEventListener("DOMContentLoaded", ensureHeaderNav);

window.RapidTech.getRoleFromToken = getRoleFromToken;
window.RapidTech.logout = logout;
