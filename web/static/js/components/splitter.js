import { DEFAULT_INBOX_COLUMNS } from "../store.js";

const MIN_LIST = 280;
const MIN_DETAIL = 320;
const STEP = 16;
const STEP_LARGE = 48;

// attachVerticalSplitter wires a drag/keyboard handle that resizes `pane`
// inside `layout`. Returns a detach() function for unmount.
export function attachVerticalSplitter({ layout, pane, splitter, getWidth, setWidth }) {
	function bounds() {
		const total = layout.clientWidth;
		return { min: MIN_LIST, max: Math.max(MIN_LIST, total - MIN_DETAIL) };
	}

	function apply(px, persist) {
		const { min, max } = bounds();
		const w = Math.round(Math.min(max, Math.max(min, px)));
		pane.style.width = `${w}px`;
		splitter.setAttribute("aria-valuemin", String(min));
		splitter.setAttribute("aria-valuemax", String(max));
		splitter.setAttribute("aria-valuenow", String(w));
		if (persist) setWidth(w);
		return w;
	}

	function restore() {
		const saved = getWidth();
		if (saved) {
			apply(saved, false);
			return;
		}
		pane.style.width = "";
		const { min, max } = bounds();
		splitter.setAttribute("aria-valuemin", String(min));
		splitter.setAttribute("aria-valuemax", String(max));
		splitter.setAttribute("aria-valuenow", String(Math.round(pane.getBoundingClientRect().width)));
	}

	function reset() {
		setWidth(0);
		restore();
	}

	let dragging = false;

	function onPointerDown(e) {
		if (e.pointerType === "mouse" && e.button !== 0) return;
		dragging = true;
		layout.classList.add("is-resizing");
		splitter.setPointerCapture(e.pointerId);
		e.preventDefault();
	}

	function onPointerMove(e) {
		if (!dragging) return;
		apply(e.clientX - layout.getBoundingClientRect().left, false);
	}

	function onPointerUp() {
		if (!dragging) return;
		dragging = false;
		layout.classList.remove("is-resizing");
		setWidth(Math.round(pane.getBoundingClientRect().width));
	}

	function onKeyDown(e) {
		const current = pane.getBoundingClientRect().width;
		const step = e.shiftKey ? STEP_LARGE : STEP;
		if (e.key === "ArrowLeft") {
			e.preventDefault();
			apply(current - step, true);
		} else if (e.key === "ArrowRight") {
			e.preventDefault();
			apply(current + step, true);
		} else if (e.key === "Home") {
			e.preventDefault();
			reset();
		}
	}

	function onWindowResize() {
		const current = pane.getBoundingClientRect().width;
		if (pane.style.width) apply(current, false);
		else restore();
	}

	splitter.addEventListener("pointerdown", onPointerDown);
	splitter.addEventListener("pointermove", onPointerMove);
	splitter.addEventListener("pointerup", onPointerUp);
	splitter.addEventListener("pointercancel", onPointerUp);
	splitter.addEventListener("keydown", onKeyDown);
	splitter.addEventListener("dblclick", reset);
	window.addEventListener("resize", onWindowResize);

	restore();

	return () => {
		dragging = false;
		layout.classList.remove("is-resizing");
		splitter.removeEventListener("pointerdown", onPointerDown);
		splitter.removeEventListener("pointermove", onPointerMove);
		splitter.removeEventListener("pointerup", onPointerUp);
		splitter.removeEventListener("pointercancel", onPointerUp);
		splitter.removeEventListener("keydown", onKeyDown);
		splitter.removeEventListener("dblclick", reset);
		window.removeEventListener("resize", onWindowResize);
	};
}

const COL_MIN = { date: 64, from: 64, att: 20 };
const COL_MAX_ATT = 48;
const SUBJECT_MIN = 80;
const COL_DIR = { date: 1, from: -1, att: -1 };
function clamp(n, min, max) {
	return Math.round(Math.min(max, Math.max(min, n)));
}

