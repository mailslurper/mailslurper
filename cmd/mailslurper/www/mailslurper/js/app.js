// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

$(document).ready(function () {
	var appURL = window.SettingsService.getAppURL();
	
	if (!window.SettingsService.serviceSettingsExistInLocalStore()) {
		window.SettingsService.getServiceSettings()
			.then(function (serviceSettings) {
				window.SettingsService.storeServiceSettings(serviceSettings);
				return serviceSettings;
			})
			.then(function (serviceSettings) {
				if (serviceSettings.authenticationScheme !== "") {
					if (!window.AuthService.tokenExistsInStorage()) {
						window.location = appURL + "/login";
					}
				}
			})
			.catch(function (err) {
				window.AlertService.error("There was an error getting service settings. See the console for more information.");
				window.AlertService.logMessage(err, "error");
			});
	} else {
		var serviceURL = window.SettingsService.getServiceURL();

		$("#logOutLink").on("click", function () {
			window.AuthService.logout(serviceURL)
				.then(function () {
					window.location = appURL + "/logout";
				})
				.catch(function (err) {
					alert("There was an error logging out: " + err)
				});
		});


	}
});

