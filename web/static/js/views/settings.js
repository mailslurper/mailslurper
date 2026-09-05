import * as apiClient from "../api.js";
import { el } from "../util/dom.js";
import * as store from "../store.js";
import { toast } from "../components/toast.js";

const THEMES = [
	["system", "System"],
	["light", "Light"],
	["dark", "Dark"],
];
const DATE_FORMATS = [
	["locale", "Browser default"],
	["iso", "ISO 8601"],
	["us", "US (M/D/Y)"],
	["eu", "EU (D/M/Y)"],
];

export function mount(container) {
	container.className = "view view-settings";

	const page = el("div", { class: "settings-page" });
	const appearance = el("div");
	appearance.appendChild(el("h1", {}, "Settings"));
	const fields = el("div", { class: "settings-fields" });
	fields.appendChild(buildThemeField());
	fields.appendChild(buildDateFormatField());
	fields.appendChild(buildAutoRefreshField());
	appearance.appendChild(fields);
	page.appendChild(appearance);
	page.appendChild(buildPruneSection());
	page.appendChild(buildVersionSection());
	container.appendChild(page);
}

export function unmount() {}

function buildThemeField() {
	const field = el("div", { class: "field" });
	field.appendChild(el("label", { for: "theme-select" }, "Theme"));
	const select = el("select", { id: "theme-select" });
	for (const [value, label] of THEMES) {
		const opt = el("option", { value }, label);
		if (value === store.getTheme()) opt.selected = true;
		select.appendChild(opt);
	}
	select.addEventListener("change", () => store.setTheme(select.value));
	field.appendChild(select);
	return field;
}

function buildDateFormatField() {
	const field = el("div", { class: "field" });
	field.appendChild(el("label", { for: "date-format-select" }, "Date format"));
	const select = el("select", { id: "date-format-select" });
	for (const [value, label] of DATE_FORMATS) {
		const opt = el("option", { value }, label);
		if (value === store.getDateFormat()) opt.selected = true;
		select.appendChild(opt);
	}
	select.addEventListener("change", () => store.setDateFormat(select.value));
	field.appendChild(select);
	return field;
}

function buildAutoRefreshField() {
	const field = el("div", { class: "field" });
	field.appendChild(el("label", { for: "auto-refresh-input" }, "Fallback refresh interval (minutes, 0 to disable)"));
	field.appendChild(el("p", { class: "field-help" }, "Used only when the live connection is interrupted."));
	const input = el("input", { id: "auto-refresh-input", type: "number", min: "0", step: "1" });
	input.value = String(store.getAutoRefreshMinutes());
	input.addEventListener("change", () => store.setAutoRefreshMinutes(Number(input.value) || 0));
	field.appendChild(input);
	return field;
}

function buildPruneSection() {
	const section = el("div");
	section.appendChild(el("h2", {}, "Delete mail"));
	const row = el("div", { class: "prune-row" });
	const select = el("select", { id: "prune-select" });
	row.appendChild(select);

	apiClient.getPruneOptions().then((options) => {
		for (const opt of options) {
			select.appendChild(el("option", { value: opt.code }, opt.description));
		}
	}).catch((err) => toast.error(`Could not load prune options: ${err.message}`));

	const btn = el("button", { class: "btn btn-danger", type: "button" }, "Delete");
	btn.addEventListener("click", async () => {
		const code = select.value;
		if (!code) return;
		const label = select.selectedOptions[0]?.textContent || code;
		if (!confirm(`Delete mail matching "${label}"? This cannot be undone.`)) return;

		try {
			const result = await apiClient.pruneMail(code);
			toast.success(`Deleted ${result.deletedCount} message${result.deletedCount === 1 ? "" : "s"}`);
		} catch (err) {
			toast.error(`Could not delete mail: ${err.message}`);
		}
	});
	row.appendChild(btn);
	section.appendChild(row);
	return section;
}

function buildVersionSection() {
	const section = el("div");
	section.appendChild(el("h2", {}, "Version"));
	const current = el("p", { class: "settings-version" }, "Loading…");
	section.appendChild(current);

	apiClient.getVersion()
		.then((v) => { current.textContent = `Running version ${v.version}`; })
		.catch(() => { current.textContent = "Could not determine running version."; });

	const checkBtn = el("button", { class: "btn", type: "button" }, "Check latest release on GitHub");
	const result = el("p", { class: "field-help" });
	checkBtn.addEventListener("click", async () => {
		result.textContent = "Checking…";
		try {
			const res = await fetch("https://api.github.com/repos/p0vidl0/mylslurper/releases/latest");
			if (!res.ok) throw new Error(`GitHub responded with ${res.status}`);
			const data = await res.json();
			result.textContent = `Latest release: ${data.tag_name}`;
		} catch (err) {
			result.textContent = `Could not check for updates: ${err.message}`;
		}
	});
	section.appendChild(checkBtn);
	section.appendChild(result);
	return section;
}
