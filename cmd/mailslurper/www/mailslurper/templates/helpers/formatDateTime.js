// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
"use strict";

Handlebars.registerHelper("formatDateTime", function (date) {
	var dateFormat = window.SettingsService.retrieveSettings().dateFormat;
	return moment(date).format(dateFormat);
});
