(function () {
  "use strict";

  const tg = window.Telegram && window.Telegram.WebApp;
  const initData = tg ? tg.initData : "";

  const listEl = document.getElementById("list");
  const emptyEl = document.getElementById("empty-state");
  const statusEl = document.getElementById("status");
  const formEl = document.getElementById("add-form");
  const formErrorEl = document.getElementById("form-error");
  const addToggleEl = document.getElementById("add-toggle");
  const cancelAddEl = document.getElementById("cancel-add");
  const subjectInput = document.getElementById("field-subject");
  const titleInput = document.getElementById("field-title");
  const deadlineInput = document.getElementById("field-deadline");

  if (tg) {
    tg.ready();
    tg.expand();
  }

  if (!initData) {
    setStatus("Откройте это приложение через Telegram, чтобы увидеть свои задания.");
  } else {
    loadHomework();
  }

  addToggleEl.addEventListener("click", function () {
    formEl.classList.toggle("hidden");
    if (!formEl.classList.contains("hidden")) {
      subjectInput.focus();
    }
  });

  cancelAddEl.addEventListener("click", function () {
    formEl.reset();
    hideFormError();
    formEl.classList.add("hidden");
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

    createHomework(payload);
  });

  function loadHomework() {
    setStatus("Загрузка…");
    apiFetch("/api/homework")
      .then(function (items) {
        renderList(items || []);
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
        formEl.reset();
        formEl.classList.add("hidden");
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
          itemEl.remove();
          if (!listEl.children.length) {
            emptyEl.classList.remove("hidden");
          }
        })
        .catch(function (err) {
          setStatus("Не удалось удалить: " + err.message);
        });
    });
  }

  function renderList(items) {
    listEl.innerHTML = "";
    if (!items.length) {
      emptyEl.classList.remove("hidden");
      return;
    }
    emptyEl.classList.add("hidden");

    items.forEach(function (item) {
      listEl.appendChild(renderItem(item));
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

    const deleteButton = document.createElement("button");
    deleteButton.className = "delete-button";
    deleteButton.type = "button";
    deleteButton.setAttribute("aria-label", "Удалить");
    deleteButton.textContent = "✕";
    deleteButton.addEventListener("click", function () {
      deleteHomework(item.id, li);
    });

    li.appendChild(body);
    li.appendChild(deleteButton);
    return li;
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

  function haptic(type) {
    if (tg && tg.HapticFeedback && tg.HapticFeedback.notificationOccurred) {
      tg.HapticFeedback.notificationOccurred(type === "notificationSuccess" ? "success" : "error");
    }
  }
})();
