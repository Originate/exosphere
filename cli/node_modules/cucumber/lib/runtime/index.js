'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _bluebird = require('bluebird');

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _event_broadcaster = require('./event_broadcaster');

var _event_broadcaster2 = _interopRequireDefault(_event_broadcaster);

var _features_runner = require('./features_runner');

var _features_runner2 = _interopRequireDefault(_features_runner);

var _stack_trace_filter = require('./stack_trace_filter');

var _stack_trace_filter2 = _interopRequireDefault(_stack_trace_filter);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Runtime = function () {
  // options - {dryRun, failFast, filterStacktraces, strict}
  function Runtime(_ref) {
    var features = _ref.features,
        listeners = _ref.listeners,
        options = _ref.options,
        supportCodeLibrary = _ref.supportCodeLibrary;
    (0, _classCallCheck3.default)(this, Runtime);

    this.features = features || [];
    this.listeners = listeners || [];
    this.options = options || {};
    this.supportCodeLibrary = supportCodeLibrary;
    this.stackTraceFilter = new _stack_trace_filter2.default();
  }

  (0, _createClass3.default)(Runtime, [{
    key: 'start',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* () {
        var eventBroadcaster = new _event_broadcaster2.default({
          listenerDefaultTimeout: this.supportCodeLibrary.defaultTimeout,
          listeners: this.listeners.concat(this.supportCodeLibrary.listeners)
        });
        var featuresRunner = new _features_runner2.default({
          eventBroadcaster: eventBroadcaster,
          features: this.features,
          options: this.options,
          supportCodeLibrary: this.supportCodeLibrary
        });

        if (this.options.filterStacktraces) {
          this.stackTraceFilter.filter();
        }

        var result = yield featuresRunner.run();

        if (this.options.filterStacktraces) {
          this.stackTraceFilter.unfilter();
        }

        return result;
      });

      function start() {
        return _ref2.apply(this, arguments);
      }

      return start;
    }()
  }, {
    key: 'attachListener',
    value: function attachListener(listener) {
      this.listeners.push(listener);
    }
  }]);
  return Runtime;
}();

exports.default = Runtime;