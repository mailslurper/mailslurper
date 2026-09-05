import * as apiClient from "../api.js";
import { UnauthorizedError } from "../api.js";
import { cloneTemplate, fill, el, clear } from "../util/dom.js";
import { formatDateTime, formatCountdown } from "../util/date.js";
import { queryStringToState, stateToQueryString } from "../util/query.js";
import * as pager from "../components/pager.js";
import * as searchModal from "../components/searchModal.js";
import * as store from "../store.js";
import { toast } from "../components/toast.js";
import { show as showLoading } from "../components/loadingOverlay.js";
import * as mailDetail from "./mailDetail.js";
import * as liveEvents from "../events.js";
import { attachVerticalSplitter, attachColumnSplitters } from "../components/splitter.js";

const SORT_FIELDS = ["date", "subject", "from"];

let searchState = queryStringToState("");
let selectedId = null;
let refreshTimer = null;
let countdownTimer = null;
let nextRefreshAt = 0;
let sseConnected = false;
let disconnectSSE = null;
let unsubscribeConnection = null;
let unsubscribeAutoRefresh = null;
let detachSplitter = null;
let detachColumns = null;

let els = {};

export function mount(container, params) {
	searchState = queryStringToState(params.queryString || "");
	selectedId = params.selectedId || null;

	container.className = "view view-inbox";
	container.appendChild(buildLayout());
	wireEvents();
	updateSortIndicators();

	load();
	startLiveUpdates();
}

export function unmount() {
	stopLiveUpdates();
	if (detachSplitter) detachSplitter();
	detachSplitter = null;
	if (detachColumns) detachColumns();
	detachColumns = null;
	selectedId = null;
	els = {};
}

function buildLayout() {
	const layout = el("div", { class: "inbox-layout", id: "inbox-layout" });

	const listPane = el("div", { class: "inbox-list-pane" });
	const toolbar = el("div", { class: "list-toolbar" });
	els.searchBtn = el("button", { class: "btn btn-search", type: "button" });
	els.searchBtn.appendChild(el("span", { class: "search-icon", "aria-hidden": "true" }));
	els.searchBtn.appendChild(document.createTextNode("Search"));
	const meta = el("div", { class: "toolbar-meta" });
	els.countdown = el("span", { class: "pill" }, "");
	els.countLabel = el("span", { class: "pill" }, "");
	meta.appendChild(els.countdown);
	meta.appendChild(els.countLabel);
	toolbar.appendChild(els.searchBtn);
	toolbar.appendChild(meta);
	listPane.appendChild(toolbar);

	els.head = el("div", { class: "mail-list-head" });
	els.colSplitters = [];
	const headCols = [
		{ field: "date", resize: "date", label: "Date" },
		{ field: "subject", resize: "from", label: "Subject" },
		{ field: "from", resize: "att", label: "From" },
		{ field: null, resize: null, label: "Att." },
	];
	for (const col of headCols) {
		const cell = el("div", { class: "mail-col-head" });
		if (col.field) {
			const btn = el("button", { type: "button", "data-sort": col.field }, col.label);
			btn.addEventListener("click", () => onSortClick(col.field));
			cell.appendChild(btn);
		} else {
			cell.appendChild(el("div", {}, col.label));
		}
		if (col.resize) {
			const splitter = el("button", {
				class: "col-splitter",
				type: "button",
				role: "separator",
				"aria-orientation": "vertical",
				"data-col": col.resize,
				"aria-label": `Resize ${col.resize === "att" ? "Attachments" : sortLabel(col.resize)} column`,
				title: "Drag to resize. Double-click to reset.",
			});
			cell.appendChild(splitter);
			els.colSplitters.push(splitter);
		}
		els.head.appendChild(cell);
	}
	listPane.appendChild(els.head);

	const scroll = el("div", { class: "inbox-list-scroll" });
	els.list = el("div", { class: "mail-list", role: "rowgroup" });
	scroll.appendChild(els.list);
	listPane.appendChild(scroll);

	els.pager = el("div");
	listPane.appendChild(els.pager);

	layout.appendChild(listPane);

	els.listPane = listPane;
	els.splitter = el("button", {
		class: "inbox-splitter",
		type: "button",
		role: "separator",
		"aria-orientation": "vertical",
		"aria-label": "Resize mail list",
		title: "Drag to resize. Double-click to reset.",
	});
	layout.appendChild(els.splitter);

	els.detailPane = el("div", { class: "inbox-detail-pane" });
	layout.appendChild(els.detailPane);

	els.layout = layout;
	return layout;
}

