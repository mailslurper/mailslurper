// Pure helpers converting inbox search/sort/paging state to and from a URL
// query string, used both for the API request and for the address bar hash.

export function stateToQueryString(state) {
	const params = new URLSearchParams();
	if (state.q) params.set("q", state.q);
	if (state.from) params.set("from", state.from);
	if (state.to) params.set("to", state.to);
	if (state.start) params.set("start", state.start);
	if (state.end) params.set("end", state.end);
	if (state.sort) params.set("sort", state.sort);
	if (state.dir) params.set("dir", state.dir);
	if (state.page && state.page !== 1) params.set("page", String(state.page));
	return params.toString();
}

export function queryStringToState(queryString) {
	const params = new URLSearchParams(queryString);
	const state = {
		q: params.get("q") || "",
		from: params.get("from") || "",
		to: params.get("to") || "",
		start: params.get("start") || "",
		end: params.get("end") || "",
		sort: params.get("sort") || "date",
		dir: params.get("dir") || "desc",
		page: parseInt(params.get("page") || "1", 10) || 1,
	};
	return state;
}
