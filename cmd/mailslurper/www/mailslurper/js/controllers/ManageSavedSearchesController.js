// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

(function () {
	"use strict";

	function deleteSavedSearch(savedSearchIndex) {
		window.SettingsService.deleteSavedSearch(savedSearchIndex);
	}

	function getSavedSearches() {
		var savedSearches = window.SettingsService.retrieveSavedSearches();
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
	}

	function initialize() {
		$(".deleteSavedSearch").on("click", function () {
			var savedSearchIndex = window.parseInt($(this).attr("data-index"), 10);
			deleteSavedSearch(savedSearchIndex);

			renderSavedSearches(getSavedSearches());
			initialize();
		});
	}

	function loadManageSavedSearchesTemplate() {
		return new Promise(function (resolve, reject) {
			window.TemplateService.load("manageSavedSearches")
				.then(function (template) {
					manageSavedSearchesTemplate = Handlebars.compile(template);
					return resolve();
				})
				.catch(function (err) {
					return reject(err);
				});
		});
	}

	function renderSavedSearches(savedSearches) {
		$("#savedSearches").html(manageSavedSearchesTemplate({ savedSearches: savedSearches }));
	}

	/****************************************************************************
	 * Constructor
	 ***************************************************************************/
	var manageSavedSearchesTemplate;

	window.AlertService.block("Loading...");

	loadManageSavedSearchesTemplate()
		.then(function () {
			renderSavedSearches(getSavedSearches());
			initialize();
			window.AlertService.unblock();
		})
		.catch(function (err) {
			window.AlertService.error(err);
		});
}());