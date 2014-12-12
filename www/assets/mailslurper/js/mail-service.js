"use strict";

var MailService = {
	formatMailDate: function(dateString) {
		return moment(dateString).format("MMMM Do YYYY, h:mm:ss a");
	},

	/**
	 * Performs an AJAX call to get a single mail item by ID.
	 */
	getMailItem: function(id, callback) {
		$.ajax({ url: ServiceSettings.buildUrl("/mails/" + id ) }).done(function(data) {
			callback(data.mailItem);
		});
	},

	/**
	 * Performs an AJAX call to retrieve a list of mail items
	 * by page. The returned mail items are passed to a
	 * supplied callback function.
	 */
	loadMailItems: function(page, callback) {
		$.ajax({ url: ServiceSettings.buildUrl("/mails/page/" + page) }).done(function(data) {
			callback(data.mailItems);
		});
	}
};
