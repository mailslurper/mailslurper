import { el } from "../util/dom.js";

let container = null;

function ensureContainer() {
	if (container) return container;
	container = el("div", { class: "toast-container", role: "region", "aria-label": "Notifications" });
	document.body.appendChild(container);
	return container;
}

function show(kind, message, timeoutMs = 4000) {
	const c = ensureContainer();
	const toast = el("div", {
		class: `toast toast-${kind}`,
		role: kind === "error" ? "alert" : "status",
		"aria-live": kind === "error" ? "assertive" : "polite",
	}, message);

	const closeBtn = el("button", { class: "btn btn-icon", "aria-label": "Dismiss", style: "float:right;border:none;background:none;" }, "×");
	closeBtn.addEventListener("click", () => toast.remove());
	toast.appendChild(closeBtn);

	c.appendChild(toast);
	if (timeoutMs > 0) setTimeout(() => toast.remove(), timeoutMs);
}

export const toast = {
	success: (msg) => show("success", msg),
	error: (msg) => show("error", msg, 6000),
	info: (msg) => show("info", msg),
};
