"use strict";

var app = angular.module("app", ["ngRoute", "ui.layout", "ui.bootstrap"])

.config(["$routeProvider", function($routeProvider) {
	$routeProvider
		.when("/", {
			controller: "IndexCtrl",
			templateUrl: "/assets/mailslurper/templates/mail-list.html"
		})

		.otherwise({
			redirectTo: "/"
		});
}]);
