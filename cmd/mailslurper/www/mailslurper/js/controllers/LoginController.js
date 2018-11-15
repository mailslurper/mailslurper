// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

(function () {
	"use strict";

	function getUserName() {
		return $("#userName").val();
	}

	function getPassword() {
		return $("#password").val();
	}

	function validateForm() {
		var result = {
			success: true,
			message: "",
			field: ""
		};

		try {
			if (getUserName() === "") {
				result.field = "userName";
				throw new Error("Please enter a user name");
			}

			if (getPassword() === "") {
				result.field = "password";
				throw new Error("Please enter a password");
			}
		} catch (err) {
			result.success = false;
			result.message = err.message;
		}

		return result;
	}

	function submitLogin(e) {
		var validationResult = validateForm();

		if (!validationResult.success) {
			$("#" + validationResult.field).focus();
			alert(validationResult.message);
			return false;
		}

		var serviceURL = window.SettingsService.getServiceURL();

		window.AuthService.login(serviceURL, getUserName(), getPassword())
			.then(function (token) {
				window.AuthService.storeToken(token);
				$("#frmLogin").submit();
			})
			.catch(function (err) {
				var appURL = window.SettingsService.getAppURL();
				window.location = appURL + "/login?message=Invalid user name or password";
			});
	}

	/****************************************************************************
	 * Constructor
	 ***************************************************************************/
	var serviceSettings = {};

	$("#btnSubmit").on("click", submitLogin);
	$("#userName").on("keypress", function (e) {
		if (e.which === 13) {
			submitLogin();
		}
	});
	$("#password").on("keypress", function (e) {
		if (e.which === 13) {
			submitLogin();
		}
	});

	window.SettingsService.getServiceSettings()
		.then(function (settings) {
			serviceSettings = settings;
			window.SettingsService.storeServiceSettings(serviceSettings);
		})
		.catch(function (err) {
			console.log(err);
			alert("Error getting service settings! See console.log");
		});
}());