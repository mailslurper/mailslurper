import { el } from "../util/dom.js";
import { openDialog } from "./modal.js";

let dialog = null;

function ensureDialog() {
	if (dialog) return dialog;
	dialog = el("dialog", { "aria-label": "Image preview" });
	dialog.style.maxWidth = "90vw";
	document.body.appendChild(dialog);
	dialog.addEventListener("click", (e) => {
		if (e.target === dialog) dialog.close();
	});
	return dialog;
}

export function openImage(url, fileName) {
	const d = ensureDialog();
	d.textContent = "";
	const img = el("img", { src: url, alt: fileName, style: "max-width:100%;max-height:80vh;display:block;" });
	d.appendChild(img);
	openDialog(d);
}
