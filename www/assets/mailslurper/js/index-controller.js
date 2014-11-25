"use strict";

app.controller("IndexCtrl",
	[
		"$scope",
		"$http",

		function($scope, $http) {
			$scope.message = "Hello!";
			$scope.mailItems = [];

			$http.get("http://localhost:8085/v1/mails/page/1").
				success(function(data) {
					$scope.mailItems = data.mailItems;
				});
		}
	]
);