function sortLabel(field) {
	const names = { date: "Date", subject: "Subject", from: "From" };
	return names[field];
}

function wireEvents() {
	els.searchBtn.addEventListener("click", () => {
		searchModal.open(searchState, (filters) => {
			searchState = { ...searchState, ...filters, page: 1 };
			syncURL();
			load();
		});
	});

	detachSplitter = attachVerticalSplitter({
		layout: els.layout,
		pane: els.listPane,
		splitter: els.splitter,
		getWidth: store.getInboxListWidth,
		setWidth: store.setInboxListWidth,
	});

	detachColumns = attachColumnSplitters({
		root: els.listPane,
		splitters: els.colSplitters,
		getWidths: store.getInboxColumnWidths,
		setWidths: store.setInboxColumnWidths,
	});
}

function onSortClick(field) {
	if (searchState.sort === field) {
		searchState.dir = searchState.dir === "asc" ? "desc" : "asc";
	} else {
		searchState.sort = field;
		searchState.dir = "asc";
	}
	searchState.page = 1;
	syncURL();
	updateSortIndicators();
	load();
}

function updateSortIndicators() {
	els.layout.querySelectorAll(".mail-list-head button[data-sort]").forEach((btn) => {
		const field = btn.getAttribute("data-sort");
		if (field === searchState.sort) {
			btn.setAttribute("aria-sort", searchState.dir === "asc" ? "ascending" : "descending");
			btn.textContent = `${sortLabel(field)} ${searchState.dir === "asc" ? "▲" : "▼"}`;
		} else {
			btn.removeAttribute("aria-sort");
			btn.textContent = sortLabel(field);
		}
	});
}

function syncURL() {
	const qs = stateToQueryString(searchState);
	const base = selectedId ? `#/mail/${encodeURIComponent(selectedId)}` : "#/inbox";
	const newHash = qs ? `${base}?${qs}` : base;
	if (location.hash !== newHash) {
		history.replaceState(null, "", newHash);
	}
}

async function load(quiet = false) {
	const hideList = quiet ? null : showLoading(els.layout.querySelector(".inbox-list-scroll"));
	try {
		const result = await apiClient.getMails(searchState);
		renderRows(result.mailItems);
		pager.render(els.pager, {
			page: result.page,
			totalPages: result.totalPages,
			onChange: (page) => {
				searchState.page = page;
				syncURL();
				load();
			},
		});
		els.countLabel.textContent = `${result.totalRecords} message${result.totalRecords === 1 ? "" : "s"}`;
		els.layout.classList.toggle("showing-detail", Boolean(selectedId));

		if (selectedId) {
			const stillVisible = result.mailItems.some((item) => item.id === selectedId);
			if (stillVisible) {
				await loadDetail(selectedId);
			} else {
				selectedId = null;
				syncURL();
				mailDetail.renderEmpty(els.detailPane);
			}
		} else {
			mailDetail.renderEmpty(els.detailPane);
		}
	} catch (err) {
		if (err instanceof UnauthorizedError) {
			location.hash = "#/login";
			return;
		}
		toast.error(`Could not load mail: ${err.message}`);
	} finally {
		if (hideList) hideList();
	}
}

function loadQuiet() {
	return load(true);
}

function renderRows(items) {
	clear(els.list);
	if (items.length === 0) {
		els.head.hidden = true;
		els.list.appendChild(buildEmptyState(searchStateHasFilters()));
		return;
	}
	els.head.hidden = false;

	for (const item of items) {
		const row = cloneTemplate("tpl-mail-row");
		fill(row, {
			dateSent: formatDateTime(item.dateSent, store.getDateFormat()),
			subject: item.subject || "(no subject)",
			fromAddress: item.fromAddress,
		});
		if (item.attachmentCount > 0) {
			row.querySelector(".mail-row-att").appendChild(el("span", { class: "att-dot", title: `${item.attachmentCount} attachment${item.attachmentCount === 1 ? "" : "s"}` }));
		}
		row.dataset.id = item.id;
		if (item.id === selectedId) row.classList.add("selected");
		row.addEventListener("click", () => selectMail(item.id));
		row.addEventListener("keydown", (e) => {
			if (e.key === "Enter" || e.key === " ") {
				e.preventDefault();
				selectMail(item.id);
			}
		});
		els.list.appendChild(row);
	}
}

