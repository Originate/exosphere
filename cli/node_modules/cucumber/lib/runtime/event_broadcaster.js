'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _bluebird = require('bluebird');

var _bluebird2 = _interopRequireDefault(_bluebird);

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _user_code_runner = require('../user_code_runner');

var _user_code_runner2 = _interopRequireDefault(_user_code_runner);

var _verror = require('verror');

var _verror2 = _interopRequireDefault(_verror);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var EventBroadcaster = function () {
  function EventBroadcaster(_ref) {
    var cwd = _ref.cwd,
        listenerDefaultTimeout = _ref.listenerDefaultTimeout,
        listeners = _ref.listeners;
    (0, _classCallCheck3.default)(this, EventBroadcaster);

    this.cwd = cwd;
    this.listenerDefaultTimeout = listenerDefaultTimeout;
    this.listeners = listeners;
  }

  (0, _createClass3.default)(EventBroadcaster, [{
    key: 'broadcastAroundEvent',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* (event, fn) {
        yield this.broadcastEvent(event.buildBeforeEvent());
        yield fn();
        yield this.broadcastEvent(event.buildAfterEvent());
      });

      function broadcastAroundEvent(_x, _x2) {
        return _ref2.apply(this, arguments);
      }

      return broadcastAroundEvent;
    }()
  }, {
    key: 'broadcastEvent',
    value: function broadcastEvent(event) {
      var _this = this;

      return _bluebird2.default.each(this.listeners, function () {
        var _ref3 = (0, _bluebird.coroutine)(function* (listener) {
          var fnName = 'handle' + event.name;
          var handler = listener[fnName];
          if (handler) {
            var timeout = listener.timeout || _this.listenerDefaultTimeout;

            var _ref4 = yield _user_code_runner2.default.run({
              argsArray: [event.data],
              fn: handler,
              thisArg: listener,
              timeoutInMilliseconds: timeout
            }),
                error = _ref4.error;

            if (error) {
              var location = _this.getListenerErrorLocation({ fnName: fnName, listener: listener });
              throw new _verror2.default(error, 'a handler errored, process exiting: ' + location);
            }
          }
        });

        return function (_x3) {
          return _ref3.apply(this, arguments);
        };
      }());
    }
  }, {
    key: 'getListenerErrorLocation',
    value: function getListenerErrorLocation(_ref5) {
      var fnName = _ref5.fnName,
          listener = _ref5.listener;

      return listener.relativeUri || listener.constructor.name + '::' + fnName;
    }
  }]);
  return EventBroadcaster;
}();

exports.default = EventBroadcaster;