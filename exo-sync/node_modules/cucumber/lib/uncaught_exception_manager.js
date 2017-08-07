'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var UncaughtExceptionManager = function () {
  function UncaughtExceptionManager() {
    (0, _classCallCheck3.default)(this, UncaughtExceptionManager);
  }

  (0, _createClass3.default)(UncaughtExceptionManager, null, [{
    key: 'registerHandler',
    value: function registerHandler(handler) {
      if (typeof window === 'undefined') {
        process.addListener('uncaughtException', handler);
      } else {
        window.onerror = handler;
      }
    }
  }, {
    key: 'unregisterHandler',
    value: function unregisterHandler(handler) {
      if (typeof window === 'undefined') {
        process.removeListener('uncaughtException', handler);
      } else {
        window.onerror = void 0;
      }
    }
  }]);
  return UncaughtExceptionManager;
}();

exports.default = UncaughtExceptionManager;