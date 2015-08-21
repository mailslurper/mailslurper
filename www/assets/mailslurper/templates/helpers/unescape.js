define(
	[
		"hbs/handlebars"
	],
	function(Handlebars) {
		"use strict";

		var replacements = [
			{ encoded: "&#39;", replacement: "'" },
			{ encoded: "&amp;", replacement: "&" },
			{ encoded: "&lt;", replacement: "<" },
			{ encoded: "&gt;", replacement: ">" },
			{ encoded: "&quot;", replacement: "\"" },
			{ encoded: "&#96;", replacement: "`" }
		];

		var helper = function(escapedString) {
			for (var index = 0; index < replacements.length; index++) {
				escapedString = escapedString.replace(replacements[index].encoded, replacements[index].replacement);
			}

			return escapedString;
		};

		Handlebars.registerHelper("unescape", helper);
		return helper;
	}
);
