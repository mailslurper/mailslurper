// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

define(
	[
		"jquery"
	],
	function($) {
		"use strict";

		var service = {
			getVersionFromGithub: function(serviceURL) {
				return $.ajax({
					url: serviceURL + "/version",
					method: "GET",
					cache: false
				});
			},

			getServerVersion: function() {
				return $.ajax({
					url: "/version",
					method: "GET",
					cache: false
				});
			}
		};

		return service;
	}
);
