define(
	[
		"jquery",
		"services/AlertService"
	],
	function($, AlertService) {
		"use strict";

		var service = {
			/**
			 * addSavedSearch will append a new saved search
			 */
			addSavedSearch: function(name, searchCriteria) {
				searchCriteria.name = name;
				var currentSavedSearches = service.retrieveSavedSearches();

				currentSavedSearches.push(searchCriteria);
				service.storeSavedSearches(currentSavedSearches);
			},

			deleteSavedSearch: function(index) {
				var currentSavedSearches = service.retrieveSavedSearches();

				if (index > -1 && index < currentSavedSearches.length) {
					currentSavedSearches.splice(index, 1);
					service.storeSavedSearches(currentSavedSearches);
				}
			},

			/**
			 * getSavedSearchByIndex returns a saved search based on its location
			 * in the array.
			 */
			getSavedSearchByIndex: function(index) {
				var searches = service.retrieveSavedSearches();

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
			getServiceSettings: function() {
				return $.ajax({
					method: "GET",
					url: "/servicesettings"
				});
			},

			/**
			 * getServiceURL returns a fully formatted service URL as a key named
			 * "serviceURL" in the context object.
			 */
			getServiceURL: function(context) {
				var serviceSettings = service.retrieveServiceSettings();
				return "//" + serviceSettings.serviceAddress + ":" + serviceSettings.servicePort;
			},

			/**
			 * getServiceURLNow returns a fully formatted service URL as a key named
			 * "serviceURL" directly instead of via a promise.
			 */
			getServiceURLNow: function() {
				var serviceSettings = service.retrieveServiceSettings();

				var serviceURL = "http://" + serviceSettings.serviceAddress + ":" + serviceSettings.servicePort + "/" + serviceSettings.version;

				return serviceURL;
			},

			/**
			 * retrieveSavedSearches reads saved searches from local storage
			 */
			retrieveSavedSearches: function() {
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
			retrieveServiceSettings: function() {
				return JSON.parse(localStorage["serviceSettings"]);
			},

			/**
			 * retrieveSettings reads user settings from local storage.
			 */
			retrieveSettings: function() {
				if (localStorage["settings"]) {
					return JSON.parse(localStorage["settings"]);
				} else {
					return {
						dateFormat: "YYYY-MM-DD hh:mm A",
						autoRefresh: 0
					};
				}
			},

			/**
			 * serviceSettingsExistInLocalStore returns true/false if the
			 * MailSlurper service settings are in the user's local storage.
			 */
			serviceSettingsExistInLocalStore: function() {
				return (localStorage["serviceSettings"] !== undefined);
			},

			/**
			 * storeSavedSearches stores user saved searches in local storage
			 */
			storeSavedSearches: function(savedSearches) {
				localStorage["savedSearches"] = JSON.stringify(savedSearches);
			},

			/**
			 * storeServiceSettings writes the MailSlurper service settings to
			 * the user's local storage.
			 */
			storeServiceSettings: function(serviceSettings) {
				localStorage["serviceSettings"] = JSON.stringify(serviceSettings);
			},

			/**
			 * storeSettings writes user settings to local storage
			 */
			storeSettings: function(settings) {
				localStorage["settings"] = JSON.stringify(settings);
			}
		};

		return service;
	}
);
