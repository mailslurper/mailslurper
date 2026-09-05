import { openDialog, closeDialog, wireDismiss } from "./modal.js";
import { openPicker } from "./savedSearchPicker.js";
import * as store from "../store.js";
import { toast } from "./toast.js";

let dialog, form, fields;
let applyCallback = null;

function init() {
	if (dialog) return;
	dialog = document.getElementById("search-modal");
	form = document.getElementById("search-form");
	fields = {
		q: document.getElementById("search-q"),
		from: document.getElementById("search-from"),
		to: document.getElementById("search-to"),
		start: document.getElementById("search-start"),
		end: document.getElementById("search-end"),
	};

	wireDismiss(dialog);

	form.addEventListener("submit", (e) => {
		e.preventDefault();
		applyCallback?.(readFields());
		closeDialog(dialog);
	});

	document.getElementById("search-clear-btn").addEventListener("click", () => {
		for (const input of Object.values(fields)) input.value = "";
	});

	document.getElementById("search-save-btn").addEventListener("click", () => {
		const name = prompt("Name this search:");
		if (!name) return;
		store.saveSavedSearch({ name, ...readFields() });
		toast.success(`Saved search "${name}"`);
	});

	document.getElementById("search-saved-btn").addEventListener("click", () => {
		openPicker((search) => writeFields(search));
	});
}

function readFields() {
	return {
		q: fields.q.value.trim(),
		from: fields.from.value.trim(),
		to: fields.to.value.trim(),
		start: fields.start.value,
		end: fields.end.value,
	};
}

function writeFields(state) {
	fields.q.value = state.q || "";
	fields.from.value = state.from || "";
	fields.to.value = state.to || "";
	fields.start.value = state.start || "";
	fields.end.value = state.end || "";
}

export function open(currentState, onApply) {
	init();
	writeFields(currentState);
	applyCallback = onApply;
	openDialog(dialog);
}
