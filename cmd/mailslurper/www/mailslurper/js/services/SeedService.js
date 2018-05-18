// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

window.SeedService = {
	getAutoRefreshOptions: function () {
		return [
			{ value: 0, description: "Don't auto refresh" },
			{ value: 1, description: "Every minute" },
			{ value: 5, description: "Every 5 minutes" },
			{ value: 10, description: "Every 10 minutes" }
		];
	},

	getDateFormatOptions: function () {
		return [
			{ dateFormat: "YYYY-MM-DD HH:mm", description: "International" },
			{ dateFormat: "MM/DD/YYYY hh:mm A", description: "US" }
		];
	},

	/**
	 * getPruneOptions returns email pruning options. This will place the array of
	 * options in the context with a key of "pruneOptions".
	 */
	getPruneOptions: function (serviceURL) {
		return new Promise(function (resolve, reject) {
			$.ajax(window.AuthService.decorateRequestWithAuthorization({
				url: serviceURL + "/pruneoptions",
				method: "GET"
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
	 * getThemes returns themes
	 */
	getThemes: function () {
		return [
			{ name: "Default", "theme": "default" },
			{ name: "Lumen", "theme": "lumen" },
			{ name: "Readable", "theme": "readable" },
			{ name: "Slate", "theme": "slate" },
			{ name: "SpaceLab", "theme": "spacelab" }
		];
	},

	/**
	 * validatePruneCode will reject the promise if the prune code
	 * is invalid. This function expects that the context to have
	 * a key named "pruneOptions" to validate against. Basically
	 * an AJAX call to get prune options should have occurred
	 * prior to this. It compare it to a key named "pruneCode".
	 */
	validatePruneCode: function (pruneOptions, pruneCode) {
		for (var index = 0; index < pruneOptions.length; index++) {
			if (pruneOptions[index].pruneCode === pruneCode) {
				return true;
			}
		}

		return false;
	}
};
