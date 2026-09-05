import { el, clear } from "../util/dom.js";
import { getMailMessageURL, getMailMessageRawURL } from "../api.js";
import { attachmentURL, formatSize, isImageAttachment } from "../util/attachments.js";
import { formatDateTime } from "../util/date.js";
import { openImage } from "../components/attachmentLightbox.js";
import * as store from "../store.js";

function openMailLinksInNewTab(iframe) {
	iframe.addEventListener(
		"load",
		() => {
			const doc = iframe.contentDocument;
			if (!doc) return;
			for (const a of doc.querySelectorAll("a[href]")) {
				a.target = "_blank";
				const rel = new Set((a.rel || "").split(/\s+/).filter(Boolean));
				rel.add("noopener");
				rel.add("noreferrer");
				a.rel = [...rel].join(" ");
			}
		},
		{ once: true },
	);
}

export function renderEmpty(container) {
	clear(container);
	container.appendChild(el("div", { class: "empty-state" }, "Select a message to view it here."));
}

export function renderDetail(container, item) {
	clear(container);

	const detail = el("div", { class: "mail-detail" });
	detail.appendChild(el("h1", { class: "mail-detail-subject" }, item.subject || "(no subject)"));

	const headers = el("dl", { class: "mail-headers" });
	const addHeader = (label, value, extraClass) => {
		headers.appendChild(el("dt", {}, label));
		headers.appendChild(el("dd", extraClass ? { class: extraClass } : {}, value || "—"));
	};
	addHeader("From", item.fromAddress);
	addHeader("To", (item.toAddresses || []).join(", "));
	addHeader("Date", formatDateTime(item.dateSent, store.getDateFormat()), "mail-date");
	if (item.xMailer) addHeader("X-Mailer", item.xMailer);
	detail.appendChild(headers);

	const links = el("p", { class: "mail-detail-links" });
	links.appendChild(el("a", { href: getMailMessageURL(item.id), target: "_blank", rel: "noopener" }, "Open message"));
	links.appendChild(el("span", { class: "sep" }, "·"));
	links.appendChild(el("a", { href: getMailMessageRawURL(item.id), target: "_blank", rel: "noopener" }, "View raw source"));
	detail.appendChild(links);

	if (item.attachments && item.attachments.length > 0) {
		const list = el("ul", { class: "attachment-list" });
		for (const att of item.attachments) {
			const li = el("li");
			const url = attachmentURL(item.id, att.id);
			if (isImageAttachment(att)) {
				const btn = el("button", { class: "btn", type: "button" }, `${att.fileName} (${formatSize(att.sizeBytes)})`);
				btn.addEventListener("click", () => openImage(url, att.fileName));
				li.appendChild(btn);
			} else {
				li.appendChild(el("a", { href: url }, `${att.fileName} (${formatSize(att.sizeBytes)})`));
			}
			list.appendChild(li);
		}
		detail.appendChild(list);
	}

	detail.appendChild(el("div", { class: "mail-detail-divider" }));

	if (item.htmlBody) {
		const iframe = el("iframe", {
			class: "mail-body-frame",
			sandbox: "allow-same-origin allow-popups allow-popups-to-escape-sandbox",
			title: "Message body",
		});
		detail.appendChild(iframe);
		openMailLinksInNewTab(iframe);
		iframe.srcdoc = item.htmlBody;
	} else {
		detail.appendChild(el("pre", { class: "mail-body-text" }, item.textBody || ""));
	}

	container.appendChild(detail);
}
