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

                alertObj.unblock();

                if (type === "success") {
                    messageHtml += " fa-check-circle";
                } else if (type === "error") {
                    type = "danger";
                    messageHtml += " fa-exclamation-circle";
                } else if (type === "information") {
                    type = "info";
                    messageHtml += " fa-info-circle";
                }

                alertObj.logMessage(message, type);

                messageHtml += "\"></i> " + message;
                $.bootstrapGrowl(messageHtml, { type: type });
            },

            block: function(message) {
                $.blockUI({ message: "<i class=\"fa fa-spinner\"></i> " + message });
            },

            error: function(message) {
                alertObj.alert(message, "error");
            },

            information: function(message) {
                alertObj.alert(message, "information");
            },

            success: function(message) {
                alertObj.alert(message, "success");
            },

            unblock: function(context) {
                $.unblockUI();
            },

            logMessage: function(message, type) {
                if (type === undefined) type = "danger";
                if (type === "error") type = "danger";

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