const state = {
  apiBase: resolveAPIBase(),
  tasks: [],
  filter: "all",
  editingID: null,
  pendingRequests: 0,
};

const refs = {
  apiBase: document.querySelector("#api-base"),
  createForm: document.querySelector("#create-form"),
  descriptionInput: document.querySelector("#description-input"),
  emptyState: document.querySelector("#empty-state"),
  feedback: document.querySelector("#feedback"),
  healthStatus: document.querySelector("#health-status"),
  taskCardTemplate: document.querySelector("#task-card-template"),
  taskCounter: document.querySelector("#task-counter"),
  taskList: document.querySelector("#task-list"),
  titleInput: document.querySelector("#title-input"),
  filterButtons: Array.from(document.querySelectorAll(".filter-button")),
};

refs.apiBase.textContent = state.apiBase;
bindEvents();
void refreshAll();
window.setInterval(() => {
  void refreshHealth();
}, 30000);

function bindEvents() {
  refs.createForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const payload = {
      title: refs.titleInput.value,
      description: refs.descriptionInput.value,
    };

    await withBusy(async () => {
      await apiRequest("/api/tasks", {
        method: "POST",
        body: JSON.stringify(payload),
      });

      refs.createForm.reset();
      setFeedback("Tarefa criada com sucesso.", true);
      await refreshAll();
    });
  });

  refs.filterButtons.forEach((button) => {
    button.addEventListener("click", () => {
      state.filter = button.dataset.filter || "all";
      state.editingID = null;
      refs.filterButtons.forEach((item) => item.classList.toggle("is-active", item === button));
      renderTasks();
    });
  });
}

async function refreshAll() {
  await Promise.all([refreshHealth(), refreshTasks()]);
}

async function refreshHealth() {
  try {
    await apiRequest("/health");
    setHealthStatus("ok", "status-pill--ok");
  } catch (error) {
    setHealthStatus(error.message || "degradado", "status-pill--degraded");
  }
}

async function refreshTasks() {
  await withBusy(async () => {
    const payload = await apiRequest("/api/tasks");
    state.tasks = Array.isArray(payload.tasks) ? payload.tasks : [];
    renderTasks();
  }, { quiet: true });
}

function renderTasks() {
  const visibleTasks = state.tasks.filter((task) => {
    if (state.filter === "pending") {
      return !task.completed;
    }

    if (state.filter === "done") {
      return task.completed;
    }

    return true;
  });

  refs.taskCounter.textContent = `${visibleTasks.length} tarefa(s)`;
  refs.taskList.replaceChildren();
  refs.emptyState.hidden = visibleTasks.length !== 0;

  if (visibleTasks.length === 0) {
    return;
  }

  visibleTasks.forEach((task) => {
    refs.taskList.appendChild(renderTaskCard(task));
  });
}

