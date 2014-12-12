"use strict";

$(document).ready(function() {
	/*
	 * Get service app settings
	 */
	ServiceSettings.getServiceSettings();

	/*
	 * Setup routing
	 */
	var routes = {
		"home": MailHomeController.index
	};

	Grapnel.listen(routes);
	window.location.hash = "#home";
});