function searchStateHasFilters() {
	return Boolean(searchState.q || searchState.from || searchState.to || searchState.start || searchState.end);
}

function buildEmptyState(filtered) {
	const wrap = el("div", { class: "empty-state" });
	const box = el("div", { class: "empty-inbox" });
	const icon = el("div", { class: "empty-inbox-icon" });
	icon.appendChild(el("span"));
	box.appendChild(icon);
	box.appendChild(el("div", { class: "empty-inbox-title" }, filtered ? "No messages found" : "No messages yet"));
	box.appendChild(el("div", { class: "empty-inbox-copy" },
		filtered
			? "Try a different search or clear the filters."
			: "Mail sent to your test SMTP endpoint will show up here.",
	));
	wrap.appendChild(box);
	return wrap;
}

function selectMail(id) {
	selectedId = id;
	els.list.querySelectorAll(".mail-row").forEach((row) => {
		row.classList.toggle("selected", row.dataset.id === id);
	});
	els.layout.classList.add("showing-detail");
	syncURL();
	loadDetail(id);
}

async function loadDetail(id) {
	try {
		const item = await apiClient.getMailByID(id);
		mailDetail.renderDetail(els.detailPane, item);
	} catch (err) {
		if (err instanceof UnauthorizedError) {
			location.hash = "#/login";
			return;
		}
		toast.error(`Could not load message: ${err.message}`);
	}
}

function onLiveEvent(type) {
	if (type === "mail.received" || type === "mail.pruned") {
		loadQuiet();
	}
}

function onConnectionChange(connected) {
	sseConnected = connected;
	updateStatusBadge();
	if (connected) {
		stopFallbackPolling();
	} else {
		startFallbackPolling();
	}
}

function updateStatusBadge() {
	if (sseConnected) {
		els.countdown.textContent = "Live";
		return;
	}

	const minutes = store.getAutoRefreshMinutes();
	if (!minutes || minutes <= 0) {
		els.countdown.textContent = "Reconnecting…";
		return;
	}

	const remaining = nextRefreshAt - Date.now();
	if (remaining > 0) {
		els.countdown.textContent = `Reconnecting… refresh in ${formatCountdown(remaining)}`;
	} else {
		els.countdown.textContent = "Reconnecting…";
	}
}

function startLiveUpdates() {
	disconnectSSE = liveEvents.connect(onLiveEvent);
	unsubscribeConnection = liveEvents.onConnectionChange(onConnectionChange);
	unsubscribeAutoRefresh = store.on("autoRefreshMinutes", () => {
		if (!sseConnected) {
			startFallbackPolling();
		}
	});
}

function stopLiveUpdates() {
	if (disconnectSSE) disconnectSSE();
	disconnectSSE = null;
	liveEvents.disconnect();
	if (unsubscribeConnection) unsubscribeConnection();
	unsubscribeConnection = null;
	if (unsubscribeAutoRefresh) unsubscribeAutoRefresh();
	unsubscribeAutoRefresh = null;
	stopFallbackPolling();
}

function startFallbackPolling() {
	stopFallbackPolling();
	const minutes = store.getAutoRefreshMinutes();
	if (!minutes || minutes <= 0) {
		updateStatusBadge();
		return;
	}

	nextRefreshAt = Date.now() + minutes * 60_000;
	refreshTimer = setInterval(() => {
		nextRefreshAt = Date.now() + minutes * 60_000;
		loadQuiet();
	}, minutes * 60_000);

	countdownTimer = setInterval(updateStatusBadge, 1000);
	updateStatusBadge();
}

function stopFallbackPolling() {
	if (refreshTimer) clearInterval(refreshTimer);
	if (countdownTimer) clearInterval(countdownTimer);
	refreshTimer = null;
	countdownTimer = null;
}
