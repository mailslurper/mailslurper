
'use strict';

/**
 * UI.Layout
 */
angular.module('ui.layout', [])
  .controller('uiLayoutCtrl', ['$scope', '$attrs', '$element', function uiLayoutCtrl($scope, $attrs, $element) {
    // Gives to the children directives the access to the parent layout.
    return {
      opts: angular.extend({}, $scope.$eval($attrs.uiLayout), $scope.$eval($attrs.options)),
      element: $element
    };
  }])

  .directive('uiLayout', ['$parse', function ($parse) {

    var splitBarElem_htmlTemplate = '<div class="stretch ui-splitbar"></div>';

    function convertNumericDataTypesToPencents(numberVairousTypeArray, minSizes, maxSizes, parentSize){
      var _i, _n;
      var _res = []; _res.length = numberVairousTypeArray.length;
      var _commonSizeIndex = [];
      var _minSizes = [];
      var _maxSizes = [];
      var _remainingSpace = 100;

      for (_i = 0, _n = numberVairousTypeArray.length; _i < _n; ++_i) {
        var minSize = parseInt(minSizes[_i], 10);
        if(minSize) {
          var minType = minSizes[_i].match(/\d+\s*(px|%)\s*$/i);
          if(!isNaN(minSize) && minType) {
            if(minType.length > 1 && 'px' === minType[1]) {
               minSize = + (minSize / parentSize * 100).toFixed(5);
            }
          }
        }
        _minSizes.push(minSize);

        var maxSize = parseInt(maxSizes[_i], 10);
        if(maxSize) {
          var maxType = maxSizes[_i].match(/\d+\s*(px|%)\s*$/i);
          if(!isNaN(maxSize) && maxType) {
            if(maxType.length > 1 && 'px' === maxType[1]) {
               maxSize = + (maxSize / parentSize * 100).toFixed(5);
            }
          }
        }
        _maxSizes.push(maxSize);

        var rawSize = numberVairousTypeArray[_i];
        var value = parseInt(rawSize, 10);
        // should only support pixels and pencent data type
        var type = rawSize.match(/\d+\s*(px|%)\s*$/i);
        if (!isNaN(value) && type){
          if (type.length > 1 && 'px' === type[1]){
            value = + (value / parentSize * 100).toFixed(5);
          }

          if(minSize) value = Math.max(value, minSize);
          if(maxSize) value = Math.min(value, maxSize);

          _res[_i] = value;
          _remainingSpace -= value;
        }  else {
          rawSize = 'auto';
        }

        if (/^\s*auto\s*$/.test(rawSize)){
          _commonSizeIndex.push(_i); continue;
        }
      }

      if (_commonSizeIndex.length > 0){
        var _commonSize = _remainingSpace / _commonSizeIndex.length;
        var _modifiedSizeIndex = [];
        // Check to make sure the common size isn't outside the bounds of min/max-size
        for (_i = 0, _n = _commonSizeIndex.length; _i < _n; ++_i) {
          var _cid = _commonSizeIndex[_i];
          var _minSize = _minSizes[_cid];
          var _maxSize = _maxSizes[_cid];
          if(_commonSize <= _minSize) {
            _remainingSpace -= _minSize;
            _res[_cid] = _minSize;
          } else if(_commonSize >= _maxSize) {
            _remainingSpace -= _maxSize;
            _res[_cid] = _maxSize;
          } else {
            _modifiedSizeIndex.push(_cid);
          }
        }
        _commonSize = _remainingSpace / _modifiedSizeIndex.length;
        for (_i = 0, _n = _modifiedSizeIndex.length; _i < _n; ++_i) {
          var cid = _modifiedSizeIndex[_i];
          if(cid !== null) {
            _res[cid] = _commonSize;
          }
        }
      }

      parentSize;

      return _res;
    }

    return {
      restrict: 'AE',
      compile: function compile(tElement, tAttrs) {

        var _i, _childens = tElement.children(), _child_len = _childens.length;
        var _sizes, _position;

        // Parse `ui-layout` or `options` attributes (with no scope...)
        var opts = angular.extend({}, $parse(tAttrs.uiLayout)(), $parse(tAttrs.options)());
        var isUsingColumnFlow = opts.flow === 'column';

        tElement
          // Force the layout to fill the parent space
          // fix no height layout...
          .addClass('stretch')
          // set the layout css class
          .addClass('ui-layout-' + (opts.flow || 'row'));

        // Initial global size definition
        opts.sizes = opts.sizes || [];
        opts.maxSizes = opts.maxSizes || [];
        opts.minSizes = opts.minSizes || [];
        opts.dividerSize = opts.dividerSize || '10px';

        // Preallocate the array size
        opts.sizes.length = _child_len;

        for (_i = 0; _i < _child_len; ++_i) {
          // Stretch all the children by default
          angular.element(_childens[_i]).addClass('stretch');
          // Size initialization priority
          // - the size attr on the child element
          // - the global size on the layout option
          // - 'auto' Fair separation of the remaining space
          opts.sizes[_i] = angular.element(_childens[_i]).attr('size') || opts.sizes[_i]  || 'auto';
          opts.maxSizes[_i] = angular.element(_childens[_i]).attr('max-size') || opts.maxSizes[_i] || null;
          opts.minSizes[_i] = angular.element(_childens[_i]).attr('min-size') || opts.minSizes[_i] || null;
        }

        // get the final percent sizes
        _sizes = convertNumericDataTypesToPencents(opts.sizes, opts.minSizes, opts.maxSizes, tElement[0]['offset' + (isUsingColumnFlow ? 'Width' : 'Height')]);

        if (_child_len > 1) {
          // Initialise the layout with equal sizes.
          var flowProperty = ( isUsingColumnFlow ? 'left' : 'top');
          var oppositeFlowProperty = ( isUsingColumnFlow ? 'right' : 'bottom');
          var sizeProperty = ( isUsingColumnFlow ? 'width' : 'height');
          _position = 0;
          for (_i = 0; _i < _child_len; ++_i) {
            var area = angular.element(_childens[_i])
              .css(flowProperty, _position + '%');

            _position += _sizes[_i];
            area.css(oppositeFlowProperty, (100 - _position) + '%');

            if (_i < _child_len - 1) {
              // Add a split bar
              var bar = angular.element(splitBarElem_htmlTemplate).css(flowProperty, _position + '%');
              bar.css(sizeProperty, opts.dividerSize);
              area.after(bar);
            }
          }
        }
      },
      controller: 'uiLayoutCtrl'
    };

  }])


  .directive('uiSplitbar', function () {

    // Get all the page.
    var htmlElement = angular.element(document.body.parentElement);

    return {
      require: '^uiLayout',
      restrict: 'EAC',
      link: function (scope, iElement, iAttrs, parentLayout) {

        var animationFrameRequested, lastX;
        var _cache = {};

        // Use relative mouse position
        var isUsingColumnFlow = parentLayout.opts.flow === 'column';
        var mouseProperty = ( isUsingColumnFlow ? 'clientX' : 'clientY');

        // Use bounding box / css property names
        var flowProperty = ( isUsingColumnFlow ? 'left' : 'top');
        var oppositeFlowProperty = ( isUsingColumnFlow ? 'right' : 'bottom');
        var sizeProperty = ( isUsingColumnFlow ? 'width' : 'height');

        // Use bounding box properties
        var barElm = iElement[0];
        var previousElement = barElm.previousElementSibling;
        var nextElement = barElm.nextElementSibling;

        // Stores the layout values for some seconds to not recalculate it during the animation
        function _cached_layout_values() {
          // layout bounding box
          var layout_bb = parentLayout.element[0].getBoundingClientRect();

          // split bar bounding box
          var bar_bb = barElm.getBoundingClientRect();

          _cache.time = +new Date();
          _cache.barSize = bar_bb[sizeProperty];
          _cache.layoutSize = layout_bb[sizeProperty];
          _cache.layoutOrigine = layout_bb[flowProperty];
          _cache.previousElement = previousElement.getBoundingClientRect();
          _cache.previousElement.min = parseInt(previousElement.getAttribute('min-size'),10);
          _cache.previousElement.max = parseInt(previousElement.getAttribute('max-size'),10);
          _cache.nextElement = nextElement.getBoundingClientRect();
          _cache.nextElement.min = parseInt(nextElement.getAttribute('min-size'),10);
          _cache.nextElement.max = parseInt(nextElement.getAttribute('max-size'),10);

          var dividerSize = isNaN(bar_bb[sizeProperty]) ? bar_bb[sizeProperty] : bar_bb[sizeProperty] + 'px';
          var _dividerSize = parseInt(dividerSize, 10);
          var _dividerType = dividerSize.match(/\d+\s*(px|%)\s*$/i);
          if(!isNaN(_dividerSize) && _dividerType) {
            if(_dividerType.length > 1 && 'px' === _dividerType[1]) {
              _dividerSize = + (_dividerSize / _cache.layoutSize * 100).toFixed(5);
            }
          }

          if(_cache.previousElement.min) {
            var minType = previousElement.getAttribute('min-size').match(/\d+\s*(px|%)\s*$/i);
            if(!isNaN(_cache.previousElement.min) && minType) {
              if(minType.length > 1 && 'px' === minType[1]) {
                 _cache.previousElement.min = + (_cache.previousElement.min / _cache.layoutSize * 100).toFixed(5);
              }
            }
            // ensure the min size isn't smaller than the divider size
            if(_dividerSize && _cache.previousElement.min < _dividerSize) _cache.previousElement.min = _dividerSize;
          } else {
            _cache.previousElement.min = _dividerSize;
          }

          if( _cache.previousElement.max) {
            var maxType = previousElement.getAttribute('max-size').match(/\d+\s*(px|%)\s*$/i);
            if(!isNaN( _cache.previousElement.max) && maxType) {
              if(maxType.length > 1 && 'px' === maxType[1]) {
                  _cache.previousElement.max = + ( _cache.previousElement.max / _cache.layoutSize * 100).toFixed(5);
              }
            }
          }

          if(_cache.nextElement.min) {
            var _minType = nextElement.getAttribute('min-size').match(/\d+\s*(px|%)\s*$/i);
            if(!isNaN(_cache.nextElement.min) && _minType) {
              if(_minType.length > 1 && 'px' === _minType[1]) {
                 _cache.nextElement.min = + (_cache.nextElement.min / _cache.layoutSize * 100).toFixed(5);
              }
            }
            // ensure the min size isn't smaller than the divider size
            if(_dividerSize && _cache.nextElement.min < _dividerSize) _cache.nextElement.min = _dividerSize;
          } else {
            _cache.nextElement.min = _dividerSize;
          }

          if(_cache.nextElement.max) {
            var _maxType = nextElement.getAttribute('max-size').match(/\d+\s*(px|%)\s*$/i);
            if(!isNaN(_cache.nextElement.max) && _maxType) {
              if(_maxType.length > 1 && 'px' === _maxType[1]) {
                 _cache.nextElement.max = + (_cache.nextElement.max / _cache.layoutSize * 100).toFixed(5);
              }
            }
          }
        }

        function _draw() {
          var the_pos = (lastX - _cache.layoutOrigine) / _cache.layoutSize * 100;

          // Keep the bar in the window (no left/top 100%)
          the_pos = Math.min(the_pos, 100 - _cache.barSize / _cache.layoutSize * 100);

          // Keep the bar from going past the previous elements max/min sizes
          var previousElementValue = _cache.previousElement[flowProperty] / _cache.layoutSize * 100;
          if(!isNaN(_cache.previousElement.min) && the_pos < (previousElementValue + _cache.previousElement.min)) the_pos = (previousElementValue + _cache.previousElement.min);
          if(!isNaN(_cache.previousElement.max) && the_pos > (previousElementValue + _cache.previousElement.max)) the_pos = (previousElementValue + _cache.previousElement.max);

          // Keep the bar from going past the next elements max/min sizes
          var nextElementValue = (_cache.nextElement[flowProperty] + _cache.nextElement[sizeProperty]) / _cache.layoutSize * 100;
          var nextElementMinValue = nextElementValue - _cache.nextElement.min;
          var nextElementMaxValue = nextElementValue - _cache.nextElement.max;
          if(!isNaN(_cache.nextElement.max) && the_pos < nextElementMaxValue) the_pos = nextElementMaxValue;
          if(!isNaN(_cache.nextElement.min) && the_pos > nextElementMinValue) the_pos = nextElementMinValue;

          // The the bar in the near beetween the near by area
          the_pos = Math.max(the_pos, parseInt(barElm.previousElementSibling.style[flowProperty], 10));
          if (barElm.nextElementSibling.nextElementSibling) {
            the_pos = Math.min(the_pos, parseInt(barElm.nextElementSibling.nextElementSibling.style[flowProperty], 10));
          }

          // change the position of the bar and the next area
          barElm.style[flowProperty] = barElm.nextElementSibling.style[flowProperty] = the_pos + '%';
          // change the size of the previous area
          barElm.previousElementSibling.style[oppositeFlowProperty] = (100 - the_pos) + '%';

          // Enable a new animation frame
          animationFrameRequested = null;
        }

        function _resize(mouseEvent) {
          // Store the mouse position for later

          // FIX :
          // - with touch events, when using jQuery, the mouseEvent is in fact a jQueryEvent. So we use originalEvent here.
          // - real touch events comes in the _targetTouches_ array
          lastX = mouseEvent[mouseProperty] ||
            (mouseEvent.originalEvent && mouseEvent.originalEvent[mouseProperty]) ||
            (mouseEvent.targetTouches ? mouseEvent.targetTouches[0][mouseProperty] : 0);

          // Cancel previous rAF call
          if (animationFrameRequested) {
            window.cancelAnimationFrame(animationFrameRequested);
          }

          if (!_cache.time || +new Date() > _cache.time + 1000) { // after ~60 frames
            _cached_layout_values();
          }

          // Animate the page outside the event
          animationFrameRequested = window.requestAnimationFrame(_draw);
        }


        // Bind the click on the bar then you can move it all over the page.
        iElement.on('mousedown touchstart', function (e) {
          e.preventDefault();
          e.stopPropagation();
          htmlElement.on('mousemove touchmove', _resize);
          return false;
        });
        htmlElement.on('mouseup touchend', function () {
          htmlElement.off('mousemove touchmove');
        });
      }
    };
  });
