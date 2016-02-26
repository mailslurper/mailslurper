// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
			deleteMailItems: function(serviceURL, pruneCode) {
				return $.ajax({
					method: "DELETE",
					url: serviceURL + "/mail",
					data: JSON.stringify({
						pruneCode: pruneCode
					})
				});
			},

			/**
			 * getAttachment retrieves a specified attachment from a given mail ID.
			 * Context is expected to have "mailID" and "attachmentID". The context
			 * will store "attachment" when the promise is fullfilled.
			 */
			getAttachment: function(serviceURL, mailID, attachmentID) {
				return $.ajax({
					method: "GET",
					url: serviceURL + "/mail/" + mailID + "/attachment/" + attachmentID
				});
			},

			/**
			 * getMailByID returns a specific mail item. The context must contain
			 * a key named "mailID" and will return a key named "mail".
			 */
			getMailByID: function(serviceURL, mailID) {
				return $.ajax({
					method: "GET",
					url: serviceURL + "/mail/" + mailID
				});
			},

			/**
			 * getMailCount returns the number of mail items in storage. This will put
			 * the count into a key named "mailCount" in the context object.
			 */
			getMailCount: function(serviceURL) {
				return $.ajax({
					method: "GET",
					url: serviceURL + "/mailcount"
				});
			},

			/**
			 * getMails returns a page of stored email. The page number must be a key
			 * named "page" in the context object. This will return mail items as an
			 * array in a key named "mails" in the context object.
			 */
			getMails: function(serviceURL, page, searchCriteria) {
				var url = serviceURL + "/mail?pageNumber=" + page;

				if (searchCriteria.message != "") {
					url += "&message=" + searchCriteria.searchMessage;
				}

				if (searchCriteria.searchStart) {
					url += "&start=" + searchCriteria.searchStart.format("YYYY-MM-DD");
				}

				if (searchCriteria.searchEnd) {
					url += "&end=" + searchCriteria.searchEnd.format("YYYY-MM-DD");
				}

				if (searchCriteria.searchFrom) {
					url += "&from=" + searchCriteria.searchFrom;
				}

				if (searchCriteria.searchTo) {
					url += "&to=" + searchCriteria.searchTo;
				}

				return $.ajax({
					method: "GET",
					url: url
				});
			}
		};

		return service;
	}
);
