// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

window.SettingsService = {
	/**
	 * addSavedSearch will append a new saved search
	 */
	addSavedSearch: function (name, searchCriteria) {
		searchCriteria.name = name;
		var currentSavedSearches = window.SettingsService.retrieveSavedSearches();

		currentSavedSearches.push(searchCriteria);
		window.SettingsService.storeSavedSearches(currentSavedSearches);
	},

	deleteSavedSearch: function (index) {
		var currentSavedSearches = window.SettingsService.retrieveSavedSearches();

		if (index > -1 && index < currentSavedSearches.length) {
			currentSavedSearches.splice(index, 1);
			window.SettingsService.storeSavedSearches(currentSavedSearches);
		}
	},

	/**
	 * getSavedSearchByIndex returns a saved search based on its location
	 * in the array.
	 */
	getSavedSearchByIndex: function (index) {
		var searches = window.SettingsService.retrieveSavedSearches();

		if (index < 0 && index >= searches.length) {
			AlertService.error({ message: "There is no saved search by that ID" });
			return {
				name: "",
				searchMessage: ""
			};
		}

		return searches[index];
	},

	/**
	 * getServiceSettings will return the MailSlurper service tier address
	 * and port.
	 */
	getServiceSettings: function () {
		return new Promise(function (resolve, reject) {
			$.ajax({
				method: "GET",
				url: "/servicesettings"
			}).then(
				function (result) {
					return resolve(result);
				},
				function (err) {
					return reject(err);
				}
			);
		})
	},

	/**
	 * getServiceURL returns a fully formatted service URL
	 */
	getServiceURL: function () {
		var serviceSettings = window.SettingsService.retrieveServiceSettings();
		var serviceURL = serviceSettings.url;

		serviceURL = serviceURL.replace('0.0.0.0', window.location.hostname);

		return serviceURL;
	},

	/**
	 * getWWURL return the fully formatted app URL
	 */
	getAppURL: function () {
		return $('meta[name=app-url]').attr('content')
	},

	/**
	 * retrieveSavedSearches reads saved searches from local storage
	 */
	retrieveSavedSearches: function () {
		if (localStorage["savedSearches"]) {
			return JSON.parse(localStorage["savedSearches"]);
		} else {
			return [];
		}
	},

	/**
	 * retrieveServiceSettings reads the MailSlurper service settings
	 * from the user's local storage.
	 */
	retrieveServiceSettings: function () {
		return JSON.parse(localStorage["serviceSettings"]);
	},

	/**
	 * retrieveSettings reads user settings from local storage.
	 */
	retrieveSettings: function () {
		if (localStorage["settings"]) {
			return JSON.parse(localStorage["settings"]);
		} else {
			return {
				dateFormat: "YYYY-MM-DD hh:mm A",
				autoRefresh: 0,
				theme: "default"
			};
		}
	},

	/**
	 * serviceSettingsExistInLocalStore returns true/false if the
	 * MailSlurper service settings are in the user's local storage.
	 */
	serviceSettingsExistInLocalStore: function () {
		return (localStorage["serviceSettings"] !== undefined);
	},

	/**
	 * storeSavedSearches stores user saved searches in local storage
	 */
	storeSavedSearches: function (savedSearches) {
		localStorage["savedSearches"] = JSON.stringify(savedSearches);
	},

	/**
	 * storeServiceSettings writes the MailSlurper service settings to
	 * the user's local storage.
	 */
	storeServiceSettings: function (serviceSettings) {
		localStorage["serviceSettings"] = JSON.stringify(serviceSettings);
	},

	/**
	 * storeSettings writes user settings to local storage and updates
	 * the theme in the MailSlurper config file
	 */
	storeSettings: function (settings) {
		return new Promise(function (resolve, reject) {
			localStorage["settings"] = JSON.stringify(settings);

			$.ajax({
				url: "/theme",
				method: "POST",
				dataType: "text",
				contentType: "application/json; charset=utf-8",
				data: JSON.stringify({
					theme: settings.theme
				})
			}).then(
				function () {
					return resolve();
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	}
};
