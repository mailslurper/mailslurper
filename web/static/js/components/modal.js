// Thin wrapper around native <dialog> that restores focus to the invoking
// element on close (native <dialog> gives focus-trap + Esc for free, but
// not focus-restore).

export function openDialog(dialogEl) {
	const invoker = document.activeElement;
	dialogEl.showModal();

	const firstField = dialogEl.querySelector("input, select, textarea, button:not([data-dismiss])");
	if (firstField) firstField.focus();

	const onClose = () => {
		dialogEl.removeEventListener("close", onClose);
		if (invoker && typeof invoker.focus === "function") invoker.focus();
	};
	dialogEl.addEventListener("close", onClose);
}

export function closeDialog(dialogEl) {
	if (dialogEl.open) dialogEl.close();
}

// wireDismiss closes the dialog when any element with [data-dismiss] inside
// it is clicked.
export function wireDismiss(dialogEl) {
	dialogEl.querySelectorAll("[data-dismiss]").forEach((btn) => {
		btn.addEventListener("click", () => closeDialog(dialogEl));
	});
}
