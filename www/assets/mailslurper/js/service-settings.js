"use strict";

var ServiceSettings = {
	address: "",
	port: 0,

	buildUrl: function(path) {
		return "http://" + ServiceSettings.address + ":" + ServiceSettings.port + "/v1" + path;
	},

	getServiceSettings: function() {
		$.ajax({ url: "/servicesettings" }).done(function(data) {
			ServiceSettings.address = data.serviceAddress;
			ServiceSettings.port = data.servicePort;
		});
	}
};
