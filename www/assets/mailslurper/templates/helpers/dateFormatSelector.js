define(
	[
		"hbs/handlebars",
		"services/SeedService"
	],
	function(Handlebars, SeedService) {
		"use strict";

		var helper = function(elementName, selectedDateFormat) {
			var html = "<select id=\"" + elementName + "\" class=\"form-control\">";
			var dateFormatOptions = SeedService.getDateFormatOptions();

			for (var index = 0; index < dateFormatOptions.length; index++) {
				html += "<option value=\"" + dateFormatOptions[index].dateFormat + "\"";
				html += (selectedDateFormat === dateFormatOptions[index].dateFormat) ? " selected=\"selected\"" : "";
				html += ">";
				html += dateFormatOptions[index].description;
				html += "</option>";
			}

			html += "</select>";

			return html;
		};

		Handlebars.registerHelper("dateFormatSelector", helper);
		return helper;
	}
);
