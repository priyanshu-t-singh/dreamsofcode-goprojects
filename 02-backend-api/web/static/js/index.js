const BASE_URL = "http://localhost:8080";

// Check API Health on load
async function checkHealth() {
  const dot = document.getElementById("status-dot");
  const text = document.getElementById("status-text");
  try {
    const res = await fetch(`${BASE_URL}/`);
    if (res.ok) {
      dot.className = "status-dot online";
      text.textContent = "API Online";
    } else {
      throw new Error();
    }
  } catch {
    dot.className = "status-dot offline";
    text.textContent = "API Offline";
  }
}

// Helper to build headers with optional Authorization
function getHeaders() {
  const token = document.getElementById("bearer-token").value.trim();
  const headers = { "Content-Type": "application/json" };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  return headers;
}

// Tab Switcher
function switchMode(mode) {
  document
    .querySelectorAll(".tab")
    .forEach((t) => t.classList.remove("active"));
  if (mode === "basic") {
    document.querySelectorAll(".tab")[0].classList.add("active");
    document.getElementById("mode-basic").classList.remove("hidden");
    document.getElementById("mode-sum").classList.add("hidden");
  } else {
    document.querySelectorAll(".tab")[1].classList.add("active");
    document.getElementById("mode-basic").classList.add("hidden");
    document.getElementById("mode-sum").classList.remove("hidden");
  }
}

// Handle Basic Operations (add, subtract, multiply, divide)
async function executeOp(operation) {
  const resultEl = document.getElementById("result");
  const a = parseInt(document.getElementById("num-a").value, 10);
  const b = parseInt(document.getElementById("num-b").value, 10);

  if (isNaN(a) || isNaN(b)) {
    resultEl.className = "result-value error";
    resultEl.textContent = "Enter valid numbers";
    return;
  }

  // Adjust field names for divide vs others based on OpenAPI spec
  let payload = {};
  if (operation === "divide") {
    payload = { dividend: a, divisor: b };
  } else {
    payload = { a: a, b: b };
  }

  try {
    const res = await fetch(`${BASE_URL}/api/v1/${operation}`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(payload),
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.message || `HTTP ${res.status}`);
    }

    resultEl.className = "result-value";
    resultEl.textContent = data.result;
  } catch (err) {
    resultEl.className = "result-value error";
    resultEl.textContent = err.message || "Error";
  }
}

// Handle Array Sum
async function executeSum() {
  const resultEl = document.getElementById("result");
  const rawInput = document.getElementById("array-input").value;

  // Parse comma-separated string into integer array
  const numbers = rawInput
    .split(",")
    .map((n) => parseInt(n.trim(), 10))
    .filter((n) => !isNaN(n));

  if (numbers.length === 0) {
    resultEl.className = "result-value error";
    resultEl.textContent = "Provide numbers";
    return;
  }

  try {
    const res = await fetch(`${BASE_URL}/api/v1/sum`, {
      method: "POST",
      headers: getHeaders(),
      body: JSON.stringify(numbers),
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.message || `HTTP ${res.status}`);
    }

    resultEl.className = "result-value";
    resultEl.textContent = data.result;
  } catch (err) {
    resultEl.className = "result-value error";
    resultEl.textContent = err.message || "Error";
  }
}

checkHealth();
