// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

define(
	[
		"jquery",
		"services/SettingsService"
	],
	function($, SettingsService) {
		"use strict";

		var service = {
			/**
			 * applySavedTheme will load the user's selected theme and apply it
			 */
			applySavedTheme: function() {
				var settings = SettingsService.retrieveSettings();
				$("#themeBootstrapStylesheet").attr("href", "/www/mailslurper/themes/" + settings.theme + "/bootstrap.css");
				$("#themeStylesheet").attr("href", "/www/mailslurper/themes/" + settings.theme + "/style.css");
			}
		};

		return service;
	}
);