function renderTaskCard(task) {
  const fragment = refs.taskCardTemplate.content.cloneNode(true);
  const card = fragment.querySelector(".task-card");
  const status = fragment.querySelector('[data-field="status"]');
  const title = fragment.querySelector('[data-field="title"]');
  const description = fragment.querySelector('[data-field="description"]');
  const updatedAt = fragment.querySelector('[data-field="updated-at"]');
  const toggleButton = fragment.querySelector('[data-action="toggle"]');
  const editToggleButton = fragment.querySelector('[data-action="edit-toggle"]');
  const deleteButton = fragment.querySelector('[data-action="delete"]');
  const cancelButton = fragment.querySelector('[data-action="edit-cancel"]');
  const editForm = fragment.querySelector(".edit-form");
  const titleField = editForm.elements.namedItem("title");
  const descriptionField = editForm.elements.namedItem("description");

  status.textContent = task.completed ? "concluida" : "pendente";
  title.textContent = task.title;
  description.textContent = task.description || "Sem descricao adicional.";
  updatedAt.textContent = `Atualizada em ${formatDate(task.updatedAt)}`;
  toggleButton.textContent = task.completed ? "Reabrir" : "Concluir";
  titleField.value = task.title;
  descriptionField.value = task.description || "";

  card.classList.toggle("task-card--done", Boolean(task.completed));
  editForm.hidden = state.editingID !== task.id;

  toggleButton.addEventListener("click", async () => {
    await withBusy(async () => {
      await apiRequest(`/api/tasks/${task.id}`, {
        method: "PUT",
        body: JSON.stringify({ completed: !task.completed }),
      });

      setFeedback("Estado da tarefa atualizado.", true);
      await refreshAll();
    });
  });

  editToggleButton.addEventListener("click", () => {
    state.editingID = state.editingID === task.id ? null : task.id;
    renderTasks();
  });

  cancelButton.addEventListener("click", () => {
    state.editingID = null;
    renderTasks();
  });

  deleteButton.addEventListener("click", async () => {
    const confirmed = window.confirm(`Excluir a tarefa \"${task.title}\"?`);
    if (!confirmed) {
      return;
    }

    await withBusy(async () => {
      await apiRequest(`/api/tasks/${task.id}`, { method: "DELETE" });
      setFeedback("Tarefa removida.", true);
      await refreshAll();
    });
  });

  editForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    await withBusy(async () => {
      await apiRequest(`/api/tasks/${task.id}`, {
        method: "PUT",
        body: JSON.stringify({
          title: titleField.value,
          description: descriptionField.value,
        }),
      });

      state.editingID = null;
      setFeedback("Tarefa atualizada.", true);
      await refreshAll();
    });
  });

  return fragment;
}

async function apiRequest(pathname, options = {}) {
  const response = await fetch(`${state.apiBase}${pathname}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const contentType = response.headers.get("content-type") || "";
  const isJSON = contentType.includes("application/json");
  const payload = isJSON ? await response.json() : null;

  if (!response.ok) {
    const message = payload && typeof payload.error === "string"
      ? payload.error
      : payload && typeof payload.message === "string"
        ? payload.message
        : "Nao foi possivel concluir a operacao.";
    throw new Error(message);
  }

  return payload;
}

async function withBusy(action, options = {}) {
  state.pendingRequests += 1;
  updateBusyState();

  try {
    await action();
  } catch (error) {
    if (!options.quiet) {
      setFeedback(error.message || "Falha inesperada.", false);
    }
  } finally {
    state.pendingRequests -= 1;
    updateBusyState();
  }
}

function updateBusyState() {
  const busy = state.pendingRequests > 0;
  document.querySelectorAll("button, input, textarea").forEach((element) => {
    if (element === refs.titleInput || element === refs.descriptionInput || element.form !== refs.createForm) {
      if (busy && element.tagName === "BUTTON") {
        element.disabled = true;
      } else if (!busy) {
        element.disabled = false;
      }
      return;
    }

    element.disabled = busy;
  });
}

function setFeedback(message, success) {
  refs.feedback.hidden = false;
  refs.feedback.textContent = message;
  refs.feedback.classList.toggle("feedback--success", Boolean(success));
}

function setHealthStatus(message, className) {
  refs.healthStatus.textContent = message;
  refs.healthStatus.className = `status-pill ${className}`;
}

function formatDate(value) {
  if (!value) {
    return "agora";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "agora";
  }

  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(date);
}

function resolveAPIBase() {
  const params = new URLSearchParams(window.location.search);
  const queryValue = params.get("api");
  if (queryValue) {
    return stripTrailingSlash(queryValue);
  }

  if (typeof window.GOTODOLIST_API_BASE_URL === "string" && window.GOTODOLIST_API_BASE_URL.trim() !== "") {
    return stripTrailingSlash(window.GOTODOLIST_API_BASE_URL);
  }

  const { protocol, hostname } = window.location;
  if (!hostname) {
    return "http://localhost:8081";
  }

  return `${protocol}//${hostname}:8081`;
}

function stripTrailingSlash(value) {
  return value.replace(/\/$/, "");
}