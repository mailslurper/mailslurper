define(
	[
		"handlebarsLibrary",
		"moment"
	],
	function(Handlebars, moment) {
		"use strict";

		/*
		 * Provides a helper to format a date in the format of
		 * Sunday, January 1st 2015 @ 05:00 PM. Usage is like:
		 * {{formatDateTime someDate}}
		 */
		Handlebars.registerHelper("formatDateTime", function(date) {
			return moment(date).format("YYYY-MM-DD hh:mm A");
		});

		/*
		 * Provides a helper to format a date in the format of
		 * January 8th. Usage is like: {{formatMonthDay someDate}}
		 */
		Handlebars.registerHelper("formatMonthDay", function(date) {
			return moment(date).format("MMMM Do");
		});

		return Handlebars;
	}
);
