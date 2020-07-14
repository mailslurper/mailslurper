// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

(function () {

	var savedSearchesModalTemplate;
	var saveSearchModalTemplate;

	window.TemplateService.load("savedSearchesModal").then(function (template) {
		savedSearchesModalTemplate = Handlebars.compile(template);
	});

	window.TemplateService.load("saveSearchModal").then(function (template) {
		saveSearchModalTemplate = Handlebars.compile(template);
	});

	window.SavedSearchesWidget = {
		showPicker: function (callback) {
			var dialogRef = BootstrapDialog.show({
				title: "Saved Searches",
				message: savedSearchesModalTemplate({ savedSearches: window.SettingsService.retrieveSavedSearches() }),
				closable: true,
				nl2br: false,
				buttons: [
					{
						id: "btnManageSavedSearches",
						label: "Manage",
						cssClass: "btn-default",
						action: function () {
							var appURL = window.SettingsService.getAppURL();
							window.location = appURL + "/savedsearches";
						}
					},
					{
						id: "btnCancelSavedSearch",
						label: "Cancel",
						cssClass: "btn-default",
						action: function (dialogRef) {
							dialogRef.close();
						}
					},
					{
						id: "btnOK",
						label: "OK",
						cssClass: "btn-primary",
						action: function () {
							var savedSearch = window.SettingsService.getSavedSearchByIndex(window.parseInt($("#savedSearchID option:selected").val(), 10));
							callback(savedSearch);
							dialogRef.close();
						}
					}
				]
			});
		},

		showSaveSearchModal: function (callback) {
			var dialogRef = BootstrapDialog.show({
				title: "Save Search",
				message: saveSearchModalTemplate(),
				closable: true,
				nl2br: false,
				buttons: [
					{
						id: "btnCancelSaveSearch",
						label: "Cancel",
						cssClass: "btn-default",
						action: function (dialogRef) {
							dialogRef.close();
						}
					},
					{
						id: "btnSaveSearchOK",
						label: "OK",
						cssClass: "btn-primary",
						action: function (dialogRef) {
							var saveSearchName = $("#txtSaveSearchName").val();
							if (saveSearchName.length <= 0) {
								alert("Please enter a name for your search!");
							} else {
								dialogRef.close();
								callback(saveSearchName);
							}
						}
					}
				],
				onshown: function (dialogRef) {
					$("#txtSaveSearchName").focus();
				}
			});
		}
	};
}());