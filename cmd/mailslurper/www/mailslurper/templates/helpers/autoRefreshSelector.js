// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
"use strict";

Handlebars.registerHelper("autoRefreshSelector", function (elementName, selectedValue) {
	var html = "<select id=\"" + elementName + "\" class=\"form-control\">";
	var options = window.SeedService.getAutoRefreshOptions();

	for (var index = 0; index < options.length; index++) {
		html += "<option value=\"" + options[index].value + "\"";
		html += (selectedValue === options[index].value) ? " selected=\"selected\"" : "";
		html += ">";
		html += options[index].description;
		html += "</option>";
	}

	html += "</select>";

	return html;
});
