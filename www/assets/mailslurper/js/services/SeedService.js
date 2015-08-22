define(
	[
		"jquery"
	],
	function($) {
		"use strict";

		var service = {
			/**
			 * getPruneOptions returns email pruning options. This will place the array of
			 * options in the context with a key of "pruneOptions".
			 */
			getPruneOptions: function(context) {
				return new Promise(function(resolve, reject) {
					$.ajax({
						url: context.serviceURL + "/pruneoptions",
						method: "GET"
					}).then(
						function(pruneOptions) {
							context.pruneOptions = pruneOptions;
							resolve(context);
						},

						function(err) {
							reject(err);
						}
					);
				});
			},

			/**
			 * validatePruneCode will reject the promise if the prune code
			 * is invalid. This function expects that the context to have
			 * a key named "pruneOptions" to validate against. Basically
			 * an AJAX call to get prune options should have occurred
			 * prior to this. It compare it to a key named "pruneCode".
			 */
			validatePruneCode: function(context) {
				if (!context.hasOwnProperty("pruneOptions")) {
					context.message = "Prune codes not retrieved!";
					return Promise.reject(context);
				} else {
					for (var index = 0; index < context.pruneOptions.length; index++) {
						if (context.pruneOptions[index].pruneCode === context.pruneCode) {
							return Promise.resolve(context);
						}
					}

					context.message = "Invalid option for pruning emails";
					return Promise.reject(context);
				}
			}
		};

		return service;
	}
)
