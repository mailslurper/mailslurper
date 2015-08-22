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
			}
		};

		return service;
	}
)
