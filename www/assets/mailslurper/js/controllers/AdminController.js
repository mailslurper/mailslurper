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
			Dialog.confirm({
				message: "Are you sure you wish to prune old emails?",
				title: "WARNING",
				type: Dialog.TYPE_WARNING,
				callback: function(result) {
					if (result) {
						context.pruneCode = $("#pruneRange option:selected").val();

						AlertService.block(context)
							.then(SeedService.validatePruneCode)
							.then(MailService.deleteMailItems)
							.then(MailService.getMailCount)
							.then(renderPruneTemplate)
							.then(initialize)
							.then(AlertService.unblock)
							.then(showPruneSuccessMessage)
							.catch(AlertService.error);
					}
				}
			});
		};

		var renderPruneTemplate = function(context) {
			var html = adminPruneTemplate({
				totalEmailCount: context.mailCount,
				pruneOptions: context.pruneOptions
			});

			$("#adminPrune").html(html);

			return Promise.resolve(context);
		};

		var showPruneSuccessMessage = function(context) {
			context.message = "Emails pruned successfully!";
			AlertService.success(context);
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
