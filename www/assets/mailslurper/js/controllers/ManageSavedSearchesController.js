require(
	[
		"jquery",
		"services/SettingsService",
		"services/AlertService",
		"bootstrap-dialog",
		"hbs!templates/manageSavedSearches"
	],
	function(
		$,
		SettingsService,
		AlertService,
		Dialog,
		manageSavedSearchesTemplate
	) {
		"use strict";

		var deleteSavedSearch = function(context) {
			SettingsService.deleteSavedSearch(context.savedSearchIndex);
			return Promise.resolve(context);
		};

		var getSavedSearches = function(context) {
			var savedSearches = SettingsService.retrieveSavedSearches();
			var grouped = [];

			for (var outerIndex = 0; outerIndex < savedSearches.length; outerIndex += 2) {
				var innerGroup = [];

				for (var innerIndex = 0; innerIndex < 2; innerIndex++) {
					var actualIndex = outerIndex + innerIndex;

					if (actualIndex < savedSearches.length) {
						innerGroup.push(savedSearches[actualIndex]);
					}
				}

				grouped.push(innerGroup);
			}

			context.savedSearches = grouped;
			return Promise.resolve(context);
		};

		var initialize = function(context) {
			$(".deleteSavedSearch").on("click", function() {
				context.savedSearchIndex = window.parseInt($(this).attr("data-index"), 10);
				deleteSavedSearch(context)
					.then(getSavedSearches)
					.then(renderSavedSearches)
					.then(initialize)
					.catch(AlertService.error);
			});

			return Promise.resolve(context);
		};

		var renderSavedSearches = function(context) {
			$("#savedSearches").html(manageSavedSearchesTemplate({ savedSearches: context.savedSearches }));
			return Promise.resolve(context);
		};

		/****************************************************************************
		 * Constructor
		 ***************************************************************************/
		var context = {
			savedSearches: [],
			message: "Loading"
		};

		AlertService.block(context)
			.then(getSavedSearches)
			.then(renderSavedSearches)
			.then(initialize)
			.then(AlertService.unblock)
			.catch(AlertService.error);

	}
);
