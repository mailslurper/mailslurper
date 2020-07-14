// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

window.AuthService = {
	login: function (serviceURL, userName, password) {
		return new Promise(function (resolve, reject) {
			$.ajax({
				url: serviceURL + "/login",
				method: "POST",
				data: {
					userName: userName,
					password: password
				}
			}).then(
				function (token) {
					return resolve(token);
				},

				function (response) {
					return reject(response.responseText);
				}
			);
		});
	},

	logout: function (serviceURL) {
		return new Promise(function (resolve, reject) {
			$.ajax({
				url: serviceURL + "/logout",
				method: "DELETE"
			}).then(
				function () {
					return resolve();
				},

				function (response) {
					return reject(response.responseText);
				}
			);
		});
	},

	storeToken: function (token) {
		localStorage["jwt"] = token;
	},

	getToken: function () {
		return localStorage["jwt"];
	},

	tokenExistsInStorage: function () {
		return (localStorage["jwt"] !== undefined) ? true : false;
	},

	decorateRequestWithAuthorization: function (requestParameters) {
		if (window.AuthService.tokenExistsInStorage()) {
			requestParameters.beforeSend = function (xhr) {
				xhr.setRequestHeader("Authorization", "Bearer " + window.AuthService.getToken());
			}
		}

		return requestParameters;
	},

	isUnauthorized: function (err) {
		if (err === "Forbidden") {
			return true;
		}

		return false;
	},

	gotoLogin: function () {
		var appURL = window.SettingsService.getAppURL();
		window.location = appURL + "/login";
	}
};
