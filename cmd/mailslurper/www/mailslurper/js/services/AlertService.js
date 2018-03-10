// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

"use strict";

window.AlertService = {
    alert: function (message, type) {
        var messageHtml = "<i class=\"fa";

        window.AlertService.unblock();

        if (type === "success") {
            messageHtml += " fa-check-circle";
        } else if (type === "error") {
            type = "danger";
            messageHtml += " fa-exclamation-circle";
        } else if (type === "information") {
            type = "info";
            messageHtml += " fa-info-circle";
        }

        window.AlertService.logMessage(message, type);

        messageHtml += "\"></i> " + message;
        $.bootstrapGrowl(messageHtml, { type: type });
    },

    block: function (message) {
        $.blockUI({ message: "<i class=\"fa fa-spinner\"></i> " + message });
    },

    error: function (message) {
        window.AlertService.alert(message, "error");
    },

    information: function (message) {
        window.AlertService.alert(message, "information");
    },

    success: function (message) {
        window.AlertService.alert(message, "success");
    },

    unblock: function (context) {
        $.unblockUI();
    },

    logMessage: function (message, type) {
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
