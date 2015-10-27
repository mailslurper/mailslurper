define(
	[
		"jquery"
	],
	function($) {
		"use strict";

		var service = {
			getDateFormatOptions: function() {
				return [
					{ dateFormat: "YYYY-MM-DD", description: "International" },
					{ dateFormat: "MM/DD/YYYY", description: "US" }
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
