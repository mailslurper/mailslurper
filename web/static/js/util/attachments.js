const IMAGE_TYPES = new Set(["image/png", "image/jpeg", "image/gif", "image/webp", "image/svg+xml", "image/bmp"]);

export function isImageAttachment(attachment) {
	return IMAGE_TYPES.has((attachment.contentType || "").toLowerCase());
}

export function attachmentURL(mailId, attachmentId) {
	return `/api/mail/${encodeURIComponent(mailId)}/attachments/${encodeURIComponent(attachmentId)}`;
}

export function formatSize(bytes) {
	if (bytes < 1024) return `${bytes} B`;
	if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
	return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
