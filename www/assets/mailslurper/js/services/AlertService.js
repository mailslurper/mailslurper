define(
	[
		"jquery",

		"blockui",
		"bootstrap-growl"
	],
	function($) {
		"use strict";

		var alertObj = {
			alert: function(message, type) {
				var messageHtml = "<i class=\"fa";
				var messageString = "";

				alertObj.unblock();

				if (message.hasOwnProperty("length")) {
					messageString = message;
				} else if (message.hasOwnProperty("message")) {
					messageString = message.message;
				} else {
					throw ("Please provide a string message or object with a string in a key named 'message'");
				}

				if (type === "success") {
					alertObj.logMessage(messageString, type);
					messageHtml += " fa-check-circle";
				} else if (type === "error") {
					type = "danger";
					alertObj.logMessage(messageString, type);
					messageHtml += " fa-exclamation-circle";
				} else if (type === "information") {
					type = "info";
					alertObj.logMessage(messageString, type);
					messageHtml += " fa-info-circle";
				}

				messageHtml += "\"></i> " + messageString;

				if (window.hasOwnProperty("console")) {
					console.error(messageString);
				}
				
				$.bootstrapGrowl(messageHtml, { type: type });
			},

			block: function(message) {
				if (message.hasOwnProperty("length")) {
					$.blockUI({ message: "<i class=\"fa fa-spinner\"></i> " + message });
				} else if (message.hasOwnProperty("message")) {
					$.blockUI({ message: "<i class=\"fa fa-spinner\"></i> " + message.message });
				} else {
					throw ("Please provide a string message, or an object with a message key");
				}

				return Promise.resolve(message);
			},

			error: function(message) {
				alertObj.alert(message, "error");
				return Promise.resolve(message);
			},

			information: function(message) {
				alertObj.alert(message, "information");
				return Promise.resolve(message);
			},

			success: function(message) {
				alertObj.alert(message, "success");
				return Promise.resolve(message);
			},

			unblock: function(context) {
				$.unblockUI();
				return Promise.resolve(context);
			},

			logMessage: function(message, type) {
				if ("console" in window) {
					if (type === "success") {
						console.log(message);
					} else if (type === "danger") {
						console.error(message);
					} else if (type === "info") {
						console.info(message);
					}
				}
			}
		};

		return alertObj;
	}
);
