const connectionListeners = new Set();

let source = null;
let reconnectTimer = null;
let reconnectDelay = 1000;
let eventHandler = null;
let stopped = false;

const MAX_RECONNECT_DELAY = 30_000;

function notifyConnection(connected) {
	for (const fn of connectionListeners) fn(connected);
}

function scheduleReconnect() {
	if (stopped || reconnectTimer) return;
	reconnectTimer = setTimeout(() => {
		reconnectTimer = null;
		if (!stopped) connect(eventHandler);
	}, reconnectDelay);
	reconnectDelay = Math.min(reconnectDelay * 2, MAX_RECONNECT_DELAY);
}

function handleOpen() {
	reconnectDelay = 1000;
	notifyConnection(true);
}

function handleError() {
	notifyConnection(false);
	if (source) {
		source.close();
		source = null;
	}
	scheduleReconnect();
}

/**
 * Opens (or reopens) the SSE stream and forwards parsed events to onEvent.
 * Returns a disconnect function.
 */
export function connect(onEvent) {
	disconnect();
	stopped = false;
	eventHandler = onEvent;

	source = new EventSource("/api/events");

	source.addEventListener("mail.received", (e) => {
		try {
			onEvent("mail.received", JSON.parse(e.data));
		} catch {
			// ignore malformed payloads
		}
	});

	source.addEventListener("mail.pruned", (e) => {
		try {
			onEvent("mail.pruned", JSON.parse(e.data));
		} catch {
			// ignore malformed payloads
		}
	});

	source.onopen = handleOpen;
	source.onerror = handleError;

	return disconnect;
}

/** Stops the SSE stream and cancels pending reconnects. */
export function disconnect() {
	stopped = true;
	if (reconnectTimer) {
		clearTimeout(reconnectTimer);
		reconnectTimer = null;
	}
	if (source) {
		source.close();
		source = null;
	}
	notifyConnection(false);
}

/** Subscribe to connection state changes (true = live, false = disconnected). */
export function onConnectionChange(fn) {
	connectionListeners.add(fn);
	return () => connectionListeners.delete(fn);
}
