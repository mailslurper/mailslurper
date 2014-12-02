"use strict";

$(document).ready(function() {
	/*
	 * Setup routing
	 */
	var routes = {
		"home": MailHomeController.index
	};

	Grapnel.listen(routes);
	window.location.hash = "#home";
});
