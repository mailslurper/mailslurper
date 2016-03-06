// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

require(
	[
		"jquery",
		"services/SettingsService",
		"services/MailService",
		"services/SeedService",
		"services/AlertService",
		"bootstrap-dialog",

		"hbs!templates/adminPrune",
		"hbs!templates/adminSettings"
	],
	function(
		$,
		settingsService,
		MailService,
		SeedService,
		alertService,
		Dialog,
		adminPruneTemplate,
		adminSettings
	) {
		"use strict";

		var getSettingsFromForm = function() {
			var settings = {
				dateFormat: $("#dateFormat option:selected").val(),
				autoRefresh: window.parseInt($("#autoRefresh option:selected").val(), 10)
			};

			return settings;
		};

		var initialize = function() {
			$("#btnRemove").on("click", function() { onBtnRemoveClick(); });
			$("#btnSaveSettings").on("click", function() { onBtnSaveSettings(); });
		};

		var onBtnRemoveClick = function() {
			Dialog.confirm({
				message: "Are you sure you wish to prune old emails?",
				title: "WARNING",
				type: Dialog.TYPE_WARNING,
				callback: function(result) {
					if (result) {
						var pruneCode = $("#pruneRange option:selected").val();

						alertService.block("Pruning...");

						if (!SeedService.validatePruneCode(pruneOptions, pruneCode)) {
							alertService.error("There was an error with the selected prune option.");
							return;
						}

						MailService.deleteMailItems(serviceURL, pruneCode).then(
							function() {
								MailService.getMailCount(serviceURL).then(function(response) {
									renderPruneTemplate(pruneOptions, response.mailCount);
									initialize();

									alertService.unblock();
									showPruneSuccessMessage();
								});
							},

							function() {
								alertService.error("There was an error deleting mail items.");
							}
						);
					}
				}
			});
		};

		var onBtnSaveSettings = function() {
			var settings = getSettingsFromForm();
			settingsService.storeSettings(settings);

			alertService.success("Settings saved!");
		};

		var renderPruneTemplate = function(pruneOptions, mailCount) {
			var html = adminPruneTemplate({
				totalEmailCount: mailCount,
				pruneOptions: pruneOptions
			});

			$("#adminPrune").html(html);
		};

		var renderSettingsTemplate = function(settings, dateFormatOptions) {
			var html = adminSettings({
				dateFormat: settings.dateFormat,
				dateFormatOptions: dateFormatOptions,
				autoRefresh: settings.autoRefresh
			});

			$("#adminSettings").html(html);
		};

		var showPruneSuccessMessage = function() {
			alertService.success("Emails pruned successfully!");
		};

		/****************************************************************************
		 * Constructor
		 ***************************************************************************/
		var serviceURL = settingsService.getServiceURL();
		var pruneOptions = [];

		SeedService.getPruneOptions(serviceURL).then(
			function(response) {
				pruneOptions = response;

				MailService.getMailCount(serviceURL).then(
					function(response) {
						var settings = settingsService.retrieveSettings();
						var dateFormatOptions = SeedService.getDateFormatOptions();

						renderPruneTemplate(pruneOptions, response.mailCount);
						renderSettingsTemplate(settings, dateFormatOptions);
						initialize();

						alertService.unblock();
					}
				);
			}
		);
	}
);
