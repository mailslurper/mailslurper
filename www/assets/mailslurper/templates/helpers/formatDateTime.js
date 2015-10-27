define(
	[
		"hbs/handlebars",
		"moment",
		"services/SettingsService"
	],
	function(Handlebars, moment, settingsService) {
		"use strict";

		var helper = function(date) {
			var dateFormat = settingsService.retrieveSettings().dateFormat;
			return moment(date).format(dateFormat);
		};

		Handlebars.registerHelper("formatDateTime", helper);
		return helper;
	}
);
