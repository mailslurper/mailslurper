// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
define(
	[
		"hbs/handlebars",
		"services/SeedService"
	],
	function(Handlebars, seedService) {
		"use strict";

		var helper = function(elementName, selectedTheme) {
			var html = "<select id=\"" + elementName + "\" class=\"form-control\">";
			var themes = seedService.getThemes();

			for (var index = 0; index < themes.length; index++) {
				html += "<option value=\"" + themes[index].theme + "\"";
				html += (selectedTheme === themes[index].theme) ? " selected=\"selected\"" : "";
				html += ">";
				html += themes[index].name;
				html += "</option>";
			}

			html += "</select>";

			return html;
		};

		Handlebars.registerHelper("themeSelector", helper);
		return helper;
	}
);
