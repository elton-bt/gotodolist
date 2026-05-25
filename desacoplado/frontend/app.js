const state = {
  apiBase: resolveAPIBase(),
  version: resolveAppVersion(),
  frontendInstance: "carregando",
  frontendIP: "carregando",
  apiInstance: "aguardando resposta",
  apiIP: "aguardando",
  tasks: [],
  filter: "all",
  editingID: null,
  pendingRequests: 0,
};

const refs = {
  apiBase: document.querySelector("#api-base"),
  apiInstance: document.querySelector("#api-instance"),
  apiIP: document.querySelector("#api-ip"),
  appVersion: document.querySelector("#app-version"),
  createForm: document.querySelector("#create-form"),
  descriptionInput: document.querySelector("#description-input"),
  emptyState: document.querySelector("#empty-state"),
  feedback: document.querySelector("#feedback"),
  frontendInstance: document.querySelector("#frontend-instance"),
  frontendIP: document.querySelector("#frontend-ip"),
  healthStatus: document.querySelector("#health-status"),
  taskCardTemplate: document.querySelector("#task-card-template"),
  taskCounter: document.querySelector("#task-counter"),
  taskList: document.querySelector("#task-list"),
  titleInput: document.querySelector("#title-input"),
  filterButtons: Array.from(document.querySelectorAll(".filter-button")),
};

refs.apiBase.textContent = state.apiBase;
refs.appVersion.textContent = state.version;
setIdentityCard(refs.frontendInstance, refs.frontendIP, state.frontendInstance, state.frontendIP);
setIdentityCard(refs.apiInstance, refs.apiIP, state.apiInstance, state.apiIP);
bindEvents();
void refreshAll();
window.setInterval(() => {
  void Promise.all([refreshFrontendInfo(), refreshHealth()]);
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
  await Promise.all([refreshFrontendInfo(), refreshHealth(), refreshTasks()]);
}

async function refreshFrontendInfo() {
  try {
    const response = await fetch("/health", {
      cache: "no-store",
    });

    applyFrontendIdentity(response.headers);
  } catch {
    setIdentityCard(refs.frontendInstance, refs.frontendIP, "indisponivel", "indisponivel");
  }
}

async function refreshHealth() {
  try {
    await apiRequest("/health", { cache: "no-store" });
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

  applyAPIIdentity(response.headers);

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

function applyFrontendIdentity(headers) {
  const instance = headers.get("X-GoToDoList-Frontend-Instance");
  const ip = headers.get("X-GoToDoList-Frontend-IP");

  if (!instance && !ip) {
    return;
  }

  state.frontendInstance = instance || state.frontendInstance;
  state.frontendIP = ip || state.frontendIP;
  setIdentityCard(refs.frontendInstance, refs.frontendIP, state.frontendInstance, state.frontendIP);
}

function applyAPIIdentity(headers) {
  const instance = headers.get("X-GoToDoList-Instance");
  const ip = headers.get("X-GoToDoList-IP");

  if (!instance && !ip) {
    return;
  }

  state.apiInstance = instance || state.apiInstance;
  state.apiIP = ip || state.apiIP;
  setIdentityCard(refs.apiInstance, refs.apiIP, state.apiInstance, state.apiIP);
}

function setIdentityCard(instanceElement, ipElement, instance, ip) {
  instanceElement.textContent = instance;
  ipElement.textContent = `IP ${ip}`;
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
  if (
    window.GOTODOLIST_CONFIG &&
    typeof window.GOTODOLIST_CONFIG.apiBase === "string" &&
    window.GOTODOLIST_CONFIG.apiBase.trim() !== ""
  ) {
    return stripTrailingSlash(window.GOTODOLIST_CONFIG.apiBase);
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

function resolveAppVersion() {
  if (
    window.GOTODOLIST_CONFIG &&
    typeof window.GOTODOLIST_CONFIG.version === "string" &&
    window.GOTODOLIST_CONFIG.version.trim() !== ""
  ) {
    return window.GOTODOLIST_CONFIG.version.trim();
  }

  if (typeof window.GOTODOLIST_APP_VERSION === "string" && window.GOTODOLIST_APP_VERSION.trim() !== "") {
    return window.GOTODOLIST_APP_VERSION.trim();
  }

  return "dev";
}

function stripTrailingSlash(value) {
  return value.replace(/\/$/, "");
}