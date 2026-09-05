import { el, clear } from "../util/dom.js";

// render(container, { page, totalPages, onChange }) draws first/prev/next/last
// controls and a "page X of Y" label.
export function render(container, { page, totalPages, onChange }) {
	clear(container);
	container.className = "pager";

	const makeBtn = (label, targetPage, disabled) => {
		const btn = el("button", { class: "btn btn-icon", type: "button" }, label);
		btn.disabled = disabled;
		btn.addEventListener("click", () => onChange(targetPage));
		return btn;
	};

	container.appendChild(makeBtn("«", 1, page <= 1));
	container.appendChild(makeBtn("‹", page - 1, page <= 1));
	container.appendChild(el("span", { class: "pager-label" }, totalPages > 0 ? `Page ${page} of ${totalPages}` : "No results"));
	container.appendChild(makeBtn("›", page + 1, page >= totalPages));
	container.appendChild(makeBtn("»", totalPages, page >= totalPages || totalPages === 0));
}
