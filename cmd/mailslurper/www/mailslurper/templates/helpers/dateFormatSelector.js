// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

Handlebars.registerHelper("dateFormatSelector", function (elementName, selectedDateFormat) {
	var html = "<select id=\"" + elementName + "\" class=\"form-control\">";
	var dateFormatOptions = window.SeedService.getDateFormatOptions();

	for (var index = 0; index < dateFormatOptions.length; index++) {
		html += "<option value=\"" + dateFormatOptions[index].dateFormat + "\"";
		html += (selectedDateFormat === dateFormatOptions[index].dateFormat) ? " selected=\"selected\"" : "";
		html += ">";
		html += dateFormatOptions[index].description;
		html += "</option>";
	}

	html += "</select>";

	return html;
});
