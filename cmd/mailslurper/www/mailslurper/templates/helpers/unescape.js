// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
"use strict";

Handlebars.registerHelper("unescape", function (escapedString) {
	var replacements = [
		{ encoded: /&#39;/g, replacement: "'" },
		{ encoded: /&amp;/g, replacement: "&" },
		{ encoded: /&lt;/g, replacement: "<" },
		{ encoded: /&gt;/g, replacement: ">" },
		{ encoded: /&quot;/g, replacement: "\"" },
		{ encoded: /&#96;/g, replacement: "`" },
		{ encoded: /&#34;/g, replacement: "\"" }
	];

	if (escapedString === undefined) {
		return "";
	}

	for (var index = 0; index < replacements.length; index++) {
		escapedString = escapedString.replace(replacements[index].encoded, replacements[index].replacement);
	}

	return escapedString;
});
