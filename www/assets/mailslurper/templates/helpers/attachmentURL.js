define(
	[
		"hbs/handlebars",
		"services/SettingsService"
	],
	function(Handlebars, SettingsService) {
		"use strict";

		var helper = function(attachment) {
			return SettingsService.getServiceURLNow() + "/mail/" + attachment.mailId + "/attachment/" + attachment.id;
		};

		Handlebars.registerHelper("attachmentURL", helper);
		return helper;
	}
);