// attachColumnSplitters lets the user drag the right edge of Date / Subject /
// From to resize that border. Subject stays flexible; Date, From, and Att
// are stored as pixel widths. Returns detach().
export function attachColumnSplitters({ root, splitters, getWidths, setWidths }) {
	let widths = { ...DEFAULT_INBOX_COLUMNS, ...getWidths() };

	function innerWidth() {
		const head = root.querySelector(".mail-list-head");
		if (!head || head.hidden) return 0;
		const style = getComputedStyle(head);
		const pad = parseFloat(style.paddingLeft) + parseFloat(style.paddingRight);
		const gap = parseFloat(style.columnGap || style.gap || "0") * 3;
		return Math.max(0, head.clientWidth - pad - gap);
	}

	function apply(next, persist) {
		const inner = innerWidth();
		const w = { ...next };
		if (inner > 0) {
			const maxDate = inner - COL_MIN.from - COL_MIN.att - SUBJECT_MIN;
			w.date = clamp(w.date, COL_MIN.date, Math.max(COL_MIN.date, maxDate));
			const maxFrom = inner - w.date - COL_MIN.att - SUBJECT_MIN;
			w.from = clamp(w.from, COL_MIN.from, Math.max(COL_MIN.from, maxFrom));
			const maxAtt = Math.min(COL_MAX_ATT, inner - w.date - w.from - SUBJECT_MIN);
			w.att = clamp(w.att, COL_MIN.att, Math.max(COL_MIN.att, maxAtt));
		}
		root.style.setProperty("--col-date", `${w.date}px`);
		root.style.setProperty("--col-from", `${w.from}px`);
		root.style.setProperty("--col-att", `${w.att}px`);
		widths = w;
		for (const splitter of splitters) {
			const col = splitter.dataset.col;
			splitter.setAttribute("aria-valuenow", String(w[col]));
			splitter.setAttribute("aria-valuemin", String(COL_MIN[col]));
		}
		if (persist) setWidths(w);
		return w;
	}

	function resetColumn(col) {
		apply({ ...widths, [col]: DEFAULT_INBOX_COLUMNS[col] }, true);
	}

	function resetAll() {
		setWidths(null);
		apply({ ...DEFAULT_INBOX_COLUMNS }, false);
	}

	const cleanups = [];
	let dragging = null;

	for (const splitter of splitters) {
		const col = splitter.dataset.col;

		function onPointerDown(e) {
			if (e.pointerType === "mouse" && e.button !== 0) return;
			dragging = { col, startX: e.clientX, start: widths[col] };
			root.classList.add("is-col-resizing");
			splitter.setPointerCapture(e.pointerId);
			e.preventDefault();
			e.stopPropagation();
		}

		function onPointerMove(e) {
			if (!dragging || dragging.col !== col) return;
			const dx = (e.clientX - dragging.startX) * COL_DIR[col];
			apply({ ...widths, [col]: dragging.start + dx }, false);
		}

		function onPointerUp() {
			if (!dragging || dragging.col !== col) return;
			dragging = null;
			root.classList.remove("is-col-resizing");
			setWidths(widths);
		}

		function onKeyDown(e) {
			const step = (e.shiftKey ? STEP_LARGE : STEP) * COL_DIR[col];
			if (e.key === "ArrowLeft") {
				e.preventDefault();
				apply({ ...widths, [col]: widths[col] - step }, true);
			} else if (e.key === "ArrowRight") {
				e.preventDefault();
				apply({ ...widths, [col]: widths[col] + step }, true);
			} else if (e.key === "Home") {
				e.preventDefault();
				resetAll();
			}
		}

		function onDblClick(e) {
			e.preventDefault();
			resetColumn(col);
		}

		splitter.addEventListener("pointerdown", onPointerDown);
		splitter.addEventListener("pointermove", onPointerMove);
		splitter.addEventListener("pointerup", onPointerUp);
		splitter.addEventListener("pointercancel", onPointerUp);
		splitter.addEventListener("keydown", onKeyDown);
		splitter.addEventListener("dblclick", onDblClick);
		cleanups.push(() => {
			splitter.removeEventListener("pointerdown", onPointerDown);
			splitter.removeEventListener("pointermove", onPointerMove);
			splitter.removeEventListener("pointerup", onPointerUp);
			splitter.removeEventListener("pointercancel", onPointerUp);
			splitter.removeEventListener("keydown", onKeyDown);
			splitter.removeEventListener("dblclick", onDblClick);
		});
	}

	const observer = new ResizeObserver(() => apply(widths, false));
	observer.observe(root);
	apply(widths, false);

	return () => {
		dragging = null;
		root.classList.remove("is-col-resizing");
		observer.disconnect();
		for (const fn of cleanups) fn();
	};
}
