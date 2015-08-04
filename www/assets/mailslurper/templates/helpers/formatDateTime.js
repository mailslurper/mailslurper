define(
	[
		"hbs/handlebars",
		"moment"
	],
	function(Handlebars, moment) {
		"use strict";

		var helper = function(date) {
			return moment(date).format("YYYY-MM-DD hh:mm A");
		};

		Handlebars.registerHelper("formatDateTime", helper);
		return helper;
	}
);
