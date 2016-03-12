// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

define(
	[
		"jquery"
	],
	function($) {
		"use strict";

		var service = {
			getAutoRefreshOptions: function() {
				return [
					{ value: 0, description: "Don't auto refresh" },
					{ value: 1, description: "Every minute" },
					{ value: 5, description: "Every 5 minutes" },
					{ value: 10, description: "Every 10 minutes" }
				];
			},

			getDateFormatOptions: function() {
				return [
					{ dateFormat: "YYYY-MM-DD HH:mm", description: "International" },
					{ dateFormat: "MM/DD/YYYY hh:mm A", description: "US" }
				];
			},

			/**
			 * getPruneOptions returns email pruning options. This will place the array of
			 * options in the context with a key of "pruneOptions".
			 */
			getPruneOptions: function(serviceURL) {
				return $.ajax({
					url: serviceURL + "/pruneoptions",
					method: "GET"
				});
			},

			/**
			 * validatePruneCode will reject the promise if the prune code
			 * is invalid. This function expects that the context to have
			 * a key named "pruneOptions" to validate against. Basically
			 * an AJAX call to get prune options should have occurred
			 * prior to this. It compare it to a key named "pruneCode".
			 */
			validatePruneCode: function(pruneOptions, pruneCode) {
				for (var index = 0; index < pruneOptions.length; index++) {
					if (pruneOptions[index].pruneCode === pruneCode) {
						return true;
					}
				}

				return false;
			}
		};

		return service;
	}
);
