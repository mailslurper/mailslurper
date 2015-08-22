define(
	[
		"jquery",
		"moment",
		"services/SettingsService"
	],
	function($, moment, SettingsService) {
		"use strict";

		var service = {
			/**
			 * deleteMailItems deletes a set of mail items. The criteria is defined
			 * by a "pruneCode", which is one of:
			 *    * 60plus
			 *    * 30plus
			 *    * 2wksplus
			 *    * all
			 */
			deleteMailItems: function(context) {
				return new Promise(function(resolve, reject) {
					$.ajax({
						method: "DELETE",
						url: context.serviceURL + "/mail",
						data: JSON.stringify({
							pruneCode: context.pruneCode
						})
					}).then(
						function() {
							resolve(context);
						},

						function(err) {
							reject(err);
						}
					);
				});
			},

			/**
			 * getAttachment retrieves a specified attachment from a given mail ID.
			 * Context is expected to have "mailID" and "attachmentID". The context
			 * will store "attachment" when the promise is fullfilled.
			 */
			getAttachment: function(context) {
				return new Promise(function(resolve, reject) {
					$.ajax({
						method: "GET",
						url: context.serviceURL + "/mail/" + context.mailID + "/attachment/" + context.attachmentID
					}).then(
						function(response) {
							context.attachment = response;
							resolve(context);
						},

						function(error) {
							reject(error);
						}
					)
				});
			},

			/**
			 * getMailByID returns a specific mail item. The context must contain
			 * a key named "mailID" and will return a key named "mail".
			 */
			getMailByID: function(context) {
				return new Promise(function(resolve, reject) {
					$.ajax({
						method: "GET",
						url: context.serviceURL + "/mail/" + context.mailID
					}).then(
						function(response) {
							context.mail = response;
							resolve(context);
						},

						function(error) {
							reject(error);
						}
					);
				});
			},

			/**
			 * getMailCount returns the number of mail items in storage. This will put
			 * the count into a key named "mailCount" in the context object.
			 */
			getMailCount: function(context) {
				return new Promise(function(resolve, reject) {
					$.ajax({
						method: "GET",
						url: context.serviceURL + "/mailcount"
					}).then(
						function(response) {
							context.mailCount = response.mailCount;
							resolve(context);
						},

						function(error) {
							reject(error);
						}
					);
				});
			},

			/**
			 * getMails returns a page of stored email. The page number must be a key
			 * named "page" in the context object. This will return mail items as an
			 * array in a key named "mails" in the context object.
			 */
			getMails: function(context) {
				return new Promise(function(resolve, reject) {
					var url = context.serviceURL + "/mails/" + context.page + "?";
					url += "message=" + (context.searchMessage || "");

					if (context.searchStart) {
						url += "&start=" + context.searchStart.format("YYYY-MM-DD");
					}

					if (context.searchEnd) {
						url += "&end=" + context.searchEnd.format("YYYY-MM-DD");
					}

					if (context.searchFrom) {
						url += "&from=" + context.searchFrom;
					}

					if (context.searchTo) {
						url += "&to=" + context.searchTo;
					}

					$.ajax({
						method: "GET",
						url: url
					}).then(
						function(response, status, xhr) {
							context.mails = response.mailItems;
							context.totalPages = window.parseInt(xhr.getResponseHeader("X-Total-Pages"), 10);
							context.totalMailCount = window.parseInt(xhr.getResponseHeader("X-Total-Mail-Count"), 10);

							resolve(context);
						},

						function(error) {
							context.message = error;
							reject(context);
						}
					);
				});
			}
		};

		return service;
	}
);
