(function () {
  "use strict";

  const tg = window.Telegram && window.Telegram.WebApp;
  const initData = tg ? tg.initData : "";

  const listEl = document.getElementById("list");
  const emptyEl = document.getElementById("empty-state");
  const statusEl = document.getElementById("status");
  const sectionButtons = Array.from(document.querySelectorAll(".section-button"));
  const analyticsEl = document.getElementById("analytics");
  const weekRangeEl = document.getElementById("week-range");
  const statCompletedEl = document.getElementById("stat-completed");
  const statMonthCompletedEl = document.getElementById("stat-month-completed");
  const statUncompletedEl = document.getElementById("stat-uncompleted");
  const statOverdueEl = document.getElementById("stat-overdue");
  const monthPickerEl = document.getElementById("month-picker");
  const monthPrevEl = document.getElementById("month-prev");
  const monthNextEl = document.getElementById("month-next");
  const monthChartTitleEl = document.getElementById("month-chart-title");
  const weekChartEl = document.getElementById("week-chart");
  const monthChartEl = document.getElementById("month-chart");

  const formEl = document.getElementById("add-form");
  const formTitleEl = document.getElementById("form-title");
  const formErrorEl = document.getElementById("form-error");
  const addToggleEl = document.getElementById("add-toggle");
  const cancelAddEl = document.getElementById("cancel-add");
  const subjectInput = document.getElementById("field-subject");
  const titleInput = document.getElementById("field-title");
  const deadlineInput = document.getElementById("field-deadline");

  const settingsToggleEl = document.getElementById("settings-toggle");
  const settingsFormEl = document.getElementById("settings-form");
  const settingsErrorEl = document.getElementById("settings-error");
  const cancelSettingsEl = document.getElementById("cancel-settings");
  const remindTimeInput = document.getElementById("field-remind-time");

  let editingId = null;
  let currentSection = "uncompleted";
  let homeworkItems = [];
  let selectedMonth = "";

  if (tg) {
    tg.ready();
    tg.expand();
  }

  if (!initData) {
    setStatus("Откройте это приложение через Telegram, чтобы увидеть свои задания.");
  } else {
    loadHomework();
  }

  sectionButtons.forEach(function (button) {
    button.addEventListener("click", function () {
      currentSection = button.dataset.section;
      sectionButtons.forEach(function (item) { item.classList.toggle("active", item === button); });
      renderCurrentSection();
    });
  });

  monthPickerEl.addEventListener("change", function () {
    if (!monthPickerEl.value) return;
    selectedMonth = monthPickerEl.value;
    loadAnalytics();
  });
  monthPrevEl.addEventListener("click", function () { shiftSelectedMonth(-1); });
  monthNextEl.addEventListener("click", function () { shiftSelectedMonth(1); });

  addToggleEl.addEventListener("click", function () {
    const wasHidden = formEl.classList.contains("hidden");
    settingsFormEl.classList.add("hidden");
    if (wasHidden) {
      openAddForm();
    } else {
      formEl.classList.add("hidden");
    }
  });

  cancelAddEl.addEventListener("click", function () {
    closeAddForm();
  });

  formEl.addEventListener("submit", function (event) {
    event.preventDefault();
    hideFormError();

    const payload = {
      subject: subjectInput.value.trim(),
      title: titleInput.value.trim(),
      deadline: deadlineInput.value,
    };
    if (!payload.subject || !payload.title || !payload.deadline) {
      showFormError("Заполните все поля.");
      return;
    }

    if (editingId) {
      updateHomework(editingId, payload);
    } else {
      createHomework(payload);
    }
  });

  settingsToggleEl.addEventListener("click", function () {
    const wasHidden = settingsFormEl.classList.contains("hidden");
    formEl.classList.add("hidden");
    if (wasHidden) {
      openSettingsForm();
    } else {
      settingsFormEl.classList.add("hidden");
    }
  });

  cancelSettingsEl.addEventListener("click", function () {
    hideSettingsError();
    settingsFormEl.classList.add("hidden");
  });

  settingsFormEl.addEventListener("submit", function (event) {
    event.preventDefault();
    hideSettingsError();
    if (!remindTimeInput.value) {
      showSettingsError("Укажите время.");
      return;
    }
    saveSettings(remindTimeInput.value);
  });

  function openAddForm() {
    editingId = null;
    formEl.reset();
    hideFormError();
    formTitleEl.textContent = "Новое задание";
    formEl.querySelector('button[type="submit"]').textContent = "Сохранить";
    formEl.classList.remove("hidden");
    subjectInput.focus();
  }

  function openEditForm(item) {
    editingId = item.id;
    hideFormError();
    subjectInput.value = item.subject;
    titleInput.value = item.title;
    deadlineInput.value = item.deadline;
    formTitleEl.textContent = "Редактирование задания №" + item.id;
    formEl.querySelector('button[type="submit"]').textContent = "Сохранить изменения";
    settingsFormEl.classList.add("hidden");
    formEl.classList.remove("hidden");
    subjectInput.focus();
  }

  function closeAddForm() {
    editingId = null;
    formEl.reset();
    hideFormError();
    formEl.classList.add("hidden");
  }

  function openSettingsForm() {
    hideSettingsError();
    settingsFormEl.classList.remove("hidden");
    remindTimeInput.value = "";
    apiFetch("/api/settings")
      .then(function (settings) {
        remindTimeInput.value = (settings && settings.reminder_time) || "09:00";
      })
      .catch(function () {
        remindTimeInput.value = "09:00";
      });
  }

  function saveSettings(reminderTime) {
    const submitButton = settingsFormEl.querySelector('button[type="submit"]');
    submitButton.disabled = true;
    apiFetch("/api/settings", { method: "PUT", body: JSON.stringify({ reminder_time: reminderTime }) })
      .then(function () {
        haptic("notificationSuccess");
        settingsFormEl.classList.add("hidden");
      })
      .catch(function (err) {
        showSettingsError(err.message);
      })
      .finally(function () {
        submitButton.disabled = false;
      });
  }

  function loadHomework() {
    setStatus("Загрузка…");
    Promise.all([apiFetch("/api/homework"), fetchAnalytics()])
      .then(function (results) {
        homeworkItems = results[0] || [];
        renderCurrentSection();
        renderAnalytics(results[1] || {});
        setStatus("");
      })
      .catch(function (err) {
        setStatus("Не удалось загрузить список: " + err.message);
      });
  }

  function createHomework(payload) {
    const submitButton = formEl.querySelector('button[type="submit"]');
    submitButton.disabled = true;
    apiFetch("/api/homework", { method: "POST", body: JSON.stringify(payload) })
      .then(function () {
        closeAddForm();
        haptic("notificationSuccess");
        loadHomework();
      })
      .catch(function (err) {
        showFormError(err.message);
      })
      .finally(function () {
        submitButton.disabled = false;
      });
  }

  function updateHomework(id, payload) {
    const submitButton = formEl.querySelector('button[type="submit"]');
    submitButton.disabled = true;
    apiFetch("/api/homework/" + id, { method: "PUT", body: JSON.stringify(payload) })
      .then(function () {
        closeAddForm();
        haptic("notificationSuccess");
        loadHomework();
      })
      .catch(function (err) {
        showFormError(err.message);
      })
      .finally(function () {
        submitButton.disabled = false;
      });
  }

  function deleteHomework(id, itemEl) {
    const confirmed = tg && tg.showConfirm
      ? new Promise(function (resolve) {
          tg.showConfirm("Удалить это задание?", resolve);
        })
      : Promise.resolve(window.confirm("Удалить это задание?"));

    confirmed.then(function (ok) {
      if (!ok) return;
      apiFetch("/api/homework/" + id, { method: "DELETE" })
        .then(function () {
          haptic("notificationSuccess");
          loadHomework();
        })
        .catch(function (err) {
          setStatus("Не удалось удалить: " + err.message);
        });
    });
  }

  function renderList(items) {
    listEl.innerHTML = "";
    if (!items.length) {
      emptyEl.textContent = emptyMessage();
      emptyEl.classList.remove("hidden");
      return;
    }
    emptyEl.classList.add("hidden");

    items.forEach(function (item) {
      listEl.appendChild(renderItem(item));
    });
  }

  function renderCurrentSection() {
    const analyticsMode = currentSection === "analytics";
    analyticsEl.classList.toggle("hidden", !analyticsMode);
    listEl.classList.toggle("hidden", analyticsMode);
    if (analyticsMode) {
      emptyEl.classList.add("hidden");
      return;
    }
    const items = homeworkItems.filter(function (item) {
      if (currentSection === "completed") return item.completed;
      if (currentSection === "overdue") return !item.completed && item.days_left < 0;
      return !item.completed;
    });
    renderList(items);
  }

  function emptyMessage() {
    if (currentSection === "completed") return "Пока нет сделанных заданий.";
    if (currentSection === "overdue") return "Просроченных заданий нет.";
    return "Несделанных заданий нет. Нажмите «+», чтобы добавить.";
  }

  function renderAnalytics(data) {
    statCompletedEl.textContent = data.completed_this_week || 0;
    statMonthCompletedEl.textContent = data.completed_this_month || 0;
    statUncompletedEl.textContent = data.uncompleted || 0;
    statOverdueEl.textContent = data.overdue || 0;
    weekRangeEl.textContent = "Неделя: " + formatDate(data.week_start) + " — " + formatDate(data.week_end);
    selectedMonth = data.month;
    monthPickerEl.value = data.month;
    monthPickerEl.max = data.current_month;
    monthNextEl.disabled = data.month >= data.current_month;
    monthChartTitleEl.textContent = "Выполнено по дням — " + formatMonth(data.month);
    renderChart(weekChartEl, data.week_days || [], true);
    renderChart(monthChartEl, data.month_days || [], false);
  }

  function fetchAnalytics() {
    const path = selectedMonth ? "/api/analytics?month=" + encodeURIComponent(selectedMonth) : "/api/analytics";
    return apiFetch(path);
  }

  function loadAnalytics() {
    fetchAnalytics()
      .then(renderAnalytics)
      .catch(function (err) { setStatus("Не удалось загрузить аналитику: " + err.message); });
  }

  function shiftSelectedMonth(amount) {
    if (!selectedMonth) return;
    const parts = selectedMonth.split("-");
    const date = new Date(Number(parts[0]), Number(parts[1]) - 1 + amount, 1);
    selectedMonth = date.getFullYear() + "-" + String(date.getMonth() + 1).padStart(2, "0");
    monthPickerEl.value = selectedMonth;
    loadAnalytics();
  }

  function formatMonth(month) {
    const parts = month.split("-");
    return new Date(Number(parts[0]), Number(parts[1]) - 1, 1)
      .toLocaleDateString("ru-RU", { month: "long", year: "numeric" });
  }

  function renderChart(container, days, showWeekday) {
    container.innerHTML = "";
    const max = Math.max(1, ...days.map(function (day) { return day.count; }));
    days.forEach(function (day) {
      const dateParts = day.date.split("-").map(Number);
      const date = new Date(dateParts[0], dateParts[1] - 1, dateParts[2]);
      const column = document.createElement("div");
      column.className = "chart-column";
      column.title = formatDate(day.date) + ": " + day.count;

      const value = document.createElement("span");
      value.className = "chart-value";
      value.textContent = day.count;

      const track = document.createElement("div");
      track.className = "chart-track";
      const bar = document.createElement("i");
      bar.className = "chart-bar";
      bar.style.height = (day.count === 0 ? 3 : Math.max(8, day.count / max * 100)) + "%";
      track.appendChild(bar);

      const label = document.createElement("span");
      label.className = "chart-label";
      label.textContent = showWeekday
        ? date.toLocaleDateString("ru-RU", { weekday: "short" })
        : String(date.getDate());

      column.appendChild(value);
      column.appendChild(track);
      column.appendChild(label);
      container.appendChild(column);
    });
  }

  function renderItem(item) {
    const li = document.createElement("li");
    li.className = "item";

    const body = document.createElement("div");
    body.className = "item-body";

    const subject = document.createElement("p");
    subject.className = "item-subject";
    subject.textContent = item.subject;

    const title = document.createElement("p");
    title.className = "item-title";
    title.textContent = item.title;

    const deadline = document.createElement("p");
    deadline.className = "item-deadline " + deadlineClass(item.days_left);
    deadline.textContent = formatDeadline(item.deadline, item.days_left);

    body.appendChild(subject);
    body.appendChild(title);
    body.appendChild(deadline);

    if (item.completed) {
      const completedDate = document.createElement("p");
      completedDate.className = "item-completed";
      completedDate.textContent = item.completed_at ? "Сделано: " + formatDate(item.completed_at) : "Сделано";
      body.appendChild(completedDate);
    }

    const actions = document.createElement("div");
    actions.className = "item-actions";

    const completeButton = document.createElement("button");
    completeButton.className = "icon-action complete-button" + (item.completed ? " is-completed" : "");
    completeButton.type = "button";
    completeButton.setAttribute("aria-label", item.completed ? "Вернуть в несделанные" : "Отметить сделанным");
    completeButton.textContent = item.completed ? "↶" : "✓";
    completeButton.addEventListener("click", function () {
      apiFetch("/api/homework/" + item.id + "/status", {
        method: "PUT",
        body: JSON.stringify({ completed: !item.completed }),
      }).then(function () {
        haptic("notificationSuccess");
        loadHomework();
      }).catch(function (err) {
        setStatus("Не удалось изменить статус: " + err.message);
      });
    });

    const editButton = document.createElement("button");
    editButton.className = "icon-action edit-button";
    editButton.type = "button";
    editButton.setAttribute("aria-label", "Редактировать");
    editButton.textContent = "✎";
    editButton.addEventListener("click", function () {
      openEditForm(item);
    });

    const deleteButton = document.createElement("button");
    deleteButton.className = "icon-action delete-button";
    deleteButton.type = "button";
    deleteButton.setAttribute("aria-label", "Удалить");
    deleteButton.textContent = "✕";
    deleteButton.addEventListener("click", function () {
      deleteHomework(item.id, li);
    });

    actions.appendChild(completeButton);
    if (!item.completed) actions.appendChild(editButton);
    actions.appendChild(deleteButton);

    li.appendChild(body);
    li.appendChild(actions);
    return li;
  }

  function formatDate(dateStr) {
    const parts = dateStr.split("-");
    return parts[2] + "." + parts[1] + "." + parts[0];
  }

  function deadlineClass(daysLeft) {
    if (daysLeft < 0) return "overdue";
    if (daysLeft === 0) return "today";
    if (daysLeft <= 2) return "soon";
    return "later";
  }

  function formatDeadline(dateStr, daysLeft) {
    const parts = dateStr.split("-");
    const formatted = parts[2] + "." + parts[1] + "." + parts[0];
    if (daysLeft < 0) return formatted + " — просрочено на " + -daysLeft + " дн.";
    if (daysLeft === 0) return formatted + " — срок сегодня";
    return formatted + " — осталось " + daysLeft + " дн.";
  }

  function apiFetch(path, options) {
    options = options || {};
    options.headers = Object.assign({}, options.headers, {
      "X-Telegram-Init-Data": initData,
      "Content-Type": "application/json",
    });
    return fetch(path, options).then(function (resp) {
      if (resp.status === 204) return null;
      if (!resp.ok) {
        return resp.text().then(function (text) {
          throw new Error(text || ("HTTP " + resp.status));
        });
      }
      return resp.json();
    });
  }

  function setStatus(text) {
    statusEl.textContent = text;
  }

  function showFormError(text) {
    formErrorEl.textContent = text;
    formErrorEl.classList.remove("hidden");
  }

  function hideFormError() {
    formErrorEl.classList.add("hidden");
    formErrorEl.textContent = "";
  }

  function showSettingsError(text) {
    settingsErrorEl.textContent = text;
    settingsErrorEl.classList.remove("hidden");
  }

  function hideSettingsError() {
    settingsErrorEl.classList.add("hidden");
    settingsErrorEl.textContent = "";
  }

  function haptic(type) {
    if (tg && tg.HapticFeedback && tg.HapticFeedback.notificationOccurred) {
      tg.HapticFeedback.notificationOccurred(type === "notificationSuccess" ? "success" : "error");
    }
  }
})();
