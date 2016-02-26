// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
		settingsService,
		alertService,
		Dialog,
		manageSavedSearchesTemplate
	) {
		"use strict";

		var deleteSavedSearch = function(savedSearchIndex) {
			settingsService.deleteSavedSearch(savedSearchIndex);
		};

		var getSavedSearches = function() {
			var savedSearches = settingsService.retrieveSavedSearches();
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

			return grouped;
		};

		var initialize = function() {
			$(".deleteSavedSearch").on("click", function() {
				var savedSearchIndex = window.parseInt($(this).attr("data-index"), 10);
				deleteSavedSearch(savedSearchIndex);

				renderSavedSearches(getSavedSearches());
				initialize();
			});
		};

		var renderSavedSearches = function(savedSearches) {
			$("#savedSearches").html(manageSavedSearchesTemplate({ savedSearches: savedSearches }));
		};

		/****************************************************************************
		 * Constructor
		 ***************************************************************************/
		alertService.block("Loading...");
		renderSavedSearches(getSavedSearches());
		initialize();

		alertService.unblock();
	}
);
