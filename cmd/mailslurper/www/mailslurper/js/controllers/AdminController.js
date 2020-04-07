// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

(function () {
	"use strict";

	function checkVersion() {
		var thisMailSlurperVersion = "";
		var currentMailSlurperVersion = "";

		window.VersionService.getServerVersion()
			.then(function (data) {
				thisMailSlurperVersion = data.version;

				window.VersionService.getVersionFromGithub(serviceURL).then(function (data) {
					currentMailSlurperVersion = data.version;

					if (window.VersionService.isVersionOlder(thisMailSlurperVersion, currentMailSlurperVersion)) {
						$("#yourVersion").html(thisMailSlurperVersion);
						$("#currentVersion").html(currentMailSlurperVersion);
						$("#versionMessage").removeClass("hidden");
					}
				});
			})
			.catch(function (err) {
				if (window.AuthService.isUnauthorized(err)) {
					window.AuthService.gotoLogin();
				}

				window.AlertService.error(err);
			});
	}

	function getSettingsFromForm() {
		var settings = {
			dateFormat: $("#dateFormat option:selected").val(),
			autoRefresh: window.parseInt($("#autoRefresh option:selected").val(), 10),
			theme: $("#theme option:selected").val()
		};

		return settings;
	}

	function initialize() {
		$("#btnRemove").on("click", function () { onBtnRemoveClick(); });
		$("#btnSaveSettings").on("click", function () { onBtnSaveSettings(); });

		checkVersion();
	}

	function loadAdminPruneTemplate() {
		return new Promise(function (resolve, reject) {
			window.TemplateService.load("adminPrune")
				.then(function (template) {
					adminPruneTemplate = Handlebars.compile(template);
					return resolve();
				})
				.catch(function (err) {
					return reject(err);
				});
		});
	}

	function loadAdminSettings() {
		return new Promise(function (resolve, reject) {
			window.TemplateService.load("adminSettings")
				.then(function (template) {
					adminSettings = Handlebars.compile(template);
					return resolve();
				})
				.catch(function (err) {
					return reject(err);
				});
		});
	}

	function onBtnRemoveClick() {
		BootstrapDialog.confirm({
			message: "Are you sure you wish to prune old emails?",
			title: "WARNING",
			type: BootstrapDialog.TYPE_WARNING,
			callback: function (result) {
				if (result) {
					var pruneCode = $("#pruneRange option:selected").val();

					window.AlertService.block("Pruning...");

					if (!window.SeedService.validatePruneCode(pruneOptions, pruneCode)) {
						window.AlertService.error("There was an error with the selected prune option.");
						return;
					}

					window.MailService.deleteMailItems(serviceURL, pruneCode)
						.then(function (rowsAffected) {
							window.MailService.getMailCount(serviceURL)
								.then(function (response) {
									renderPruneTemplate(pruneOptions, response.mailCount);
									initialize();

									window.AlertService.unblock();
									showPruneSuccessMessage(window.parseInt(rowsAffected));
								})
								.catch(function (err) {
									if (window.AuthService.isUnauthorized(err)) {
										window.AuthService.gotoLogin();
									}

									window.AlertService.error("There was a problem getting mail count");
								});
						})
						.catch(function (err) {
							if (window.AuthService.isUnauthorized(err)) {
								window.AuthService.gotoLogin();
							}

							window.AlertService.error("There was an error deleting mail items.");
						});
				}
			}
		});
	}

	function onBtnSaveSettings() {
		var settings = getSettingsFromForm();

		window.SettingsService.storeSettings(settings)
			.then(function () {
				if (settings.theme != currentTheme) {
					var appURL = window.SettingsService.getAppURL();
					window.location = appURL + "/admin";
					return;
				}

				window.AlertService.success("Settings saved!");
			})
			.catch(function (err) {
				if (window.AuthService.isUnauthorized(err)) {
					window.AuthService.gotoLogin();
				}

				window.AlertService.error(err);
			});
	}

	function renderPruneTemplate(pruneOptions, mailCount) {
		var html = adminPruneTemplate({
			totalEmailCount: mailCount,
			pruneOptions: pruneOptions
		});

		$("#adminPrune").html(html);
	}

	function renderSettingsTemplate(settings, dateFormatOptions) {
		var html = adminSettings({
			dateFormat: settings.dateFormat,
			dateFormatOptions: dateFormatOptions,
			autoRefresh: settings.autoRefresh,
			theme: settings.theme
		});

		$("#adminSettings").html(html);
	}

	function showPruneSuccessMessage(rowsAffected) {
		window.AlertService.success("" + rowsAffected + " email(s) pruned");
	}

	/****************************************************************************
	 * Constructor
	 ***************************************************************************/
	var adminPruneTemplate;
	var adminSettings;

	var serviceURL = window.SettingsService.getServiceURL();
	var pruneOptions = [];
	var currentTheme = "";

	Promise.all([
		loadAdminPruneTemplate(),
		loadAdminSettings()
	])
		.catch(function (err) {
			window.AlertService.error(err);
		});

	window.SeedService.getPruneOptions(serviceURL)
		.then(function (response) {
			pruneOptions = response;

			window.MailService.getMailCount(serviceURL)
				.then(function (response) {
					var settings = window.SettingsService.retrieveSettings();
					var dateFormatOptions = window.SeedService.getDateFormatOptions();

					currentTheme = settings.theme;

					renderPruneTemplate(pruneOptions, response.mailCount);
					renderSettingsTemplate(settings, dateFormatOptions);
					initialize();

					window.AlertService.unblock();
				})
				.catch(function (err) {
					if (window.AuthService.isUnauthorized(err)) {
						window.AuthService.gotoLogin();
					}

					window.AlertService.error(err);
				});
		})
		.catch(function (err) {
			if (window.AuthService.isUnauthorized(err)) {
				window.AuthService.gotoLogin();
			}

			window.AlertService.error(err);
		});
}());
