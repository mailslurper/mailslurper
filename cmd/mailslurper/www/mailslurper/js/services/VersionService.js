// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

window.VersionService = {
	getVersionFromGithub: function (serviceURL) {
		return new Promise(function (resolve, reject) {
			$.ajax({
				url: "/masterversion",
				method: "GET",
				cache: false
			}).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	},

	getServerVersion: function () {
		return new Promise(function (resolve, reject) {
			$.ajax({
				url: "/version",
				method: "GET",
				cache: false
			}).then(
				function (result) {
					return resolve(result);
				},
				function (xhr, errorType, err) {
					return reject(err);
				}
			);
		});
	},

	isVersionOlder: function (version1, version2) {
		var split1 = version1.split(".");
		var split2 = version2.split(".");

		var major1 = window.parseInt(split1[0]);
		var minor1 = window.parseInt(split1[1]);
		var build1 = window.parseInt(split1[2]);

		var major2 = window.parseInt(split2[0]);
		var minor2 = window.parseInt(split2[1]);
		var build2 = window.parseInt(split2[2]);

		if (major1 < major2) return true;
		if (major1 > major2) return false;

		if (minor1 < minor2) return true;
		if (minor1 > minor2) return false;

		if (build1 < build2) return true;
		return false;
	}
};

