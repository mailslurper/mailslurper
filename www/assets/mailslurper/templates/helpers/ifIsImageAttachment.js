define(
	[
		"hbs/handlebars"
	],
	function(Handlebars) {
		"use strict";

		var validImageMIMETypes = [
			"image/jpg",
			"image/jpeg",
			"image/png",
			"image/gif"
		];

		var helper = function(attachment, block) {
			if (validImageMIMETypes.indexOf(attachment.headers.contentType) > -1) {
				return block.fn(this);
			} else {
				return block.inverse(this);
			}
		};

		Handlebars.registerHelper("ifIsImageAttachment", helper);
		return helper;
	}
);
