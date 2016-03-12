// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
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
