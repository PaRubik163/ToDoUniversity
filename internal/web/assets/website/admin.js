(function () {
  const tokenKey = "access_token";
  const list = document.getElementById("container-list");
  const project = document.getElementById("project");
  const message = document.getElementById("admin-message");
  const buttons = ["refresh", "restart", "delete"].map((id) => document.getElementById(id));

  async function ensureToken() {
    let token = localStorage.getItem(tokenKey);
    if (!token) {
      const response = await fetch("/api/auth/refresh", { method: "POST", credentials: "same-origin" });
      if (response.ok) {
        const result = await response.json();
        token = result.access_token;
        localStorage.setItem(tokenKey, token);
        localStorage.setItem("telegram_linked", String(result.telegram_linked));
      }
    }
    if (!token) { location.replace("/"); throw new Error("Нужно войти в аккаунт."); }
    return token;
  }

  async function api(path, options) {
    const token = await ensureToken();
    const response = await fetch(path, Object.assign({}, options || {}, { headers: Object.assign({}, (options && options.headers) || {}, { Authorization: "Bearer " + token }) }));
    if (response.status === 401) {
      localStorage.removeItem(tokenKey);
      location.replace("/");
      throw new Error("Сеанс истёк. Войдите снова.");
    }
    if (response.status === 403) throw new Error("У этого аккаунта нет доступа к администрированию.");
    if (!response.ok) throw new Error(await response.text() || "Ошибка запроса.");
    return response;
  }

  async function loadContainers() {
    message.textContent = "";
    list.replaceChildren();
    try {
      const response = await api("/api/admin/containers");
      const result = await response.json();
      project.textContent = "Compose-проект: " + result.project;
      if (!result.containers.length) {
        const empty = document.createElement("p");
        empty.textContent = "Контейнеры проекта не найдены.";
        list.append(empty);
        return;
      }
      result.containers.forEach((container) => {
        const row = document.createElement("div");
        row.className = "container-row";
        const name = document.createElement("span");
        name.className = "container-name";
        name.textContent = container.service || container.name;
        const state = document.createElement("span");
        state.className = "container-state" + (container.state === "running" ? " running" : "");
        state.textContent = container.state;
        row.append(name, state);
        list.append(row);
      });
    } catch (error) { message.textContent = error.message; }
  }

  async function perform(action) {
    const confirmation = action === "delete"
      ? "Удалить все контейнеры проекта? Бот, сайт, PostgreSQL и Redis остановятся. Docker volume с данными этой командой не удаляется."
      : "Перезапустить все контейнеры проекта?"
    if (!window.confirm(confirmation)) return;
    buttons.forEach((button) => { button.disabled = true; });
    message.textContent = action === "delete" ? "Отправлена команда на удаление контейнеров…" : "Отправлена команда на перезапуск…";
    try {
      await api("/api/admin/containers/" + action, { method: "POST" });
      message.textContent = action === "delete"
        ? "Команда принята. Для возврата контейнеров выполните docker compose up -d на сервере."
        : "Перезапуск запущен. Обновите страницу через несколько секунд.";
      if (action === "restart") setTimeout(loadContainers, 12000);
      if (action === "delete") list.replaceChildren();
    } catch (error) { message.textContent = error.message; }
    finally { buttons.forEach((button) => { button.disabled = false; }); }
  }

  document.getElementById("refresh").addEventListener("click", loadContainers);
  document.getElementById("restart").addEventListener("click", () => perform("restart"));
  document.getElementById("delete").addEventListener("click", () => perform("delete"));
  loadContainers();
})();
