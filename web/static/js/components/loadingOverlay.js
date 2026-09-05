import { el } from "../util/dom.js";

// show(container) attaches a blocking overlay with a spinner inside
// container (which must be position: relative/absolute for it to cover
// correctly) and returns a hide() function.
export function show(container, text = "Loading…") {
	const overlay = el("div", { class: "loading-overlay", role: "status", "aria-live": "polite" });
	overlay.appendChild(el("div", { class: "spinner", "aria-hidden": "true" }));
	overlay.appendChild(el("span", { class: "visually-hidden" }, text));
	container.appendChild(overlay);
	return () => overlay.remove();
}
