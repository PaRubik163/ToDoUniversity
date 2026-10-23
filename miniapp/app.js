(function () {
  "use strict";

  const translations = {
    ru: {
      appTitle: "Задания", settingsAria: "Настройки", addTaskAria: "Добавить задание",
      settingsTitle: "Настройки", reminderDaysLabel: "Напоминать за", days14: "14 дней", days7: "7 дней", days3: "3 дня", days1: "1 день",
      reminderTimesLabel: "Время напоминаний каждый день", addReminderTime: "Добавить время", removeReminderTime: "Удалить время", languageLabel: "Язык",
      designLabel: "Оформление", designClassic: "Обычное", designNeon: "Неоновое",
      cancel: "Отмена", save: "Сохранить", saveChanges: "Сохранить изменения",
      newTask: "Новое задание", subject: "Предмет", subjectPlaceholder: "Разработка безопасного ПО",
      task: "Задание", taskPlaceholder: "Практическая 1", deadline: "Дедлайн",
      sectionsAria: "Разделы заданий", uncompleted: "Несделанные", completed: "Сделанные",
      overdue: "Просроченные", analytics: "Аналитика", previousMonth: "Предыдущий месяц",
      nextMonth: "Следующий месяц", analyticsMonth: "Месяц аналитики", doneThisWeek: "Сделано за неделю",
      doneThisMonth: "Сделано за месяц", uncompletedTotal: "Несделано всего", overdueTotal: "Просрочено",
      currentWeek: "Текущая неделя", loading: "Загрузка…",
      openTelegram: "Откройте это приложение через Telegram, чтобы увидеть свои задания.",
      fillAll: "Заполните все поля.", timeRequired: "Укажите хотя бы одно время и выберите интервал.",
      loadTasksError: "Не удалось загрузить список: ", loadAnalyticsError: "Не удалось загрузить аналитику: ",
      deleteError: "Не удалось удалить: ", statusError: "Не удалось изменить статус: ",
      deleteConfirm: "Удалить это задание?", noCompleted: "Пока нет сделанных заданий.",
      noOverdue: "Просроченных заданий нет.", noUncompleted: "Несделанных заданий нет. Нажмите «+», чтобы добавить.",
      editTask: "Редактирование задания №{id}", doneDate: "Сделано: {date}", done: "Сделано",
      undoDone: "Вернуть в несделанные", markDone: "Отметить сделанным", overdueBy: "просрочено на {days} {unit}",
      editTaskAria: "Редактировать", deleteTaskAria: "Удалить",
      dueToday: "срок сегодня", daysLeft: "осталось {days} {unit}", dayUnit: "дней", dayUnitOne: "день", dayUnitFew: "дня",
      weekRange: "Неделя: {from} — {to}", monthGraphTitle: "Выполнено по дням — {month}",
      weekChartAria: "Выполненные задания за текущую неделю", monthChartAria: "Выполненные задания за выбранный месяц",
      chartCount: "{date}: {count}"
    },
    en: {
      appTitle: "Tasks", settingsAria: "Settings", addTaskAria: "Add task",
      settingsTitle: "Settings", reminderDaysLabel: "Remind me", days14: "14 days before", days7: "7 days before", days3: "3 days before", days1: "1 day before",
      reminderTimesLabel: "Reminder times each day", addReminderTime: "Add time", removeReminderTime: "Remove time", languageLabel: "Language",
      designLabel: "Appearance", designClassic: "Default", designNeon: "Neon",
      cancel: "Cancel", save: "Save", saveChanges: "Save changes",
      newTask: "New task", subject: "Subject", subjectPlaceholder: "Secure Software Development",
      task: "Task", taskPlaceholder: "Assignment 1", deadline: "Deadline",
      sectionsAria: "Task sections", uncompleted: "To do", completed: "Done",
      overdue: "Overdue", analytics: "Analytics", previousMonth: "Previous month",
      nextMonth: "Next month", analyticsMonth: "Analytics month", doneThisWeek: "Done this week",
      doneThisMonth: "Done this month", uncompletedTotal: "To do", overdueTotal: "Overdue",
      currentWeek: "Current week", loading: "Loading…",
      openTelegram: "Open this app through Telegram to see your tasks.",
      fillAll: "Please fill in every field.", timeRequired: "Choose at least one time and one interval.",
      loadTasksError: "Could not load tasks: ", loadAnalyticsError: "Could not load analytics: ",
      deleteError: "Could not delete task: ", statusError: "Could not update task: ",
      deleteConfirm: "Delete this task?", noCompleted: "No completed tasks yet.",
      noOverdue: "There are no overdue tasks.", noUncompleted: "No tasks to do. Tap “+” to add one.",
      editTask: "Edit task №{id}", doneDate: "Completed: {date}", done: "Completed",
      undoDone: "Move back to to do", markDone: "Mark as done", overdueBy: "overdue by {days} {unit}",
      editTaskAria: "Edit task", deleteTaskAria: "Delete task",
      dueToday: "due today", daysLeft: "{days} {unit} left", dayUnit: "days", dayUnitOne: "day", dayUnitFew: "days",
      weekRange: "Week: {from} — {to}", monthGraphTitle: "Completed by day — {month}",
      weekChartAria: "Tasks completed this week", monthChartAria: "Tasks completed in selected month",
      chartCount: "{date}: {count}"
    }
  };

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
  const taskReminderDaysInput = document.getElementById("task-reminder-days");
  const taskReminderValueEl = document.getElementById("task-reminder-value");

  const settingsToggleEl = document.getElementById("settings-toggle");
  const settingsFormEl = document.getElementById("settings-form");
  const settingsErrorEl = document.getElementById("settings-error");
  const cancelSettingsEl = document.getElementById("cancel-settings");
  const reminderTimesEl = document.getElementById("reminder-times");
  const addReminderTimeEl = document.getElementById("add-reminder-time");
  const themeButtons = Array.from(document.querySelectorAll(".theme-option"));
  const neonThemeStylesheet = document.getElementById("neon-theme-stylesheet");
  const languageButtons = Array.from(document.querySelectorAll(".language-option"));

  let editingId = null;
  let currentSection = "uncompleted";
  let homeworkItems = [];
  let selectedMonth = "";
  let language = "ru";
  let settingsOriginalLanguage = "ru";
  let designTheme = "classic";
  let settingsOriginalTheme = "classic";
  let languageTouched = false;
  let analyticsData = null;

  if (tg) {
    tg.ready();
    tg.expand();
  }

  if (!initData) {
    setStatus(t("openTelegram"));
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
  languageButtons.forEach(function (button) {
    button.addEventListener("click", function () {
      languageTouched = true;
      language = button.dataset.language;
      applyLanguage();
    });
  });
  themeButtons.forEach(function (button) {
    button.addEventListener("click", function () {
      designTheme = button.dataset.theme;
      applyDesignTheme();
    });
  });
  addReminderTimeEl.addEventListener("click", function () {
    const existing = Array.from(reminderTimesEl.querySelectorAll('input[type="time"]'));
    if (existing.length >= 8) return;
    const last = existing.length ? existing[existing.length - 1].value : "09:00";
    const parsed = last ? last.split(":").map(Number) : [9, 0];
    const nextHour = (parsed[0] + 1) % 24;
    addReminderTimeRow(String(nextHour).padStart(2, "0") + ":" + String(parsed[1]).padStart(2, "0"));
  });
  taskReminderDaysInput.addEventListener("input", updateTaskReminderDisplay);

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
      reminder_days: [Number(taskReminderDaysInput.value)],
    };
    if (!payload.subject || !payload.title || !payload.deadline || !payload.reminder_days.length) {
      showFormError(t("fillAll"));
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
    language = settingsOriginalLanguage;
    designTheme = settingsOriginalTheme;
    applyDesignTheme();
    languageTouched = false;
    applyLanguage();
    settingsFormEl.classList.add("hidden");
  });

  settingsFormEl.addEventListener("submit", function (event) {
    event.preventDefault();
    hideSettingsError();
    const reminderTimes = Array.from(reminderTimesEl.querySelectorAll('input[type="time"]')).map(function (input) { return input.value; }).filter(Boolean);
    if (!reminderTimes.length || new Set(reminderTimes).size !== reminderTimes.length) {
      showSettingsError(t("timeRequired"));
      return;
    }
    saveSettings(reminderTimes, language, designTheme);
  });

  function openAddForm() {
    editingId = null;
    formEl.reset();
    setTaskReminderDays(7);
    hideFormError();
    formTitleEl.textContent = t("newTask");
    formEl.querySelector('button[type="submit"]').textContent = t("save");
    formEl.classList.remove("hidden");
    subjectInput.focus();
  }

  function openEditForm(item) {
    editingId = item.id;
    hideFormError();
    subjectInput.value = item.subject;
    titleInput.value = item.title;
    deadlineInput.value = item.deadline;
    setTaskReminderDays(item.reminder_days && item.reminder_days.length ? item.reminder_days[0] : 7);
    formTitleEl.textContent = t("editTask", { id: item.id });
    formEl.querySelector('button[type="submit"]').textContent = t("saveChanges");
    settingsFormEl.classList.add("hidden");
    formEl.classList.remove("hidden");
    subjectInput.focus();
  }

  function closeAddForm() {
    editingId = null;
    formEl.reset();
    setTaskReminderDays(7);
    hideFormError();
    formEl.classList.add("hidden");
  }

  function setTaskReminderDays(days) {
    taskReminderDaysInput.value = String(Math.min(20, Math.max(1, days)));
    updateTaskReminderDisplay();
  }

  function updateTaskReminderDisplay() {
    const days = Number(taskReminderDaysInput.value);
    const unitKey = language === "en" ? "dayUnit" : days === 1 ? "dayUnitOne" : days >= 2 && days <= 4 ? "dayUnitFew" : "dayUnit";
    taskReminderValueEl.textContent = days + " " + t(unitKey);
    taskReminderValueEl.value = taskReminderValueEl.textContent;
    taskReminderDaysInput.style.setProperty("--reminder-progress", ((days - 1) / 19 * 100) + "%");
    formEl.dataset.reminderZone = days <= 5 ? "pink" : days <= 14 ? "orange" : "blue";
  }

  function openSettingsForm() {
    hideSettingsError();
    languageTouched = false;
    settingsOriginalLanguage = language;
    settingsOriginalTheme = designTheme;
    settingsFormEl.classList.remove("hidden");
    reminderTimesEl.replaceChildren();
    apiFetch("/api/settings")
      .then(function (settings) {
        renderReminderTimes((settings && settings.reminder_times) || ["09:00"]);
        designTheme = (settings && settings.design_theme) || "classic";
        settingsOriginalTheme = designTheme;
        applyDesignTheme();
        if (!languageTouched) {
          settingsOriginalLanguage = (settings && settings.language) || language;
          language = settingsOriginalLanguage;
          applyLanguage();
        }
      })
      .catch(function () {
        renderReminderTimes(["09:00"]);
        designTheme = "classic";
        settingsOriginalTheme = designTheme;
        applyDesignTheme();
        settingsOriginalLanguage = language;
      });
  }

  function saveSettings(reminderTimes, selectedLanguage, selectedTheme) {
    const submitButton = settingsFormEl.querySelector('button[type="submit"]');
    submitButton.disabled = true;
    apiFetch("/api/settings", { method: "PUT", body: JSON.stringify({ reminder_times: reminderTimes, language: selectedLanguage, design_theme: selectedTheme }) })
      .then(function () {
        haptic("notificationSuccess");
        settingsOriginalLanguage = selectedLanguage;
        settingsOriginalTheme = selectedTheme;
        languageTouched = false;
        settingsFormEl.classList.add("hidden");
      })
      .catch(function (err) {
        showSettingsError(err.message);
      })
      .finally(function () {
        submitButton.disabled = false;
      });
  }

  function applyDesignTheme() {
    const neon = designTheme === "neon";
    neonThemeStylesheet.disabled = !neon;
    document.body.dataset.theme = designTheme;
    themeButtons.forEach(function (button) {
      const selected = button.dataset.theme === designTheme;
      button.classList.toggle("active", selected);
      button.setAttribute("aria-pressed", selected ? "true" : "false");
    });
  }

  function renderReminderTimes(times) {
    reminderTimesEl.replaceChildren();
    (times.length ? times : ["09:00"]).forEach(addReminderTimeRow);
  }

  function addReminderTimeRow(value) {
    const row = document.createElement("div");
    row.className = "reminder-time-row";
    const input = document.createElement("input");
    input.type = "time";
    input.required = true;
    input.value = value;
    const remove = document.createElement("button");
    remove.type = "button";
    remove.className = "language-option remove-reminder-time";
    remove.textContent = t("removeReminderTime");
    remove.setAttribute("aria-label", t("removeReminderTime"));
    remove.addEventListener("click", function () {
      if (reminderTimesEl.children.length > 1) {
        row.remove();
        updateReminderTimeRemoveButtons();
      }
    });
    row.append(input, remove);
    reminderTimesEl.append(row);
    updateReminderTimeRemoveButtons();
  }

  function updateReminderTimeRemoveButtons() {
    const canRemove = reminderTimesEl.children.length > 1;
    reminderTimesEl.querySelectorAll(".remove-reminder-time").forEach(function (button) { button.disabled = !canRemove; });
  }

  function loadHomework() {
    setStatus(t("loading"));
    Promise.all([apiFetch("/api/homework"), fetchAnalytics(), apiFetch("/api/settings")])
      .then(function (results) {
        homeworkItems = results[0] || [];
        language = results[2].language || "ru";
        designTheme = results[2].design_theme || "classic";
        settingsOriginalTheme = designTheme;
        applyDesignTheme();
        settingsOriginalLanguage = language;
        languageTouched = false;
        applyLanguage();
        renderCurrentSection();
        renderAnalytics(results[1] || {});
        setStatus("");
      })
      .catch(function (err) {
        setStatus(t("loadTasksError") + err.message);
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
          tg.showConfirm(t("deleteConfirm"), resolve);
        })
      : Promise.resolve(window.confirm(t("deleteConfirm")));

    confirmed.then(function (ok) {
      if (!ok) return;
      apiFetch("/api/homework/" + id, { method: "DELETE" })
        .then(function () {
          haptic("notificationSuccess");
          loadHomework();
        })
        .catch(function (err) {
          setStatus(t("deleteError") + err.message);
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
    if (currentSection === "completed") return t("noCompleted");
    if (currentSection === "overdue") return t("noOverdue");
    return t("noUncompleted");
  }

  function renderAnalytics(data) {
    analyticsData = data;
    statCompletedEl.textContent = data.completed_this_week || 0;
    statMonthCompletedEl.textContent = data.completed_this_month || 0;
    statUncompletedEl.textContent = data.uncompleted || 0;
    statOverdueEl.textContent = data.overdue || 0;
    weekRangeEl.textContent = t("weekRange", { from: formatDate(data.week_start), to: formatDate(data.week_end) });
    selectedMonth = data.month;
    monthPickerEl.value = data.month;
    monthPickerEl.max = data.current_month;
    monthNextEl.disabled = data.month >= data.current_month;
    monthChartTitleEl.textContent = t("monthGraphTitle", { month: formatMonth(data.month) });
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
      .catch(function (err) { setStatus(t("loadAnalyticsError") + err.message); });
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
      .toLocaleDateString(language === "en" ? "en-US" : "ru-RU", { month: "long", year: "numeric" });
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
    svg.setAttribute("aria-label", t(showWeekday ? "weekChartAria" : "monthChartAria"));
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
      title.textContent = t("chartCount", { date: formatDate(point.day.date), count: point.day.count });
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
          ? date.toLocaleDateString(language === "en" ? "en-US" : "ru-RU", { weekday: "short" })
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
      completedDate.textContent = item.completed_at ? t("doneDate", { date: formatDate(item.completed_at) }) : t("done");
      body.appendChild(completedDate);
    }

    const actions = document.createElement("div");
    actions.className = "item-actions";

    const completeButton = document.createElement("button");
    completeButton.className = "icon-action complete-button" + (item.completed ? " is-completed" : "");
    completeButton.type = "button";
    completeButton.setAttribute("aria-label", item.completed ? t("undoDone") : t("markDone"));
    completeButton.textContent = item.completed ? "↶" : "✓";
    completeButton.addEventListener("click", function () {
      apiFetch("/api/homework/" + item.id + "/status", {
        method: "PUT",
        body: JSON.stringify({ completed: !item.completed }),
      }).then(function () {
        haptic("notificationSuccess");
        loadHomework();
      }).catch(function (err) {
        setStatus(t("statusError") + err.message);
      });
    });

    const editButton = document.createElement("button");
    editButton.className = "icon-action edit-button";
    editButton.type = "button";
    editButton.setAttribute("aria-label", t("editTaskAria"));
    editButton.textContent = "✎";
    editButton.addEventListener("click", function () {
      openEditForm(item);
    });

    const deleteButton = document.createElement("button");
    deleteButton.className = "icon-action delete-button";
    deleteButton.type = "button";
    deleteButton.setAttribute("aria-label", t("deleteTaskAria"));
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
    return language === "en"
      ? parts[1] + "/" + parts[2] + "/" + parts[0]
      : parts[2] + "." + parts[1] + "." + parts[0];
  }

  function deadlineClass(daysLeft) {
    if (daysLeft < 0) return "overdue";
    if (daysLeft === 0) return "today";
    if (daysLeft <= 2) return "soon";
    return "later";
  }

  function formatDeadline(dateStr, daysLeft) {
    const formatted = formatDate(dateStr);
    if (daysLeft < 0) {
      const overdueDays = -daysLeft;
      return formatted + " — " + t("overdueBy", { days: overdueDays, unit: dayUnit(overdueDays) });
    }
    if (daysLeft === 0) return formatted + " — " + t("dueToday");
    return formatted + " — " + t("daysLeft", { days: daysLeft, unit: dayUnit(daysLeft) });
  }

  function dayUnit(days) {
    if (language === "en") return days === 1 ? t("dayUnitOne") : t("dayUnit");
    const remainder = days % 10;
    const lastTwo = days % 100;
    if (remainder === 1 && lastTwo !== 11) return t("dayUnitOne");
    if (remainder >= 2 && remainder <= 4 && (lastTwo < 12 || lastTwo > 14)) return t("dayUnitFew");
    return t("dayUnit");
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

  function t(key, values) {
    let value = (translations[language] && translations[language][key]) || translations.ru[key] || key;
    Object.keys(values || {}).forEach(function (name) {
      value = value.replace("{" + name + "}", values[name]);
    });
    return value;
  }

  function applyLanguage() {
    document.documentElement.lang = language;
    updateTaskReminderDisplay();
    document.title = t("appTitle");
    document.querySelectorAll("[data-i18n]").forEach(function (element) {
      element.textContent = t(element.dataset.i18n);
    });
    document.querySelectorAll("[data-i18n-aria]").forEach(function (element) {
      element.setAttribute("aria-label", t(element.dataset.i18nAria));
    });
    document.querySelectorAll("[data-i18n-placeholder]").forEach(function (element) {
      element.setAttribute("placeholder", t(element.dataset.i18nPlaceholder));
    });
    languageButtons.forEach(function (button) {
      const selected = button.dataset.language === language;
      button.classList.toggle("active", selected);
      button.setAttribute("aria-pressed", selected ? "true" : "false");
    });
    document.querySelectorAll(".remove-reminder-time").forEach(function (button) {
      button.textContent = t("removeReminderTime");
      button.setAttribute("aria-label", t("removeReminderTime"));
    });
    if (editingId) {
      formTitleEl.textContent = t("editTask", { id: editingId });
      formEl.querySelector('button[type="submit"]').textContent = t("saveChanges");
    }
    renderCurrentSection();
    if (analyticsData) renderAnalytics(analyticsData);
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
