// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
"use strict";

Handlebars.registerHelper("pageSelector", function (elementName, totalPages, currentPage) {
	var html = "<select id=\"" + elementName + "\" class=\"form-control\">";

	for (var index = 1; index <= totalPages; index++) {
		html += "<option value=\"" + index + "\"";
		html += (currentPage == index) ? " selected=\"selected\"" : "";
		html += ">" + index;
		html += "</option>";
	}

	html += "</select>";

	return html;
});
