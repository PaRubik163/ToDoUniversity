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
    renderLineChart(weekChartEl, data.week_days || [], true);
    renderLineChart(monthChartEl, data.month_days || [], false);
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

  function renderLineChart(container, days, showWeekday) {
    const namespace = "http://www.w3.org/2000/svg";
    const width = showWeekday ? 520 : Math.max(520, days.length * 24);
    const height = 210;
    const left = 38;
    const right = width - 14;
    const top = 18;
    const bottom = height - 34;
    const max = Math.max(1, ...days.map(function (day) { return day.count; }));
    const svg = document.createElementNS(namespace, "svg");
    svg.setAttribute("viewBox", "0 0 " + width + " " + height);
    svg.setAttribute("role", "img");
    svg.setAttribute("aria-label", showWeekday ? "Выполненные задачи за текущую неделю" : "Выполненные задачи за выбранный месяц");
    if (showWeekday) {
      svg.setAttribute("preserveAspectRatio", "none");
    } else {
      svg.style.width = width + "px";
    }

    (max < 2 ? [0, max] : [0, max / 2, max]).forEach(function (value) {
      const y = bottom - value / max * (bottom - top);
      const gridLine = document.createElementNS(namespace, "line");
      gridLine.setAttribute("x1", left);
      gridLine.setAttribute("x2", right);
      gridLine.setAttribute("y1", y);
      gridLine.setAttribute("y2", y);
      gridLine.setAttribute("class", "chart-grid-line");
      svg.appendChild(gridLine);

      const axisLabel = document.createElementNS(namespace, "text");
      axisLabel.setAttribute("x", left - 8);
      axisLabel.setAttribute("y", y + 4);
      axisLabel.setAttribute("text-anchor", "end");
      axisLabel.setAttribute("class", "chart-axis-label");
      axisLabel.textContent = Math.round(value);
      svg.appendChild(axisLabel);
    });

    const points = days.map(function (day, index) {
      const x = days.length <= 1 ? (left + right) / 2 : left + index * (right - left) / (days.length - 1);
      const y = bottom - day.count / max * (bottom - top);
      return { day: day, x: x, y: y };
    });
    const line = document.createElementNS(namespace, "polyline");
    line.setAttribute("points", points.map(function (point) { return point.x + "," + point.y; }).join(" "));
    line.setAttribute("class", "chart-line");
    svg.appendChild(line);

    points.forEach(function (point, index) {
      const circle = document.createElementNS(namespace, "circle");
      circle.setAttribute("cx", point.x);
      circle.setAttribute("cy", point.y);
      circle.setAttribute("r", showWeekday ? 4 : 3);
      circle.setAttribute("class", "chart-point");
      const title = document.createElementNS(namespace, "title");
      title.textContent = formatDate(point.day.date) + ": " + point.day.count;
      circle.appendChild(title);
      svg.appendChild(circle);

      const dateParts = point.day.date.split("-").map(Number);
      const date = new Date(dateParts[0], dateParts[1] - 1, dateParts[2]);
      const showMonthLabel = index === 0 || index === points.length - 1 || date.getDate() % 5 === 0;
      if (showWeekday || showMonthLabel) {
        const label = document.createElementNS(namespace, "text");
        label.setAttribute("x", point.x);
        label.setAttribute("y", height - 9);
        label.setAttribute("text-anchor", "middle");
        label.setAttribute("class", "chart-date-label");
        label.textContent = showWeekday
          ? date.toLocaleDateString("ru-RU", { weekday: "short" })
          : String(date.getDate());
        svg.appendChild(label);
      }
    });

    container.replaceChildren(svg);
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
