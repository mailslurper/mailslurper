"use strict";

var MailHomeController = {
	/**
	 * Adds a new mail item to the mails array, which is bound to the interface
	 * and will display the mail item in a table.
	 */
	addMailItemToTable: function(mailItem, view) {
		MailHomeController.clearMailItemSelections();

		var mails = view.get("mails");
		mails.unshift(mailItem);

		view.set("mails", mails);
	},

	/**
	 * Clears any previously selected mail items.
	 */
	clearMailItemSelections: function() {
		$(".mailrow").removeClass("highlight-row");
	},

	/**
	 * Handler for the index page, which consists of the mail item
	 * list and detail views.
	 */
	index: function(ctx) {
		var listView;

		$("#template").load("/assets/mailslurper/templates/mail-list.html", function() {
			listView = new Ractive({
				el: "view",
				template: "#template",
				data: {
					mails: [],
					sortColumn: "dateSent",
					sortDirection: "desc",

					compressTo: function(toAddresses) {
						return toAddresses.join("; ");
					},

					/*
					 * Called when clicking on a header column to sort.
					 * This method will sort the array of data based on a passed
					 * in column and current sort order.
					 *
					 * There is an event attached to this Ractive instance
					 * that will swap the current sort direction.
					 */
					sort: function(array, column) {
						var
							dir = this.get("sortDirection"),
							result1 = (dir === "desc") ? 1 : -1,
							result2 = (dir === "desc") ? -1 : 1;

						array = array.slice();

						return array.sort(function(a, b) {
							return a[column] < b[column] ? result1 : result2;
						});
					},

					getAttachmentIcon: function(attachments) {
						var result = "";

						if (attachments.length > 0) {
							result = "<span class=\"glyphicon glyphicon-paperclip\"></span>";
						}

						return result;
					},

					/*
					 * Returns the correct CSS classes for a column
					 * based on if it is the current sort column and
					 * what the direction is.
					 */
					getSortIcon: function(column) {
						var
							sc = this.get("sortColumn"),
							sd = this.get("sortDirection");

						if (sc !== column) {
							return "";
						} else {
							if (sd === "desc") {
								return "glyphicon glyphicon-arrow-down";
							} else {
								return "glyphicon glyphicon-arrow-up";
							}
						}
					}
				},

				oncomplete: function() {
					MailHomeController.setupLayout();
					MailService.loadMailItems(1, function(mailItems) {
						MailHomeController.updateMailItemsArray(mailItems, listView);
					});
				}
			});

			/*
			 * Assign event listeners to this view
			 */
			listView.on({
				sort: MailHomeController.sortMailItemList,
				viewMailItem: function(e) { MailHomeController.viewMailItem(e, listView); }
			});

			MailHomeController.setupWebsocket(listView);
		});
	},

	/**
	 * Selects a specified mail item. This function takes a
	 * Ractive event item.
	 */
	selectMailItem: function(ractiveEventItem) {
		$(ractiveEventItem.node).addClass("highlight-row");
	},

	/**
	 * Prepare the mail view layout with a north panel for the title
	 * and menu, south panel for footer, and a center and east panel
	 * for the mail item list and detail views.
	 */
	setupLayout: function() {
		$("body").layout({
			north: {
				size: 35,
				resizable: false,
				closable: false
			},
			south: {
				resizable: false,
				closable: false,
				size: 40
			}
		});

		$("#view").layout({
			east: {
				resizable: true,
				closable: true,
				size: "40%"
			}
		});
	},

	/**
	 * Sets up a websocket connection to the web server. Hooks up the
	 * close, message, and error events. The *onmessage* event adds
	 * the passed in mail item to our table.
	 */
	setupWebsocket: function(view) {
		if (window.hasOwnProperty("WebSocket")) {
			MailHomeController.websocketConnection = new WebSocket("ws://" + location.host + "/ws");

			MailHomeController.websocketConnection.onclose = function(e) { logger("Websocket closed"); MailHomeController.websocketConnection = null; }
			MailHomeController.websocketConnection.onmessage = function(e) { MailHomeController.addMailItemToTable($.parseJSON(e.data), view); }
			MailHomeController.websocketConnection.onerror = function(e) { logger("An error occurred on the websocket. Closing."); MailHomeController.websocketConnection.close(); MailHomeController.websocketConnection = null; }
		}
	},

	sortMailItemList: function(e, column) {
		if (this.get("sortColumn") === column) {
			this.set("sortDirection", (this.get("sortDirection") === "desc") ? "asc" : "desc");
		} else {
			this.set("sortDirection", "desc");
		}

		this.set("sortColumn", column);
	},

	/**
	 * Updates the mail items array for the specified list view.
	 */
	updateMailItemsArray: function(mailItems, view) {
		view.set("mails", mailItems);
	},

	/**
	 * Updates the detail view for a specified mail item.
	 */
	updateMailItemDetailView: function(mailId, subject, dateSent, fromAddress, body, attachments, view) {
		view.set({
			subject: subject,
			dateSent: ((dateSent.length > 0) ? MailService.formatMailDate(dateSent) : ""),
			fromAddress: fromAddress,
			body: body,
			attachments: attachments
		});
	},

	/**
	 * Used to view a single mail item. This will highlight the selected
	 * item and open the detail view.
	 */
	viewMailItem: function(e, view) {
		MailService.getMailItem(e.context.id, function(mailItem) {
			MailHomeController.clearMailItemSelections();
			MailHomeController.selectMailItem(e);

			MailHomeController.updateMailItemDetailView(
				mailItem.id,
				mailItem.subject,
				mailItem.dateSent,
				mailItem.fromAddress,
				mailItem.body,
				mailItem.attachments,
				view
			);
		})
	},

	websocketConnection: null
};
