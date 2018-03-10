// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
"use strict";

Handlebars.registerHelper("ifIsImageAttachment", function (attachment, block) {
	var validImageMIMETypes = [
		"image/jpg",
		"image/jpeg",
		"image/png",
		"image/gif"
	];

	if (validImageMIMETypes.indexOf(attachment.headers.contentType) > -1) {
		return block.fn(this);
	} else {
		return block.inverse(this);
	}
});
