// Small helpers for building DOM from <template> elements without a
// templating library. Fields are filled via textContent/attributes only —
// never innerHTML — so untrusted mail content can't inject markup.

export function cloneTemplate(id) {
	const tpl = document.getElementById(id);
	if (!tpl) throw new Error(`template #${id} not found`);
	return tpl.content.firstElementChild.cloneNode(true);
}

// fill(root, values) sets textContent on every descendant with
// data-field="name" matching a key in values, and sets attributes on
// elements with data-attr-<name>="field" (e.g. data-attr-href="url").
export function fill(root, values) {
	for (const [key, value] of Object.entries(values)) {
		const target = root.querySelector(`[data-field="${key}"]`);
		if (target) target.textContent = value ?? "";
	}
	for (const el of root.querySelectorAll("[data-attr-href]")) {
		const field = el.getAttribute("data-attr-href");
		if (field in values) el.href = values[field] ?? "";
	}
	return root;
}

export function clear(el) {
	while (el.firstChild) el.removeChild(el.firstChild);
}

export function el(tag, attrs = {}, text) {
	const node = document.createElement(tag);
	for (const [k, v] of Object.entries(attrs)) {
		if (k === "class") node.className = v;
		else node.setAttribute(k, v);
	}
	if (text !== undefined) node.textContent = text;
	return node;
}
