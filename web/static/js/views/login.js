import * as auth from "../auth.js";
import { el } from "../util/dom.js";
import { toast } from "../components/toast.js";

export async function mount(container, params) {
	if (!(await auth.authRequired().catch(() => false))) {
		location.hash = params?.next || "#/inbox";
		return;
	}

	container.className = "view view-login";

	const card = el("div", { class: "login-card" });
	card.appendChild(el("h1", {}, "Sign in"));

	const form = el("form");
	const userField = el("div", { class: "field" });
	userField.appendChild(el("label", { for: "login-username" }, "Username"));
	const userInput = el("input", { id: "login-username", type: "text", autocomplete: "username", required: "required" });
	userField.appendChild(userInput);
	form.appendChild(userField);

	const passField = el("div", { class: "field" });
	passField.appendChild(el("label", { for: "login-password" }, "Password"));
	const passInput = el("input", { id: "login-password", type: "password", autocomplete: "current-password", required: "required" });
	passField.appendChild(passInput);
	form.appendChild(passField);

	const submitBtn = el("button", { class: "btn btn-primary", type: "submit" }, "Sign in");
	form.appendChild(submitBtn);

	form.addEventListener("submit", async (e) => {
		e.preventDefault();
		submitBtn.disabled = true;
		try {
			await auth.login(userInput.value, passInput.value);
			location.hash = params?.next || "#/inbox";
		} catch (err) {
			toast.error(err.message || "Login failed");
		} finally {
			submitBtn.disabled = false;
		}
	});

	card.appendChild(form);
	container.appendChild(card);
	userInput.focus();
}

export function unmount() {}
