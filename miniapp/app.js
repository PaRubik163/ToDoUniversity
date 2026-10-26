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
      chartCount: "{date}: {count}", eventHourUnit: "\u0447\u0430\u0441.",
      calendar: "\u041f\u043e\u0441\u0442\u043e\u044f\u043d\u043d\u044b\u0435", newEvent: "\u041d\u043e\u0432\u043e\u0435 \u0441\u043e\u0431\u044b\u0442\u0438\u0435", eventTitle: "\u0421\u043e\u0431\u044b\u0442\u0438\u0435", eventDateTime: "\u0414\u0430\u0442\u0430 \u0438 \u0432\u0440\u0435\u043c\u044f", remindHours: "\u041d\u0430\u043f\u043e\u043c\u043d\u0438\u0442\u044c \u0437\u0430 (\u0447\u0430\u0441\u043e\u0432)", eventsForDay: "\u0421\u043e\u0431\u044b\u0442\u0438\u044f \u0434\u043d\u044f", addEvent: "+ \u0421\u043e\u0431\u044b\u0442\u0438\u0435", noEvents: "\u041d\u0430 \u044d\u0442\u0443 \u0434\u0430\u0442\u0443 \u0441\u043e\u0431\u044b\u0442\u0438\u0439 \u043d\u0435\u0442.", eventDelete: "\u0423\u0434\u0430\u043b\u0438\u0442\u044c \u044d\u0442\u043e \u0441\u043e\u0431\u044b\u0442\u0438\u0435?", eventError: "\u041d\u0435 \u0443\u0434\u0430\u043b\u043e\u0441\u044c \u0441\u043e\u0445\u0440\u0430\u043d\u0438\u0442\u044c: ", eventDone: "\u0413\u043e\u0442\u043e\u0432\u043e", eventUndo: "\u0412\u0435\u0440\u043d\u0443\u0442\u044c", noReminder: "\u0411\u0435\u0437 \u043d\u0430\u043f\u043e\u043c\u0438\u043d\u0430\u043d\u0438\u044f"
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
      chartCount: "{date}: {count}", eventHourUnit: "h", calendar: "Calendar", newEvent: "New event", eventTitle: "Event", eventDateTime: "Date and time", remindHours: "Remind me (hours before)", eventsForDay: "Events for this day", addEvent: "+ Event", noEvents: "No events on this date.", eventDelete: "Delete this event?", eventError: "Could not save event: ", eventDone: "Done", eventUndo: "Undo", noReminder: "No reminder"
    },
    zh: {
      appTitle: "任务", settingsAria: "设置", addTaskAria: "添加任务", settingsTitle: "设置", reminderDaysLabel: "提前提醒",
      days14: "提前14天", days7: "提前7天", days3: "提前3天", days1: "提前1天", reminderTimesLabel: "每日提醒时间",
      addReminderTime: "添加时间", removeReminderTime: "删除时间", languageLabel: "语言", designLabel: "外观", designClassic: "经典", designNeon: "霓虹",
      cancel: "取消", save: "保存", saveChanges: "保存更改", newTask: "新任务", subject: "课程", subjectPlaceholder: "软件开发",
      task: "任务", taskPlaceholder: "作业1", deadline: "截止日期", sectionsAria: "任务分类", uncompleted: "待办", completed: "已完成",
      overdue: "已逾期", analytics: "统计", previousMonth: "上个月", nextMonth: "下个月", analyticsMonth: "统计月份",
      doneThisWeek: "本周完成", doneThisMonth: "本月完成", uncompletedTotal: "待办总数", overdueTotal: "逾期总数", currentWeek: "本周",
      loading: "加载中…", openTelegram: "请通过 Telegram 打开应用以查看任务。", fillAll: "请填写所有字段。", timeRequired: "请至少选择一个提醒时间。",
      loadTasksError: "无法加载任务：", loadAnalyticsError: "无法加载统计：", deleteError: "无法删除：", statusError: "无法更新状态：",
      deleteConfirm: "确定删除此任务吗？", noCompleted: "还没有已完成的任务。", noOverdue: "没有逾期任务。", noUncompleted: "没有待办任务。点击“+”添加。",
      editTask: "编辑任务 №{id}", doneDate: "完成于：{date}", done: "已完成", undoDone: "移回待办", markDone: "标记为完成",
      overdueBy: "逾期 {days}{unit}", editTaskAria: "编辑任务", deleteTaskAria: "删除任务", dueToday: "今天截止", daysLeft: "还剩 {days}{unit}",
      dayUnit: "天", dayUnitOne: "天", dayUnitFew: "天", weekRange: "本周：{from} — {to}", monthGraphTitle: "每日完成情况 — {month}",
      weekChartAria: "本周完成的任务", monthChartAria: "所选月份完成的任务", chartCount: "{date}：{count}"
    },
    es: {
      appTitle: "Tareas", settingsAria: "Ajustes", addTaskAria: "Añadir tarea", settingsTitle: "Ajustes", reminderDaysLabel: "Avisar con",
      days14: "14 días antes", days7: "7 días antes", days3: "3 días antes", days1: "1 día antes", reminderTimesLabel: "Hora diaria de los avisos",
      addReminderTime: "Añadir hora", removeReminderTime: "Quitar hora", languageLabel: "Idioma", designLabel: "Apariencia", designClassic: "Clásico", designNeon: "Neón",
      cancel: "Cancelar", save: "Guardar", saveChanges: "Guardar cambios", newTask: "Nueva tarea", subject: "Asignatura", subjectPlaceholder: "Desarrollo de software",
      task: "Tarea", taskPlaceholder: "Práctica 1", deadline: "Fecha límite", sectionsAria: "Secciones de tareas", uncompleted: "Pendientes", completed: "Completadas",
      overdue: "Atrasadas", analytics: "Estadísticas", previousMonth: "Mes anterior", nextMonth: "Mes siguiente", analyticsMonth: "Mes de estadísticas",
      doneThisWeek: "Completadas esta semana", doneThisMonth: "Completadas este mes", uncompletedTotal: "Pendientes", overdueTotal: "Atrasadas", currentWeek: "Esta semana",
      loading: "Cargando…", openTelegram: "Abre esta aplicación desde Telegram para ver tus tareas.", fillAll: "Completa todos los campos.", timeRequired: "Elige al menos una hora de aviso.",
      loadTasksError: "No se pudieron cargar las tareas: ", loadAnalyticsError: "No se pudieron cargar las estadísticas: ", deleteError: "No se pudo eliminar: ", statusError: "No se pudo actualizar el estado: ",
      deleteConfirm: "¿Eliminar esta tarea?", noCompleted: "Aún no hay tareas completadas.", noOverdue: "No hay tareas atrasadas.", noUncompleted: "No hay tareas pendientes. Pulsa «+» para añadir una.",
      editTask: "Editar tarea n.º {id}", doneDate: "Completada: {date}", done: "Completada", undoDone: "Volver a pendientes", markDone: "Marcar como completada",
      overdueBy: "con {days} {unit} de retraso", editTaskAria: "Editar tarea", deleteTaskAria: "Eliminar tarea", dueToday: "vence hoy", daysLeft: "quedan {days} {unit}",
      dayUnit: "días", dayUnitOne: "día", dayUnitFew: "días", weekRange: "Semana: {from} — {to}", monthGraphTitle: "Tareas completadas por día — {month}",
      weekChartAria: "Tareas completadas esta semana", monthChartAria: "Tareas completadas en el mes seleccionado", chartCount: "{date}: {count}"
    },
    ar: {
      appTitle: "المهام", settingsAria: "الإعدادات", addTaskAria: "إضافة مهمة", settingsTitle: "الإعدادات", reminderDaysLabel: "التذكير قبل",
      days14: "قبل 14 يومًا", days7: "قبل 7 أيام", days3: "قبل 3 أيام", days1: "قبل يوم", reminderTimesLabel: "وقت التذكير اليومي",
      addReminderTime: "إضافة وقت", removeReminderTime: "حذف الوقت", languageLabel: "اللغة", designLabel: "المظهر", designClassic: "كلاسيكي", designNeon: "نيون",
      cancel: "إلغاء", save: "حفظ", saveChanges: "حفظ التغييرات", newTask: "مهمة جديدة", subject: "المادة", subjectPlaceholder: "تطوير البرمجيات",
      task: "المهمة", taskPlaceholder: "واجب 1", deadline: "الموعد النهائي", sectionsAria: "أقسام المهام", uncompleted: "قيد الإنجاز", completed: "مكتملة",
      overdue: "متأخرة", analytics: "الإحصاءات", previousMonth: "الشهر السابق", nextMonth: "الشهر التالي", analyticsMonth: "شهر الإحصاءات",
      doneThisWeek: "المكتمل هذا الأسبوع", doneThisMonth: "المكتمل هذا الشهر", uncompletedTotal: "المهام المتبقية", overdueTotal: "المهام المتأخرة", currentWeek: "هذا الأسبوع",
      loading: "جارٍ التحميل…", openTelegram: "افتح التطبيق من Telegram لعرض مهامك.", fillAll: "يرجى ملء جميع الحقول.", timeRequired: "اختر وقت تذكير واحدًا على الأقل.",
      loadTasksError: "تعذر تحميل المهام: ", loadAnalyticsError: "تعذر تحميل الإحصاءات: ", deleteError: "تعذر الحذف: ", statusError: "تعذر تحديث الحالة: ",
      deleteConfirm: "هل تريد حذف هذه المهمة؟", noCompleted: "لا توجد مهام مكتملة بعد.", noOverdue: "لا توجد مهام متأخرة.", noUncompleted: "لا توجد مهام قيد الإنجاز. اضغط «+» لإضافة مهمة.",
      editTask: "تعديل المهمة رقم {id}", doneDate: "اكتملت: {date}", done: "مكتملة", undoDone: "إعادتها إلى المهام الجارية", markDone: "وضع علامة مكتملة",
      overdueBy: "متأخرة {days} {unit}", editTaskAria: "تعديل المهمة", deleteTaskAria: "حذف المهمة", dueToday: "موعدها اليوم", daysLeft: "متبقي {days} {unit}",
      dayUnit: "أيام", dayUnitOne: "يوم", dayUnitFew: "أيام", weekRange: "الأسبوع: {from} — {to}", monthGraphTitle: "المهام المكتملة يوميًا — {month}",
      weekChartAria: "المهام المكتملة هذا الأسبوع", monthChartAria: "المهام المكتملة في الشهر المحدد", chartCount: "{date}: {count}"
    },
    hy: {
      appTitle: "Առաջադրանքներ", settingsAria: "Կարգավորումներ", addTaskAria: "Ավելացնել առաջադրանք", settingsTitle: "Կարգավորումներ", reminderDaysLabel: "Հիշեցնել առաջ",
      days14: "14 օր առաջ", days7: "7 օր առաջ", days3: "3 օր առաջ", days1: "1 օր առաջ", reminderTimesLabel: "Ամենօրյա հիշեցման ժամերը",
      addReminderTime: "Ավելացնել ժամ", removeReminderTime: "Հեռացնել ժամը", languageLabel: "Լեզու", designLabel: "Տեսք", designClassic: "Դասական", designNeon: "Նեոնային",
      cancel: "Չեղարկել", save: "Պահպանել", saveChanges: "Պահպանել փոփոխությունները", newTask: "Նոր առաջադրանք", subject: "Առարկա", subjectPlaceholder: "Ծրագրային ապահովման մշակում",
      task: "Առաջադրանք", taskPlaceholder: "Գործնական աշխատանք 1", deadline: "Վերջնաժամկետ", sectionsAria: "Առաջադրանքների բաժիններ", uncompleted: "Չկատարված", completed: "Կատարված",
      overdue: "Ժամկետանց", analytics: "Վիճակագրություն", previousMonth: "Նախորդ ամիս", nextMonth: "Հաջորդ ամիս", analyticsMonth: "Վիճակագրության ամիս",
      doneThisWeek: "Կատարված այս շաբաթ", doneThisMonth: "Կատարված այս ամիս", uncompletedTotal: "Չկատարված ընդհանուր", overdueTotal: "Ժամկետանց ընդհանուր", currentWeek: "Այս շաբաթ",
      loading: "Բեռնվում է…", openTelegram: "Առաջադրանքները տեսնելու համար բացեք հավելվածը Telegram-ում։", fillAll: "Լրացրեք բոլոր դաշտերը։", timeRequired: "Ընտրեք առնվազն մեկ հիշեցման ժամ։",
      loadTasksError: "Չհաջողվեց բեռնել առաջադրանքները․ ", loadAnalyticsError: "Չհաջողվեց բեռնել վիճակագրությունը․ ", deleteError: "Չհաջողվեց ջնջել․ ", statusError: "Չհաջողվեց փոխել կարգավիճակը․ ",
      deleteConfirm: "Ջնջե՞լ այս առաջադրանքը։", noCompleted: "Կատարված առաջադրանքներ դեռ չկան։", noOverdue: "Ժամկետանց առաջադրանքներ չկան։", noUncompleted: "Չկատարված առաջադրանքներ չկան։ Սեղմեք «+»՝ ավելացնելու համար։",
      editTask: "Խմբագրել առաջադրանք №{id}", doneDate: "Կատարված է՝ {date}", done: "Կատարված է", undoDone: "Վերադարձնել չկատարվածների մեջ", markDone: "Նշել որպես կատարված",
      overdueBy: "ուշացել է {days} {unit}-ով", editTaskAria: "Խմբագրել առաջադրանքը", deleteTaskAria: "Ջնջել առաջադրանքը", dueToday: "ժամկետը այսօր է", daysLeft: "մնացել է {days} {unit}",
      dayUnit: "օր", dayUnitOne: "օր", dayUnitFew: "օր", weekRange: "Շաբաթ՝ {from} — {to}", monthGraphTitle: "Օրվա ընթացքում կատարված — {month}",
      weekChartAria: "Այս շաբաթ կատարված առաջադրանքները", monthChartAria: "Ընտրված ամսում կատարված առաջադրանքները", chartCount: "{date}՝ {count}"
    }
  };
  const extraCalendarTranslations = {
    zh: { calendar: "日历", newEvent: "新事件", eventTitle: "事件", eventDateTime: "日期和时间", remindHours: "提前提醒（小时）", eventsForDay: "当天事件", addEvent: "+ 事件", noEvents: "这一天没有事件。", eventDelete: "删除此事件？", eventError: "保存失败：", eventDone: "完成", eventUndo: "撤销", noReminder: "不提醒", eventHourUnit: "小时" },
    es: { calendar: "Calendario", newEvent: "Nuevo evento", eventTitle: "Evento", eventDateTime: "Fecha y hora", remindHours: "Avisar con (horas de antelación)", eventsForDay: "Eventos del día", addEvent: "+ Evento", noEvents: "No hay eventos en esta fecha.", eventDelete: "¿Eliminar este evento?", eventError: "No se pudo guardar: ", eventDone: "Hecho", eventUndo: "Deshacer", noReminder: "Sin aviso", eventHourUnit: "h" },
    ar: { calendar: "التقويم", newEvent: "حدث جديد", eventTitle: "الحدث", eventDateTime: "التاريخ والوقت", remindHours: "التذكير قبل (ساعات)", eventsForDay: "أحداث اليوم", addEvent: "+ حدث", noEvents: "لا توجد أحداث في هذا التاريخ.", eventDelete: "هل تريد حذف هذا الحدث؟", eventError: "تعذر الحفظ: ", eventDone: "تم", eventUndo: "تراجع", noReminder: "بدون تذكير", eventHourUnit: "س" },
    hy: { calendar: "Օրացույց", newEvent: "Նոր իրադարձություն", eventTitle: "Իրադարձություն", eventDateTime: "Ամսաթիվ և ժամ", remindHours: "Հիշեցնել (ժամ առաջ)", eventsForDay: "Օրվա իրադարձությունները", addEvent: "+ Իրադարձություն", noEvents: "Այս օրը իրադարձություններ չկան։", eventDelete: "Ջնջե՞լ իրադարձությունը։", eventError: "Չհաջողվեց պահպանել․ ", eventDone: "Կատարված է", eventUndo: "Հետարկել", noReminder: "Առանց հիշեցման", eventHourUnit: "ժ" }
  };
  Object.keys(extraCalendarTranslations).forEach(function (locale) {
    Object.assign(translations[locale], extraCalendarTranslations[locale]);
  });

  const tg = window.Telegram && window.Telegram.WebApp;
  const initData = tg ? tg.initData : "";

  const listEl = document.getElementById("list");
  const emptyEl = document.getElementById("empty-state");
  const statusEl = document.getElementById("status");
  const sectionButtons = Array.from(document.querySelectorAll(".section-button"));
  const analyticsEl = document.getElementById("analytics");
  const calendarEl = document.getElementById("calendar");
  const calendarGridEl = document.getElementById("calendar-grid");
  const calendarMonthEl = document.getElementById("calendar-month");
  const calendarDayTitleEl = document.getElementById("calendar-day-title");
  const eventListEl = document.getElementById("event-list");
  const eventEmptyEl = document.getElementById("event-empty");
  const eventFormEl = document.getElementById("event-form");
  const eventTitleInput = document.getElementById("event-title");
  const eventDateInput = document.getElementById("event-date-time");
  const eventReminderInput = document.getElementById("event-reminder-hours");
  const eventErrorEl = document.getElementById("event-error");
  const weekRangeEl = document.getElementById("week-range");
  const statCompletedEl = document.getElementById("stat-completed");
  const statMonthCompletedEl = document.getElementById("stat-month-completed");
  const statUncompletedEl = document.getElementById("stat-uncompleted");
  const statOverdueEl = document.getElementById("stat-overdue");
  const monthPickerEl = document.getElementById("month-picker");
  const monthPickerLabelEl = document.getElementById("month-picker-label");
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
  const deadlineDisplayEl = document.getElementById("deadline-display");
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
  const languageSelectEl = document.getElementById("language-select");

  let editingId = null;
  let currentSection = "uncompleted";
  let homeworkItems = [];
  let selectedMonth = "";
  let language = telegramLanguage();
  let settingsOriginalLanguage = language;
  let designTheme = "classic";
  let settingsOriginalTheme = "classic";
  let languageTouched = false;
  let analyticsData = null;
  let calendarMonth = new Date().getFullYear() + "-" + String(new Date().getMonth() + 1).padStart(2, "0");
  let selectedCalendarDate = new Date().getFullYear() + "-" + String(new Date().getMonth() + 1).padStart(2, "0") + "-" + String(new Date().getDate()).padStart(2, "0");
  let calendarEvents = [];

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
      if (currentSection !== "calendar") { eventFormEl.classList.add("hidden"); formEl.classList.add("hidden"); settingsFormEl.classList.add("hidden"); }
      renderCurrentSection();
    });
  });

  monthPickerEl.addEventListener("change", function () {
    if (!monthPickerEl.value) return;
    selectedMonth = monthPickerEl.value;
    updateMonthPickerLabel();
    loadAnalytics();
  });
  monthPrevEl.addEventListener("click", function () { shiftSelectedMonth(-1); });
  monthNextEl.addEventListener("click", function () { shiftSelectedMonth(1); });
  document.getElementById("calendar-prev").addEventListener("click", function () { shiftCalendarMonth(-1); });
  document.getElementById("calendar-next").addEventListener("click", function () { shiftCalendarMonth(1); });
  document.getElementById("calendar-add-event").addEventListener("click", openEventForm);
  document.getElementById("cancel-event").addEventListener("click", function () { eventFormEl.classList.add("hidden"); });
  eventFormEl.addEventListener("submit", submitEvent);
  languageSelectEl.addEventListener("change", function () {
    languageTouched = true;
    language = languageSelectEl.value;
    applyLanguage();
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
  deadlineInput.addEventListener("change", updateDeadlineDisplay);

  addToggleEl.addEventListener("click", function () {
    if (currentSection === "calendar") { openEventForm(); return; }
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
    updateDeadlineDisplay();
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
    updateDeadlineDisplay();
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
    updateDeadlineDisplay();
    setTaskReminderDays(7);
    hideFormError();
    formEl.classList.add("hidden");
  }

  function setTaskReminderDays(days) {
    taskReminderDaysInput.value = String(Math.min(20, Math.max(1, days)));
    updateTaskReminderDisplay();
  }

  function updateDeadlineDisplay() {
    const placeholders = {
      ru: "ДД.ММ.ГГГГ", en: "MM/DD/YYYY", zh: "YYYY/MM/DD", es: "DD/MM/AAAA",
      ar: "\u064a\u0648\u0645/\u0634\u0647\u0631/\u0633\u0646\u0629", hy: "\u0555\u0555.\u0531\u0531.\u054f\u054f\u054f\u054f"
    };
    deadlineDisplayEl.textContent = deadlineInput.value ? formatDate(deadlineInput.value) : placeholders[language] || placeholders.ru;
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
    Promise.all([apiFetch("/api/homework"), fetchAnalytics(), apiFetch("/api/settings"), fetchCalendarEvents()])
      .then(function (results) {
        homeworkItems = results[0] || [];
        language = results[2].language || telegramLanguage();
        designTheme = results[2].design_theme || "classic";
        settingsOriginalTheme = designTheme;
        applyDesignTheme();
        settingsOriginalLanguage = language;
        languageTouched = false;
        applyLanguage();
        renderCurrentSection();
        renderAnalytics(results[1] || {});
        calendarEvents = results[3] || [];
        renderCalendar();
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
    const calendarMode = currentSection === "calendar";
    addToggleEl.setAttribute("aria-label", calendarMode ? t("addEvent") : t("addTaskAria"));
    analyticsEl.classList.toggle("hidden", !analyticsMode);
    calendarEl.classList.toggle("hidden", !calendarMode);
    listEl.classList.toggle("hidden", analyticsMode || calendarMode);
    if (analyticsMode) {
      emptyEl.classList.add("hidden");
      return;
    }
    if (calendarMode) { emptyEl.classList.add("hidden"); renderCalendar(); return; }
    const items = homeworkItems.filter(function (item) {
      if (currentSection === "completed") return item.completed;
      if (currentSection === "overdue") return !item.completed && item.days_left < 0;
      return !item.completed;
    });
    renderList(items);
  }

  function shiftCalendarMonth(amount) {
    const p = calendarMonth.split("-");
    const d = new Date(Number(p[0]), Number(p[1]) - 1 + amount, 1);
    calendarMonth = d.getFullYear() + "-" + String(d.getMonth() + 1).padStart(2, "0");
    selectedCalendarDate = calendarMonth + "-01";
    fetchCalendarEvents().then(function (items) { calendarEvents = items || []; renderCalendar(); })
      .catch(function (err) { setStatus(t("eventError") + err.message); });
  }

  function fetchCalendarEvents() { return apiFetch("/api/events?month=" + encodeURIComponent(calendarMonth)); }

  function renderCalendar() {
    const p = calendarMonth.split("-");
    const year = Number(p[0]), month = Number(p[1]) - 1;
    calendarMonthEl.textContent = new Date(year, month, 1).toLocaleDateString(languageLocale(), { month: "long", year: "numeric" });
    calendarGridEl.replaceChildren();
    const weekdayNames = [];
    for (let i = 0; i < 7; i++) weekdayNames.push(new Date(2024, 0, 1 + i).toLocaleDateString(languageLocale(), { weekday: "short" }));
    weekdayNames.forEach(function (name) { const label = document.createElement("span"); label.className = "calendar-weekday"; label.textContent = name; calendarGridEl.appendChild(label); });
    const firstOffset = (new Date(year, month, 1).getDay() + 6) % 7;
    const days = new Date(year, month + 1, 0).getDate();
    for (let blank = 0; blank < firstOffset; blank++) { const spacer = document.createElement("span"); spacer.className = "calendar-blank"; calendarGridEl.appendChild(spacer); }
    for (let day = 1; day <= days; day++) {
      const key = calendarMonth + "-" + String(day).padStart(2, "0");
      const button = document.createElement("button");
      button.type = "button"; button.className = "calendar-date" + (key === selectedCalendarDate ? " selected" : "");
      const count = calendarEvents.filter(function (item) { return eventDateKey(item) === key; }).length;
      button.innerHTML = "<span>" + day + "</span>" + (count ? '<i class="calendar-dot" aria-label="' + count + '"></i>' : "");
      button.addEventListener("click", function () { selectedCalendarDate = key; renderCalendar(); });
      calendarGridEl.appendChild(button);
    }
    calendarDayTitleEl.textContent = formatDate(selectedCalendarDate);
    const dayItems = calendarEvents.filter(function (item) { return eventDateKey(item) === selectedCalendarDate; });
    eventListEl.replaceChildren();
    eventEmptyEl.textContent = t("noEvents");
    eventEmptyEl.classList.toggle("hidden", dayItems.length > 0);
    dayItems.forEach(function (item) { eventListEl.appendChild(renderCalendarEvent(item)); });
  }

  function eventDateKey(item) {
    return item.event_at.slice(0, 10);
  }

  function renderCalendarEvent(item) {
    const li = document.createElement("li"); li.className = "item calendar-event" + (item.completed ? " is-done" : "");
    const body = document.createElement("div"); body.className = "item-body";
    const title = document.createElement("p"); title.className = "item-title"; title.textContent = item.title;
    const details = document.createElement("p"); details.className = "item-deadline later";
    details.textContent = item.event_at.slice(11, 16) + " · " + (item.reminder_hours ? item.reminder_hours + " " + t("eventHourUnit") : t("noReminder"));
    body.append(title, details);
    if (item.completed && item.completed_at) { const done = document.createElement("p"); done.className = "item-completed"; done.textContent = t("doneDate", { date: formatDate(item.completed_at.slice(0, 10)) }); body.appendChild(done); }
    const actions = document.createElement("div"); actions.className = "item-actions";
    const complete = document.createElement("button"); complete.type = "button"; complete.className = "icon-action complete-button" + (item.completed ? " is-completed" : ""); complete.textContent = item.completed ? "↶" : "✓"; complete.setAttribute("aria-label", item.completed ? t("eventUndo") : t("eventDone"));
    complete.addEventListener("click", function () { apiFetch("/api/events/" + item.id + "/status", { method: "PUT", body: JSON.stringify({ completed: !item.completed }) }).then(refreshEvents).then(loadAnalytics).catch(function (err) { setStatus(t("statusError") + err.message); }); });
    const remove = document.createElement("button"); remove.type = "button"; remove.className = "icon-action delete-button"; remove.textContent = "×"; remove.setAttribute("aria-label", t("deleteTaskAria"));
    remove.addEventListener("click", function () { const ok = tg && tg.showConfirm ? new Promise(function (resolve) { tg.showConfirm(t("eventDelete"), resolve); }) : Promise.resolve(window.confirm(t("eventDelete"))); ok.then(function (yes) { if (yes) apiFetch("/api/events/" + item.id, { method: "DELETE" }).then(refreshEvents).then(loadAnalytics).catch(function (err) { setStatus(t("eventError") + err.message); }); }); });
    actions.append(complete, remove); li.append(body, actions); return li;
  }

  function openEventForm() {
    formEl.classList.add("hidden"); settingsFormEl.classList.add("hidden");
    eventFormEl.reset(); eventReminderInput.value = "1";
    eventDateInput.value = selectedCalendarDate + "T09:00";
    eventErrorEl.classList.add("hidden"); eventErrorEl.textContent = "";
    eventFormEl.classList.remove("hidden"); eventTitleInput.focus();
  }

  function submitEvent(event) {
    event.preventDefault();
    const payload = { title: eventTitleInput.value.trim(), event_at: eventDateInput.value, reminder_hours: Number(eventReminderInput.value) };
    if (!payload.title || !payload.event_at || payload.reminder_hours < 0 || payload.reminder_hours > 168) { eventErrorEl.textContent = t("fillAll"); eventErrorEl.classList.remove("hidden"); return; }
    apiFetch("/api/events", { method: "POST", body: JSON.stringify(payload) }).then(function () {
      eventFormEl.classList.add("hidden"); calendarMonth = payload.event_at.slice(0, 7); selectedCalendarDate = payload.event_at.slice(0, 10); return refreshEvents().then(loadAnalytics);
    }).catch(function (err) { eventErrorEl.textContent = t("eventError") + err.message; eventErrorEl.classList.remove("hidden"); });
  }

  function refreshEvents() { return fetchCalendarEvents().then(function (items) { calendarEvents = items || []; renderCalendar(); }); }

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
    updateMonthPickerLabel();
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
    updateMonthPickerLabel();
    loadAnalytics();
  }

  function updateMonthPickerLabel() {
    if (selectedMonth) monthPickerLabelEl.textContent = formatMonth(selectedMonth);
  }

  function formatMonth(month) {
    const parts = month.split("-");
    return new Date(Number(parts[0]), Number(parts[1]) - 1, 1)
      .toLocaleDateString(languageLocale(), { month: "long", year: "numeric" });
  }

  function languageLocale() {
    return { ru: "ru-RU", en: "en-US", zh: "zh-CN", es: "es-ES", ar: "ar", hy: "hy-AM" }[language] || "en-US";
  }

  function telegramLanguage() {
    const code = tg && tg.initDataUnsafe && tg.initDataUnsafe.user
      ? String(tg.initDataUnsafe.user.language_code || "").toLowerCase().split(/[-_]/)[0]
      : "";
    return ["ru", "en", "zh", "es", "ar", "hy"].includes(code) ? code : "en";
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
          ? date.toLocaleDateString(languageLocale(), { weekday: "short" })
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
    deadline.className = "item-deadline " + (item.completed ? "completed" : deadlineClass(item.days_left));
    deadline.textContent = item.completed ? formatDate(item.deadline) : formatDeadline(item.deadline, item.days_left);

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
    const date = new Date(Number(parts[0]), Number(parts[1]) - 1, Number(parts[2]));
    return date.toLocaleDateString(languageLocale(), { year: "numeric", month: "2-digit", day: "2-digit" });
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
    if (language !== "ru") return days === 1 ? t("dayUnitOne") : t("dayUnit");
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
    const locale = languageLocale();
    document.documentElement.lang = locale;
    document.documentElement.dir = language === "ar" ? "rtl" : "ltr";
    // Native date/month controls can use the control's own language instead
    // of inheriting the document language in some browsers.
    deadlineInput.lang = locale;
    monthPickerEl.lang = locale;
    updateDeadlineDisplay();
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
    languageSelectEl.value = language;
    addToggleEl.setAttribute("aria-label", currentSection === "calendar" ? t("addEvent") : t("addTaskAria"));
    updateMonthPickerLabel();
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
