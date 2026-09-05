// Pure date formatting helpers. dateFormat is one of "locale" (browser
// default), "iso" (ISO 8601), or a fixed set of common presentations.

export function formatDateTime(isoString, dateFormat) {
	if (!isoString) return "";
	const d = new Date(isoString);
	if (Number.isNaN(d.getTime())) return isoString;

	switch (dateFormat) {
		case "iso":
			return d.toISOString();
		case "us":
			return d.toLocaleString("en-US");
		case "eu":
			return d.toLocaleString("en-GB");
		default:
			return d.toLocaleString();
	}
}

// Returns a human countdown string like "2:05" for a given number of
// milliseconds remaining. Zero or negative returns "0:00".
export function formatCountdown(msRemaining) {
	const totalSeconds = Math.max(0, Math.floor(msRemaining / 1000));
	const minutes = Math.floor(totalSeconds / 60);
	const seconds = totalSeconds % 60;
	return `${minutes}:${String(seconds).padStart(2, "0")}`;
}
