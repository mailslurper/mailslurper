app.directive("mailLayout", function() {
	return {
		scope: {},
		link: function(scope, el, attrs) {
			$(window.body).layout({
				north: {
					size: 35,
					resizable: false,
					closable: false
				},
				south: {
					resizable: false,
					closable: false,
					size: 40
				},
				east: {
					resizable: true,
					closable: true,
					size: "40%"
				}
			});
		}
	};
});