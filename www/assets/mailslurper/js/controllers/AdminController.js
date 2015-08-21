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
			return Promise.resolve(context);
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
