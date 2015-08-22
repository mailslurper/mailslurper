require(
	[
		"jquery",
		"services/SettingsService",
		"services/MailService",
		"services/SeedService",
		"services/AlertService",
		"bootstrap-dialog",

		"hbs!templates/adminPrune"
	],
	function(
		$,
		SettingsService,
		MailService,
		SeedService,
		AlertService,
		Dialog,
		adminPruneTemplate
	) {
		"use strict";

		var initialize = function(context) {
			$("#btnRemove").on("click", function() { onBtnRemoveClick(context); });
			return Promise.resolve(context);
		};

		var onBtnRemoveClick = function(context) {
			context.pruneCode = $("#pruneRange option:selected").val();

			AlertService.block(context)
				.then(SeedService.validatePruneCode)
				.then(MailService.deleteMailItems)
				.then(MailService.getMailCount)
				.then(renderPruneTemplate)
				.then(initialize)
				.then(AlertService.unblock)
				.catch(AlertService.error);
		};

		var renderPruneTemplate = function(context) {
			var html = adminPruneTemplate({
				totalEmailCount: context.mailCount,
				pruneOptions: context.pruneOptions
			});

			$("#adminPrune").html(html);

			return Promise.resolve(context);
		};

		/****************************************************************************
		 * Constructor
		 ***************************************************************************/
		var context = {
			message: "Loading"
		};

		AlertService.block(context)
			.then(SettingsService.getServiceURL)
			.then(SeedService.getPruneOptions)
			.then(MailService.getMailCount)
			.then(renderPruneTemplate)
			.then(initialize)
			.then(AlertService.unblock)
			.catch(AlertService.error);
	}
);
