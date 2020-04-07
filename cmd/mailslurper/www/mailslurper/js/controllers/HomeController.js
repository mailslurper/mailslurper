// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

(function () {
	"use strict";

	/**
	 * Creates the markup for the filters popover
	 */
	function buildFiltersPopoverText() {
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
	function calculateWindowHeight() {
		return $(window).outerHeight(true) -
			$("#mailItemsHeader").outerHeight(true) -
			$(".navbar").outerHeight(true) -
			$("#mailSearchNav").outerHeight(true);
	};

	/**
	 * Changes sort field and/or direction
	 */
	function changeSortBy(field) {
		if (field === sortCriteria.orderByField) {
			sortCriteria.orderByDirection = (sortCriteria.orderByDirection === "asc") ? "desc" : "asc";
		} else {
			sortCriteria.orderByField = field;
			sortCriteria.orderByDirection = "asc";
		}
	};

	/**
	 * Retrieves an attachment and displays it to the user. This expects the context to
	 * have "attachmentID" and "mailID".
	 */
	function displayAttachment() {
		window.AlertService.block("Retrieving...");

		window.MailService.getAttachment(serviceURL, mailID, attachmentID)
			.then(function (response) {
				window.AlertService.unblock();
			})
			.catch(function () {
				if (window.AuthService.isUnauthorized(err)) {
					window.AuthService.gotoLogin();
				}

				window.AlertService.error("There was a problem retrieving your attachment.");
			});
	};

	/**
	 * Highlights a mail row.
	 */
	function highlightMailRow(rowID) {
		$("#" + rowID).addClass("mail-list-row-highlight");
	};

	/**
	 * Initialize the list of mail items. This will attach click events and
	 * handle resizing of the window so our scrollable content windows adjust
	 * correctly.
	 */
	function initializeMailItems() {
		$(".mailSubject").on("click", function () {
			var id = $(this).attr("data-id");

			removeAllMailRowHighlights();
			highlightMailRow(id);

			mailID = id;
			viewMailDetails();
		});

		$("#btnRefresh").on("click", function () {
			refreshMailList();
		});

		$("#btnSearch").on("click", function () {
			renderSearchMailModal();
		});

		$("#firstPage").on("click", function () {
			page = 1;
			performSearch();
		});

		$("#previousPage").on("click", function () {
			page = previousPage;
			performSearch();
		});

		$("#nextPage").on("click", function () {
			page = nextPage;
			performSearch();
		});

		$("#lastPage").on("click", function () {
			page = totalPages;
			performSearch();
		});

		$("#sortDate").on("click", function () {
			changeSortBy("date");
			performSearch();
		});

		$("#sortSubject").on("click", function () {
			changeSortBy("subject");
			performSearch();
		});

		$("#sortFrom").on("click", function () {
			changeSortBy("from");
			performSearch();
		});

		$("#goToPage").on("change", function () {
			page = window.parseInt($("#goToPage option:selected").val());
			performSearch();
		});

		resizeMailItems();
		resizeMailDetails();

		$(window).on("resize", function () {
			resizeMailItems();
			resizeMailDetails();
		});

		$("#showSearchFilters").popover({
			html: true,
			placement: "left",
			trigger: "click, focus"
		});

		$("#openInTab").on("click", function () {
			var id = $(this).attr("data-id");

			if (id === "") {
				return;
			}

			var url = window.MailService.getMailMessageURL(serviceURL, id);
			window.open(url);
		});
	};

	/*
	 * Loads the mail details template
	 */
	function loadMailDetailsTemplate() {
		return new Promise(function (resolve, reject) {
			window.TemplateService.load("mailDetails")
				.then(function (template) {
					mailDetailsTemplate = Handlebars.compile(template);
					return resolve();
				})
				.catch(function (err) {
					return reject(err);
				});
		});
	}

	/*
	 * Loads the mail list template
	 */
	function loadMailListTemplate() {
		return new Promise(function (resolve, reject) {
			window.TemplateService.load("mailList")
				.then(function (template) {
					mailListTemplate = Handlebars.compile(template);
					return resolve();
				})
				.catch(function (err) {
					return reject(err);
				});
		});
	}

	/*
	 * Loads the search mail modal template
	 */
	function loadSearchMailModalTemplate() {
		return new Promise(function (resolve, reject) {
			window.TemplateService.load("searchMailModal")
				.then(function (template) {
					searchMailModalTemplate = Handlebars.compile(template);
					return resolve();
				})
				.catch(function (err) {
					return reject(err);
				});
		});
	}

	/**
	 * Performs a search for mail items, then re-renders the mail items
	 * window.
	 */
	function performSearch() {
		return new Promise(function (resolve, reject) {
			window.AlertService.block("Searching...");

			window.MailService.getMails(serviceURL, page, searchCriteria, sortCriteria)
				.then(function (response) {
					mails = response.mailItems;
					totalPages = response.totalPages;
					totalMailCount = response.totalRecordCount;

					renderMailItems();
					initializeMailItems();
					window.AlertService.unblock();

					setRefreshTimeLeft();

					return resolve();
				})
				.catch(function (err) {
					if (window.AuthService.isUnauthorized(err)) {
						window.AuthService.gotoLogin();
					}

					window.AlertService.error("There was a problem performing your search");
					return reject(err);
				});
		});
	};

	/**
	 * Refreshes the mail list view. Basically just
	 * performs a search again.
	 */
	function refreshMailList() {
		return performSearch();
	};

	/**
	 * Removes highlights from all mail rows
	 */
	function removeAllMailRowHighlights() {
		$(".mailRow").removeClass("mail-list-row-highlight");
	};

	/**
	 * Render the date range picker widget
	 */
	function renderDateRangePicker(dialogRef) {
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
		}, function (start, end) {
			renderDateRangeSpan(start, end);
		});
	};

	function renderDateRangeSpan(start, end) {
		searchCriteria.searchStart = start;
		searchCriteria.searchEnd = end;
		$("#dateRange span").html(start.format("MMMM D, YYYY") + " - " + end.format("MMMM D, YYYY"));
	};

	/**
	 * Renders the detail view for a specific mailitem.
	 */
	function renderMailDetails(mail) {
		var html = mailDetailsTemplate({ mail: mail });
		$("#mailDetails").html(html);
		$("#openInTab").attr("data-id", mail.id);
	};

	/**
	 * Renders the list of mail items.
	 */
	function renderMailItems() {
		nextPage = (page < totalPages) ? page + 1 : totalPages;
		previousPage = (page > 1) ? page - 1 : 1;

		var chevron = (sortCriteria.orderByDirection === "desc") ? "fa fa-chevron-down" : "fa fa-chevron-up";

		var dateSortIcon = "";
		var subjectSortIcon = "";
		var fromSortIcon = "";

		if (sortCriteria.orderByField === "date") {
			dateSortIcon = " <i class=\"" + chevron + "\"></i>";
		}

		if (sortCriteria.orderByField === "subject") {
			subjectSortIcon = " <i class=\"" + chevron + "\"></i>";
		}

		if (sortCriteria.orderByField === "from") {
			fromSortIcon = " <i class=\"" + chevron + "\"></i>";
		}

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
			filtersPopover: buildFiltersPopoverText(),
			dateSortIcon: dateSortIcon,
			subjectSortIcon: subjectSortIcon,
			fromSortIcon: fromSortIcon,
			direction: sortCriteria.orderByDirection,
			page: page
		});

		$("#mailList").html(html);
	};

	/**
	 * Renders and handles events for the search modal dialog box.
	 */
	function renderSearchMailModal() {
		var dialogRef = BootstrapDialog.show({
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
					action: function () {
						var searchCriteria = {
							searchMessage: $("#txtMessage").val(),
							searchFrom: $("#txtFrom").val(),
							searchTo: $("#txtTo").val()
						};

						window.SavedSearchesWidget.showSaveSearchModal(function (saveSearchName) {
							window.SettingsService.addSavedSearch(saveSearchName, searchCriteria);
						});
					}
				},
				{
					id: "btnClearSearch",
					label: "Clear",
					cssClass: "btn-default",
					action: function () {
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
					action: function (dialogRef) {
						dialogRef.close();
					}
				},
				{
					id: "btnExecuteSearch",
					label: "Search",
					cssClass: "btn-primary",
					hotkey: 13,
					action: function (dialogRef) {
						searchCriteria.searchMessage = $("#txtMessage").val();
						searchCriteria.searchFrom = $("#txtFrom").val();
						searchCriteria.searchTo = $("#txtTo").val();

						dialogRef.close();
						performSearch();
					}
				}
			],
			onshown: function (dialogRef) {
				renderDateRangePicker();
				renderDateRangeSpan(searchCriteria.searchStart, searchCriteria.searchEnd);
				$("#btnOpenSavedSearches").on("click", function () { showSavedSearchesModal(); });

				$("#txtFrom").val(searchCriteria.searchFrom);
				$("#txtTo").val(searchCriteria.searchTo);
				$("#txtMessage").val(searchCriteria.searchMessage).focus();
			}
		});
	};

	/**
	 * Resizes the mail detail window
	 */
	function resizeMailDetails() {
		$("#mailDetailsColumn").innerHeight(calculateWindowHeight());
	};

	/**
	 * Resizes the mail items list window.
	 */
	function resizeMailItems() {
		$("#mailItemsColumn").innerHeight(calculateWindowHeight());
	};

	function setRefreshTimeLeft() {
		var settings = window.SettingsService.retrieveSettings();
		refreshTime = moment.duration(settings.autoRefresh, "minutes");
	};

	/*
	 * Sets up the auto refresh timer
	 */
	function setupAutoRefresh() {
		var settings = window.SettingsService.retrieveSettings();

		if (settings.autoRefresh > 0) {
			window.AlertService.logMessage("Auto refresh set to " + settings.autoRefresh + " minute(s)", "info");

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
	function showSavedSearchesModal() {
		window.SavedSearchesWidget.showPicker(function (savedSearch) {
			if (savedSearch) {
				$("#txtMessage").val(savedSearch.searchMessage);
				$("#txtFrom").val(savedSearch.searchFrom);
				$("#txtTo").val(savedSearch.searchTo);
			}
		});
	};

	/*
	 * Updates the auto-refresh countdown timer
	 */
	function updateAutoRefreshCountdown() {
		$("#refreshCountdownText").html("(" + refreshTime.humanize() + ")");
		refreshTime = moment.duration(refreshTime.asSeconds() - 10, "seconds");
	};

	/**
	 * Loads the details for a selected mail item, then renders them.
	 */
	function viewMailDetails() {
		window.AlertService.block("Getting details...");

		window.MailService.getMailByID(serviceURL, mailID)
			.then(function (response) {
				renderMailDetails(response);
				window.AlertService.unblock();
			})
			.catch(function () {
				if (window.AuthService.isUnauthorized(err)) {
					window.AuthService.gotoLogin();
				}

				window.AlertService.error("There was a problem getting this mail's details");
			});
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
	var sortCriteria = {
		orderByField: "date",
		orderByDirection: "desc"
	};
	var mailDetailsTemplate = null;
	var mailListTemplate = null;
	var searchMailModalTemplate = null;

	var serviceURL = window.SettingsService.getServiceURL();

	window.AlertService.block("Loading");

	Promise.all([
		loadMailDetailsTemplate(),
		loadMailListTemplate(),
		loadSearchMailModalTemplate()
	])
		.then(performSearch)
		.catch(function (err) {
			if (window.AuthService.isUnauthorized(err)) {
				window.AuthService.gotoLogin();
			}

			window.AlertService.error("There was an error getting mail items!");
		});
}());