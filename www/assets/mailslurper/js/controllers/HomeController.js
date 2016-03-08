// Copyright 2013-2016 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
		settingsService,
		mailService,
		alertService,
		SavedSearchesWidget,
		Dialog,
		moment,
		mailListTemplate,
		mailDetailsTemplate,
		searchMailModalTemplate
	) {
		"use strict";

		/**
		 * Creates the markup for the filters popover
		 */
		var buildFiltersPopoverText = function() {
			var html = "<strong>Current Page:</strong> " + page + "<br />";
			html += "<strong>Message Filter:</strong> " + searchCriteria.searchMessage + "<br />";
			html += "<strong>Date Range:</strong> " + moment(searchCriteria.searchStart).format("MMMM D, YYYY") + " - ";
			html += moment(searchCriteria.searchEnd).format("MMMM D, YYYY") + "<br />";
			html += "<strong>From:</strong> " + searchCriteria.searchFrom + "<br />";
			html += "<strong>To:</strong> " + searchCriteria.searchTo + "<br />";

			return html;
		};

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
		var displayAttachment = function() {
			alertService.block("Retrieving...");

			mailService.getAttachment(serviceURL, mailID, attachmentID).then(
				function(response) {
					alertService.unblock();
				},

				function() {
					alertService.error("There was a problem retrieving your attachment.");
				}
			);
		};

		/**
		 * Highlights a mail row.
		 */
		var highlightMailRow = function(rowID) {
			$("#" + rowID).addClass("mail-list-row-highlight");
		};

		/**
		 * Initialize the list of mail items. This will attach click events and
		 * handle resizing of the window so our scrollable content windows adjust
		 * correctly.
		 */
		var initializeMailItems = function() {
			$(".mailSubject").on("click", function() {
				var id = $(this).attr("data-id");

				removeAllMailRowHighlights();
				highlightMailRow(id);

				mailID = id;
				viewMailDetails();
			});

			$("#btnRefresh").on("click", function() {
				refreshMailList();
			});

			$("#btnSearch").on("click", function() {
				renderSearchMailModal();
			});

			$("#firstPage").on("click", function() {
				page = 1;
				performSearch();
			});

			$("#previousPage").on("click", function() {
				page = previousPage;
				performSearch();
			});

			$("#nextPage").on("click", function() {
				page = nextPage;
				performSearch();
			});

			$("#lastPage").on("click", function() {
				page = totalPages;
				performSearch();
			});

			resizeMailItems();
			resizeMailDetails();

			$(window).on("resize", function() {
				resizeMailItems();
				resizeMailDetails();
			});

			$("#showSearchFilters").popover({
				html: true,
				placement: "left",
				trigger: "click, focus"
			});
		};

		/**
		 * Performs a search for mail items, then re-renders the mail items
		 * window.
		 */
		var performSearch = function() {
			alertService.block("Searching...");

			mailService.getMails(serviceURL, page, searchCriteria).then(
				function(response, status, xhr) {
					mails = response.mailItems;
					totalPages = response.totalPages;
					totalMailCount = response.totalRecordCount;

					renderMailItems();
					initializeMailItems();
					alertService.unblock();

					setRefreshTimeLeft();
				},

				function() {
					alertService.error("There was a problem performing your search");
				}
			);
		};

		/**
		 * Refreshes the mail list view. Basically just
		 * performs a search again.
		 */
		var refreshMailList = function() {
			return performSearch();
		};

		/**
		 * Removes highlights from all mail rows
		 */
		var removeAllMailRowHighlights = function() {
			$(".mailRow").removeClass("mail-list-row-highlight");
		};

		/**
		 * Render the date range picker widget
		 */
		var renderDateRangePicker = function(dialogRef) {
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
				startDate: searchCriteria.searchStart,
				endDate: searchCriteria.searchEnd
			}, function(start, end) {
				renderDateRangeSpan(start, end);
			});
		};

		var renderDateRangeSpan = function(start, end) {
			searchCriteria.searchStart = start;
			searchCriteria.searchEnd = end;
			$("#dateRange span").html(start.format("MMMM D, YYYY") + " - " + end.format("MMMM D, YYYY"));
		};

		/**
		 * Renders the detail view for a specific mailitem.
		 */
		var renderMailDetails = function(mail) {
			var html = mailDetailsTemplate({mail: mail.mailItem});
			$("#mailDetails").html(html);
		};

		/**
		 * Renders the list of mail items.
		 */
		var renderMailItems = function() {
			nextPage = (page < totalPages) ? page + 1 : totalPages;
			previousPage = (page > 1) ? page - 1 : 1;

			var html = mailListTemplate({
				mails: mails,
				totalPages: totalPages,
				hasNavigation: (totalPages > 1) ? true : false,
				hasFirstButton: (page > 1) ? true : false,
				hasPreviousButton: (page > 1) ? true : false,
				hasNextButton: (page < totalPages) ? true : false,
				hasLastButton: (page < totalPages) ? true : false,
				previousPage: previousPage,
				nextPage: nextPage,
				filtersPopover: buildFiltersPopoverText()
			});

			$("#mailList").html(html);
		};

		/**
		 * Renders and handles events for the search modal dialog box.
		 */
		var renderSearchMailModal = function() {
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
								searchMessage: $("#txtMessage").val(),
								searchFrom: $("#txtFrom").val(),
								searchTo: $("#txtTo").val()
							};

							SavedSearchesWidget.showSaveSearchModal(function(saveSearchName) {
								settingsService.addSavedSearch(saveSearchName, searchCriteria);
							});
						}
					},
					{
						id: "btnClearSearch",
						label: "Clear",
						cssClass: "btn-default",
						action: function() {
							searchCriteria.searchStart = moment().startOf("month");
							searchCriteria.searchEnd = moment().endOf("month");
							renderDateRangePicker();
							renderDateRangeSpan(searchCriteria.searchStart, searchCriteria.searchEnd);

							$("#txtMessage").val("");
							$("#txtFrom").val("");
							$("#txtTo").val("");
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
							searchCriteria.searchMessage = $("#txtMessage").val();
							searchCriteria.searchFrom = $("#txtFrom").val();
							searchCriteria.searchTo = $("#txtTo").val();

							dialogRef.close();
							performSearch();
						}
					}
				],
				onshown: function(dialogRef) {
					renderDateRangePicker();
					renderDateRangeSpan(searchCriteria.searchStart, searchCriteria.searchEnd);
					$("#btnOpenSavedSearches").on("click", function() { showSavedSearchesModal(); });

					$("#txtFrom").val(searchCriteria.searchFrom);
					$("#txtTo").val(searchCriteria.searchTo);
					$("#txtMessage").val(searchCriteria.searchMessage).focus();
				}
			});
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

		var setRefreshTimeLeft = function() {
			var settings = settingsService.retrieveSettings();
			refreshTime = moment.duration(settings.autoRefresh, "minutes");
		};

		/*
		 * Sets up the auto refresh timer
		 */
		var setupAutoRefresh = function() {
			var settings = settingsService.retrieveSettings();

			if (settings.autoRefresh > 0) {
				alertService.logMessage("Auto refresh set to " + settings.autoRefresh + " minute(s)", "info");

				var timeLeft = settings.autoRefresh * 60 * 1000;
				setRefreshTimeLeft();
				updateAutoRefreshCountdown();

				window.setInterval(updateAutoRefreshCountdown, 10000);
				window.setInterval(performSearch, timeLeft);
			}
		};

		/**
		 * Displays the saved searches modal
		 */
		var showSavedSearchesModal = function() {
			SavedSearchesWidget.showPicker(function(savedSearch) {
				$("#txtMessage").val(savedSearch.searchMessage);
				$("#txtFrom").val(savedSearch.searchFrom);
				$("#txtTo").val(savedSearch.searchTo);
			});
		};

		/*
		 * Updates the auto-refresh countdown timer
		 */
		var updateAutoRefreshCountdown = function() {
			$("#refreshCountdownText").html("(" + refreshTime.humanize() + ")");
			refreshTime = moment.duration(refreshTime.asSeconds() - 10, "seconds");
		};

		/**
		 * Loads the details for a selected mail item, then renders them.
		 */
		var viewMailDetails = function() {
			alertService.block("Getting details...");

			mailService.getMailByID(serviceURL, mailID).then(
				function(response) {
					renderMailDetails(response);
					alertService.unblock();
				},

				function() {
					alertService.error("There was a problem getting this mail's details");
				}
			);
		};

		/****************************************************************************
		 * Constructor
		 ***************************************************************************/
		var mails = [];
		var mailID = 0;
		var previousPage = 0;
		var nextPage = 0;
		var refreshTime = 0;
		var totalPages = 0;
		var totalMailCount = 0;
		var page = 1;
		var searchCriteria = {
			searchMessage: "",
			searchStart: moment().startOf("month"),
			searchEnd: moment().endOf("month"),
			searchFrom: "",
			searchTo: ""
		};

		var serviceURL = settingsService.getServiceURL();

		alertService.block("Loading");

		mailService.getMails(serviceURL, page, searchCriteria).then(
			function(response, status, xhr) {
				mails = response.mailItems;
				totalPages = response.totalPages;
				totalMailCount = response.totalRecordCount;

				renderMailItems();
				initializeMailItems();
				alertService.unblock();

				setupAutoRefresh();
			},

			function() {
				alertService.error("There was an error getting mail items!");
			}
		);
	}
);
