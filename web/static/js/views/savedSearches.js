import * as store from "../store.js";
import { el, clear } from "../util/dom.js";
import { toast } from "../components/toast.js";

// renderList draws saved searches as a list of buttons; onSelect(search) is
// called when one is chosen. Shared by the full management page and the
// in-modal picker so there's exactly one rendering path for this data.
export function renderList(container, onSelect, { showDelete = false } = {}) {
	clear(container);
	const searches = store.getSavedSearches();

	if (searches.length === 0) {
		container.appendChild(el("p", { class: "empty-state" }, "No saved searches yet."));
		return;
	}

	const list = el("ul", { class: "saved-list" });
	for (const search of searches) {
		const item = el("li");
		const btn = el("button", { class: "btn", type: "button" }, search.name);
		btn.addEventListener("click", () => onSelect(search));
		item.appendChild(btn);

		if (showDelete) {
			const del = el("button", { class: "btn btn-danger", type: "button", "aria-label": `Delete ${search.name}` }, "Delete");
			del.addEventListener("click", () => {
				store.deleteSavedSearch(search.name);
				renderList(container, onSelect, { showDelete });
				toast.info(`Deleted "${search.name}"`);
			});
			item.appendChild(del);
		}
		list.appendChild(item);
	}
	container.appendChild(list);
}

let unsubscribe = null;

export function mount(container) {
	container.className = "view view-saved";
	const page = el("div", { class: "saved-page" });
	page.appendChild(el("h1", {}, "Saved searches"));
	const list = el("div");
	page.appendChild(list);
	container.appendChild(page);

	renderList(list, (search) => {
		location.hash = `#/inbox?${new URLSearchParams(search).toString()}`;
	}, { showDelete: true });

	unsubscribe = store.on("savedSearches", () => renderList(list, (search) => {
		location.hash = `#/inbox?${new URLSearchParams(search).toString()}`;
	}, { showDelete: true }));
}

export function unmount() {
	unsubscribe?.();
	unsubscribe = null;
}
