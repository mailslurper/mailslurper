require(
	[
		"jquery",
		"services/SettingsService",
		"services/MailService",
		"services/AlertService",
		"widgets/SavedSearchesWidget",
		"bootstrap-dialog",
		"moment",

		"hbs!templates/mailList",
		"hbs!templates/mailDetails",
		"hbs!templates/searchMailModal",

		"lightbox",
		"bootstrap-daterangepicker"
	],
	function(
		$,
		SettingsService,
		MailService,
		AlertService,
		SavedSearchesWidget,
		Dialog,
		moment,
		mailListTemplate,
		mailDetailsTemplate,
		searchMailModalTemplate
	) {
		"use strict";

		/**
		 * Calculates the height of the window minus nav bars and table headers.
		 */
		var calculateWindowHeight = function() {
			return $(window).outerHeight(true) -
				$("#mailItemsHeader").outerHeight(true) -
				$(".navbar").outerHeight(true) -
				$("#mailSearchNav").outerHeight(true);
		};

		/**
		 * Retrieves an attachment and displays it to the user. This expects the context to
		 * have "attachmentID" and "mailID".
		 */
		var displayAttachment = function(context) {
			context.message = "Retrieving...";

			AlertService.block(context)
				.then(MailService.getAttachment)
				.then(showAttachmentInLightbox)
				.then(AlertService.unblock)
				.catch(AlertService.error);

			return Promise.resolve(context);
		};

		/**
		 * Initialize the list of mail items. This will attach click events and
		 * handle resizing of the window so our scrollable content windows adjust
		 * correctly.
		 */
		var initializeMailItems = function(context) {
			$(".mailSubject").on("click", function() {
				var id = $(this).attr("data-id");
				context.mailID = id;

				viewMailDetails(context);
			});

			$("#btnRefresh").on("click", function() {
				refreshMailList(context);
			});

			$("#btnSearch").on("click", function() {
				renderSearchMailModal(context);
			});

			resizeMailItems();
			resizeMailDetails();

			$(window).on("resize", function() {
				resizeMailItems();
				resizeMailDetails();
			});

			return Promise.resolve(context);
		};

		/**
		 * Performs a search for mail items, then re-renders the mail items
		 * window.
		 */
		var performSearch = function(context) {
			context.message = "Searching...";

			AlertService.block(context)
				.then(MailService.getMails)
				.then(renderMailItems)
				.then(initializeMailItems)
				.then(AlertService.unblock)
				.catch(AlertService.error);

			return Promise.resolve(context);
		};

		/**
		 * Refreshes the mail list view. Basically just
		 * performs a search again.
		 */
		var refreshMailList = function(context) {
			return performSearch(context);
		};

		/**
		 * Render the date range picker widget
		 */
		var renderDateRangePicker = function(context, dialogRef) {
			$("#dateRange").daterangepicker({
				ranges: {
					"Today": [moment(), moment()],
					"Yesterday": [moment().subtract(1, "days"), moment().subtract(1, "days")],
					"Last 7 Days": [moment().subtract(6, "days"), moment()],
					"Last 30 Days": [moment().subtract(29, "days"), moment()],
					"This Month": [moment().startOf("month"), moment().endOf("month")],
					"Last Month": [moment().subtract(1, "month").startOf("month"), moment().subtract(1, "month").endOf("month")]
				},
				opens: "right",
				drops: "down",
				minDate: moment().subtract(1, "month").startOf("month"),
				maxDate: moment().endOf("month"),
				startDate: context.searchStart,
				endDate: context.searchEnd
			}, function(start, end) {
				renderDateRangeSpan(context, start, end);
			});
		};

		var renderDateRangeSpan = function(context, start, end) {
			context.searchStart = start;
			context.searchEnd = end;
			$("#dateRange span").html(start.format("MMMM D, YYYY") + " - " + end.format("MMMM D, YYYY"));
		};

		/**
		 * Renders the detail view for a specific mailitem.
		 */
		var renderMailDetails = function(context) {
			var html = mailDetailsTemplate({mail: context.mail.mailItem});
			$("#mailDetails").html(html);

			return Promise.resolve(context);
		};

		/**
		 * Renders the list of mail items.
		 */
		var renderMailItems = function(context) {
			var html = mailListTemplate({mails: context.mails});
			$("#mailList").html(html);

			return Promise.resolve(context);
		};

		/**
		 * Renders and handles events for the search modal dialog box.
		 */
		var renderSearchMailModal = function(context) {
			var dialogRef = Dialog.show({
				title: "Search Mail",
				message: searchMailModalTemplate(),
				closable: true,
				nl2br: false,
				data: {
					start: null,
					end: null,
				},
				buttons: [
					{
						id: "btnSave",
						label: "Save",
						cssClass: "btn-default",
						action: function() {
							var searchCriteria = {
								searchMessage: $("#txtMessage").val()
							};

							SavedSearchesWidget.showSaveSearchModal(function(saveSearchName) {
								SettingsService.addSavedSearch(saveSearchName, searchCriteria);
							});
						}
					},
					{
						id: "btnClearSearch",
						label: "Clear",
						cssClass: "btn-default",
						action: function() {
							context.searchStart = moment().startOf("month");
							context.searchEnd = moment().endOf("month");
							renderDateRangePicker(context);
							renderDateRangeSpan(context, context.searchStart, context.searchEnd);

							$("#txtMessage").val("");
						}
					},
					{
						id: "btnCancelSearch",
						label: "Cancel",
						cssClass: "btn-default",
						action: function(dialogRef) {
							dialogRef.close();
						}
					},
					{
						id: "btnExecuteSearch",
						label: "Search",
						cssClass: "btn-primary",
						hotkey: 13,
						action: function(dialogRef) {
							context.searchMessage = $("#txtMessage").val();
							dialogRef.close();
							performSearch(context);
						}
					}
				],
				onshown: function(dialogRef) {
					renderDateRangePicker(context);
					renderDateRangeSpan(context, context.searchStart, context.searchEnd);
					$("#btnOpenSavedSearches").on("click", function() { showSavedSearchesModal(context); });
					$("#txtMessage").val(context.searchMessage).focus();
				}
			});

			return Promise.resolve(context);
		};

		/**
		 * Resizes the mail detail window
		 */
		var resizeMailDetails = function() {
			$("#mailDetailsColumn").innerHeight(calculateWindowHeight());
		};

		/**
		 * Resizes the mail items list window.
		 */
		var resizeMailItems = function() {
			$("#mailItemsColumn").innerHeight(calculateWindowHeight());
		};

		/**
		 * Displays the saved searches modal
		 */
		var showSavedSearchesModal = function(context) {
			SavedSearchesWidget.showPicker(function(savedSearch) {
				$("#txtMessage").val(savedSearch.searchMessage);
			});
		};

		/**
		 * Loads the details for a selected mail item, then renders them.
		 */
		var viewMailDetails = function(context) {
			context.message = "Getting details...";

			AlertService.block(context)
				.then(MailService.getMailByID)
				.then(renderMailDetails)
				.then(AlertService.unblock)
				.catch(AlertService.error);
		};

		/****************************************************************************
		 * Constructor
		 ***************************************************************************/
		var context = {
			mails: [],
			message: "Loading",
			page: 1,
			searchMessage: "",
			searchStart: moment().startOf("month"),
			searchEnd: moment().endOf("month")
		};

		AlertService.block(context)
			.then(SettingsService.getServiceURL)
			.then(MailService.getMails)
			.then(renderMailItems)
			.then(initializeMailItems)
			.then(AlertService.unblock)
			.catch(AlertService.error);
	}
);
