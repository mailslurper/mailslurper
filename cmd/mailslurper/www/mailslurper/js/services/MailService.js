// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

window.MailService = {
	/**
	 * deleteMailItems deletes a set of mail items. The criteria is defined
	 * by a "pruneCode", which is one of:
	 *    * 60plus
	 *    * 30plus
	 *    * 2wksplus
	 *    * all
	 */
	deleteMailItems: function (serviceURL, pruneCode) {
		return new Promise(function (resolve, reject) {
			$.ajax(window.AuthService.decorateRequestWithAuthorization({
				method: "DELETE",
				url: serviceURL + "/mail",
				contentType: "application/json; charset=utf-8",
				dataType: "text",
				data: JSON.stringify({
					pruneCode: pruneCode
				})
			})).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	},

	/**
	 * getAttachment retrieves a specified attachment from a given mail ID.
	 * Context is expected to have "mailID" and "attachmentID". The context
	 * will store "attachment" when the promise is fullfilled.
	 */
	getAttachment: function (serviceURL, mailID, attachmentID) {
		return new Promise(function (resolve, reject) {
			$.ajax(window.AuthService.decorateRequestWithAuthorization({
				method: "GET",
				url: serviceURL + "/mail/" + mailID + "/attachment/" + attachmentID
			})).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	},

	/**
	 * getMailByID returns a specific mail item. The context must contain
	 * a key named "mailID" and will return a key named "mail".
	 */
	getMailByID: function (serviceURL, mailID) {
		return new Promise(function (resolve, reject) {
			$.ajax(window.AuthService.decorateRequestWithAuthorization({
				method: "GET",
				url: serviceURL + "/mail/" + mailID
			})).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	},

	/**
	 * getMailCount returns the number of mail items in storage. This will put
	 * the count into a key named "mailCount" in the context object.
	 */
	getMailCount: function (serviceURL) {
		return new Promise(function (resolve, reject) {
			$.ajax(window.AuthService.decorateRequestWithAuthorization({
				method: "GET",
				url: serviceURL + "/mailcount",
				cache: false
			})).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	},

	/**
	 * getMailMessageURL returns the full service URL to get a mail's message body
	 */
	getMailMessageURL: function (serviceURL, mailID) {
		return serviceURL + "/mail/" + mailID + "/message";
	},

	/**
	 * getMails returns a page of stored email. The page number must be a key
	 * named "page" in the context object. This will return mail items as an
	 * array in a key named "mails" in the context object.
	 */
	getMails: function (serviceURL, page, searchCriteria, sortCriteria) {
		return new Promise(function (resolve, reject) {
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

			if (sortCriteria.orderByField) {
				url += "&orderby=" + sortCriteria.orderByField;
			}

			if (sortCriteria.orderByDirection) {
				url += "&dir=" + sortCriteria.orderByDirection;
			}

			var params = window.AuthService.decorateRequestWithAuthorization({
				method: "GET",
				url: url,
				cache: false
			});

			$.ajax(params).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	}
};
