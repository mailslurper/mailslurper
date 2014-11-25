(function (angular) {
    "use strict";
    angular.module('flexyLayout.block', [])
        .provider('Block', function () {

            /**
             * A composite block made of different types of blocks that must implement the structural interface
             *
             * moveLength->change the lengthValue according to specific block rules
             * canMoveLength->tells whether the block can change his lengthValue in the current state
             * getAvailableLength->return the length the block can be reduced of
             *
             * , canMoveLength, getAvailableLength
             * @param composingBlocks
             * @constructor
             */
            function CompositeBlock(composingBlocks) {
                this.blocks = [];

                if (angular.isArray(composingBlocks)) {
                    for (var i = 0, l = composingBlocks.length; i < l; i++) {
                        //should implement structural interface
                        if (composingBlocks[i].moveLength && composingBlocks[i].canMoveLength && composingBlocks[i].getAvailableLength) {
                            this.blocks.push(composingBlocks[i]);
                        }
                    }
                }
            }

            CompositeBlock.prototype.moveLength = function (length) {

                var
                    divider = 0,
                    initialLength = length,
                    blockLength;

                for (var i = 0, l = this.blocks.length; i < l; i++) {
                    if (this.blocks[i].canMoveLength(length) === true) {
                        divider++;
                    }
                }

                for (var j = 0; divider > 0; j++) {
                    blockLength = this.blocks[j].moveLength(length / divider);
                    length -= blockLength;
                    if (Math.abs(blockLength) > 0) {
                        divider--;
                    }
                }

                return initialLength - length;
            };

            CompositeBlock.prototype.canMoveLength = function (length) {

                for (var i = 0, l = this.blocks.length; i < l; i++) {
                    if (this.blocks[i].canMoveLength(length) === true) {
                        return true;
                    }
                }

                return false;
            };

            CompositeBlock.prototype.getAvailableLength = function () {
                var length = 0;
                for (var i = 0, l = this.blocks.length; i < l; i++) {
                    length += this.blocks[i].getAvailableLength();
                }

                return length;
            };

            CompositeBlock.prototype.clean = function () {
                delete this.blocks;
            };

            /**
             * A Blokc which can be locked (ie its lengthValue can not change) this is the standard composing block
             * @constructor
             */
            function Block(initial) {
                this.initialLength = initial > 0 ? initial : 0;
                this.isLocked = false;
                this.lengthValue = 0;
                this.minLength = 0;
            }

            Block.prototype.moveLength = function (length) {

                if (this.isLocked === true) {
                    return 0;
                }

                var oldLength = this.lengthValue;
                if (angular.isNumber(length)) {
                    this.lengthValue = Math.max(0, this.lengthValue + length);
                }
                return this.lengthValue - oldLength;
            };

            Block.prototype.canMoveLength = function (length) {
                return !(this.isLocked === true || (length < 0 && (this.getAvailableLength()) === 0));
            };

            Block.prototype.getAvailableLength = function () {
                return this.isLocked === true ? 0 : this.lengthValue - this.minLength;
            };

            /**
             * Splitter a splitter block which split a set of blocks into two separate set
             * @constructor
             */
            function Splitter() {
                this.lengthValue = 5;
                this.initialPosition = { x: 0, y: 0};
                this.availableLength = {before: 0, after: 0};
                this.ghostPosition = { x: 0, y: 0};

            }

            Splitter.prototype.canMoveLength = function () {
                return false;
            };

            Splitter.prototype.moveLength = function () {
                return 0;
            };

            Splitter.prototype.getAvailableLength = function () {
                return 0;
            };

            this.$get = function () {
                return {
                    //variadic -> can call getNewComposite([block1, block2, ...]) or getNewComposite(block1, block2, ...)
                    getNewComposite: function () {
                        var args = [].slice.call(arguments);
                        if (args.length === 1 && angular.isArray(args[0])) {
                            args = args[0];
                        }
                        return new CompositeBlock(args);
                    },
                    getNewBlock: function (initialLength) {
                        return new Block(initialLength);
                    },
                    getNewSplitter: function () {
                        return new Splitter();
                    },

                    isSplitter: function (block) {
                        return block instanceof Splitter;
                    }
                };
            }
        });
})(angular);
(function (angular) {
    "use strict";
    angular.module('flexyLayout.directives', ['flexyLayout.mediator'])
        .directive('flexyLayout', function () {
            return {
                restrict: 'E',
                scope: {},
                template: '<div class="flexy-layout" ng-transclude></div>',
                replace: true,
                transclude: true,
                controller: 'mediatorCtrl',
                link: function (scope, element, attrs, ctrl) {
                    scope.$watch(function () {
                        return element[0][ctrl.lengthProperties.offsetName];
                    }, function () {
                        ctrl.init();
                    });
                }
            };
        })
        .directive('blockContainer', ['Block', function (Block) {
            return{
                restrict: 'E',
                require: '^flexyLayout',
                transclude: true,
                replace: true,
                scope: {},
                template: '<div class="block">' +
                    '<div class="block-content" ng-transclude>' +
                    '</div>' +
                    '</div>',
                link: function (scope, element, attrs, ctrl) {
                    var initialLength = scope.$eval(attrs.init);
                    scope.block = Block.getNewBlock(initialLength);
                    scope.$watch('block.lengthValue', function (newValue, oldValue) {
                        element.css(ctrl.lengthProperties.lengthName, Math.floor(newValue) + 'px');
                    });

                    ctrl.addBlock(scope.block);
                }
            };
        }])
        .directive('blockSplitter', ['Block', function (Block) {
            return{
                restrict: 'E',
                require: '^flexyLayout',
                replace: true,
                scope: {},
                template: '<div class="block splitter">' +
                    '<div class="ghost"></div>' +
                    '</div>',
                link: function (scope, element, attrs, ctrl) {
                    scope.splitter = Block.getNewSplitter();

                    var ghost = element.children()[0];
                    var mouseDownHandler = function (event) {
                        this.initialPosition.x = event.clientX;
                        this.initialPosition.y = event.clientY;
                        this.availableLength = ctrl.getSplitterRange(this);
                        ctrl.movingSplitter = this;

                        //to avoid the block content to be selected when dragging the splitter
                        event.preventDefault();
                    };

                    ctrl.addBlock(scope.splitter);

                    element.bind('mousedown', angular.bind(scope.splitter, mouseDownHandler));

                    scope.$watch('splitter.ghostPosition.' + ctrl.lengthProperties.position, function (newValue, oldValue) {
                        if (newValue !== oldValue) {
                            ghost.style[ctrl.lengthProperties.positionName] = newValue + 'px';
                        }
                    });

                }
            };
        }]);

    angular.module('flexyLayout', ['flexyLayout.directives']);

})(angular);
(function (angular) {
    "use strict";
    //TODO this guy is now big, split it, maybe the part for event handling should be moved somewhere else
    angular.module('flexyLayout.mediator', ['flexyLayout.block']).
        controller('mediatorCtrl', ['$scope', '$element', '$attrs', 'Block', function (scope, element, attrs, Block) {

            var blocks = [],
                pendingSplitter = null,
                splitterCount = 0,
                self = this,
                possibleOrientations = ['vertical', 'horizontal'],
                orientation = possibleOrientations.indexOf(attrs.orientation) !== -1 ? attrs.orientation : 'horizontal',
                className = orientation === 'horizontal' ? 'flexy-layout-column' : 'flexy-layout-row';

            element.addClass(className);

            this.lengthProperties = orientation === 'horizontal' ? {lengthName: 'width', offsetName: 'offsetWidth', positionName: 'left', position: 'x', eventProperty: 'clientX'} :
            {lengthName: 'height', offsetName: 'offsetHeight', positionName: 'top', position: 'y', eventProperty: 'clientY'};

            ///// mouse event handler /////

            this.movingSplitter = null;

            var mouseMoveHandler = function (event) {
                var length = 0,
                    eventProperty = this.lengthProperties.eventProperty,
                    position = this.lengthProperties.position;

                if (this.movingSplitter !== null) {
                    length = event[eventProperty] - this.movingSplitter.initialPosition[position];
                    if (length < 0) {
                        this.movingSplitter.ghostPosition[position] = (-1) * Math.min(Math.abs(length), this.movingSplitter.availableLength.before);
                    } else {
                        this.movingSplitter.ghostPosition[position] = Math.min(length, this.movingSplitter.availableLength.after);
                    }
                }
            };

            var mouseUpHandler = function (event) {
                var length = 0,
                    eventProperty = this.lengthProperties.eventProperty,
                    position = this.lengthProperties.position;

                if (this.movingSplitter !== null) {
                    length = event[eventProperty] - this.movingSplitter.initialPosition[position];
                    this.moveSplitterLength(this.movingSplitter, length);
                    this.movingSplitter.ghostPosition[position] = 0;
                    this.movingSplitter = null;
                }
            };

            element.bind('mouseup', function (event) {
                scope.$apply(angular.bind(self, mouseUpHandler, event));
            });

            //todo should do some throttle before calling apply
            element.bind('mousemove', function (event) {
                scope.$apply(angular.bind(self, mouseMoveHandler, event));
            });

            /////   adding blocks   ////

            this.addBlock = function (block) {

                if (!Block.isSplitter(block)) {
                    if (pendingSplitter !== null) {
                        blocks.push(pendingSplitter);
                        splitterCount++;
                        pendingSplitter = null;
                    }

                    blocks.push(block);
                    this.init();
                } else {
                    pendingSplitter = block;
                }
            };

            /**
             * to be called when flexy-layout container has been resized
             */
            this.init = function () {

                var i,
                    l = blocks.length,
                    elementLength = element[0][this.lengthProperties.offsetName],
                    block,
                    bufferBlock = Block.getNewBlock();//temporary buffer block

                blocks.push(bufferBlock);

                //reset all blocks
                for (i = 0; i < l; i++) {
                    block = blocks[i];
                    block.isLocked = false;
                    if (!Block.isSplitter(block)) {
                        block.moveLength(-10000);
                    }
                }
                //buffer block takes all available space
                bufferBlock.moveLength(elementLength - splitterCount * 5);

                for (i = 0; i < l; i++) {
                    block = blocks[i];
                    if (block.initialLength > 0) {
                        this.moveBlockLength(block, block.initialLength);
                        block.isLocked=true;
                    }
                }

                //buffer block free space for non fixed block
                this.moveBlockLength(bufferBlock, -10000);

                for (i = 0; i < l; i++) {
                    blocks[i].isLocked = false;
                }

                blocks.splice(l, 1);

            };

            ///// public api /////

            /**
             * Will move a given block length from @length
             *
             * @param block can be a block or an index (likely index of the block)
             * @param length < 0 or > 0 : decrease/increase block size of abs(length) px
             */
            this.moveBlockLength = function (block, length) {

                var
                    blockIndex = typeof block !== 'object' ? block : blocks.indexOf(block),
                    composingBlocks,
                    composite,
                    availableLength,
                    blockToMove;


                if (blockIndex < 0 || length === 0 || blockIndex >= blocks.length) {
                    return;
                }

                blockToMove = blocks[blockIndex];

                composingBlocks = (blocks.slice(0, blockIndex)).concat(blocks.slice(blockIndex + 1, blocks.length));
                composite = Block.getNewComposite(composingBlocks);

                if (composite.canMoveLength(-length) !== true || blockToMove.canMoveLength(length) !== true) {
                    return;
                }

                if (length < 0) {
                    availableLength = (-1) * blockToMove.moveLength(length);
                    composite.moveLength(availableLength);
                } else {
                    availableLength = (-1) * composite.moveLength(-length);
                    blockToMove.moveLength(availableLength);
                }

                //free memory
                composite.clean();
            };

            /**
             * move splitter it will affect all the blocks before until the previous/next splitter or the edge of area
             * @param splitter
             * @param length
             */
                //todo mutualise with moveBlockLength
            this.moveSplitterLength = function (splitter, length) {

                var
                    splitterIndex = blocks.indexOf(splitter),
                    beforeComposite,
                    afterComposite,
                    availableLength;

                if (!Block.isSplitter(splitter) || splitterIndex === -1) {
                    return;
                }

                beforeComposite = Block.getNewComposite(fromSplitterToSplitter(splitter, true));
                afterComposite = Block.getNewComposite(fromSplitterToSplitter(splitter, false));

                if (!beforeComposite.canMoveLength(length) || !afterComposite.canMoveLength(-length)) {
                    return;
                }

                if (length < 0) {
                    availableLength = (-1) * beforeComposite.moveLength(length);
                    afterComposite.moveLength(availableLength);
                } else {
                    availableLength = (-1) * afterComposite.moveLength(-length);
                    beforeComposite.moveLength(availableLength);
                }

                afterComposite.clean();
                beforeComposite.clean();

            };

            /**
             * return an object with the available length before the splitter and after the splitter
             * @param splitter
             * @returns {{before: *, after: *}}
             */
            this.getSplitterRange = function (splitter) {

                var
                    beforeSplitter = fromSplitterToSplitter(splitter, true),
                    afterSplitter = fromSplitterToSplitter(splitter, false),
                    toReturn = {
                        before: beforeSplitter.getAvailableLength(),
                        after: afterSplitter.getAvailableLength()
                    };

                beforeSplitter.clean();
                afterSplitter.clean();

                return toReturn;
            };

            /**
             * lock/unlock a given block
             * @param block block or blockIndex
             * @param lock new value for block.isLocked
             */
            this.toggleLockBlock = function (block, lock) {
                var
                    blockIndex = typeof block !== 'object' ? block : blocks.indexOf(block),
                    blockToLock;

                if (blockIndex >= 0 && blockIndex < blocks.length) {
                    blockToLock = blocks[blockIndex];
                    blockToLock.isLocked = lock;
                }

            };

            var fromSplitterToSplitter = function (splitter, before) {

                var
                    splitterIndex = blocks.indexOf(splitter),
                    blockGroup = before === true ? blocks.slice(0, splitterIndex) : blocks.slice(splitterIndex + 1, blocks.length),
                    fn = before === true ? Array.prototype.pop : Array.prototype.shift,
                    composite = [],
                    testedBlock;

                while (testedBlock = fn.apply(blockGroup)) {
                    if (Block.isSplitter(testedBlock)) {
                        break;
                    } else {
                        composite.push(testedBlock);
                    }
                }
                return Block.getNewComposite(composite);
            };
        }]);
})(angular);