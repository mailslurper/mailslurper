// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
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
