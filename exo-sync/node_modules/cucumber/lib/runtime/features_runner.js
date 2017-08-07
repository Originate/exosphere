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

var _event = require('./event');

var _event2 = _interopRequireDefault(_event);

var _features_result = require('../models/features_result');

var _features_result2 = _interopRequireDefault(_features_result);

var _scenario_result = require('../models/scenario_result');

var _scenario_result2 = _interopRequireDefault(_scenario_result);

var _scenario_runner = require('./scenario_runner');

var _scenario_runner2 = _interopRequireDefault(_scenario_runner);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var FeaturesRunner = function () {
  function FeaturesRunner(_ref) {
    var eventBroadcaster = _ref.eventBroadcaster,
        features = _ref.features,
        options = _ref.options,
        supportCodeLibrary = _ref.supportCodeLibrary;
    (0, _classCallCheck3.default)(this, FeaturesRunner);

    this.eventBroadcaster = eventBroadcaster;
    this.features = features;
    this.options = options;
    this.supportCodeLibrary = supportCodeLibrary;
    this.featuresResult = new _features_result2.default(options.strict);
  }

  (0, _createClass3.default)(FeaturesRunner, [{
    key: 'run',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* () {
        var _this = this;

        var event = new _event2.default({ data: this.features, name: _event2.default.FEATURES_EVENT_NAME });
        yield this.eventBroadcaster.broadcastAroundEvent(event, (0, _bluebird.coroutine)(function* () {
          yield _bluebird2.default.each(_this.features, _this.runFeature.bind(_this));
          yield _this.broadcastFeaturesResult();
        }));
        return this.featuresResult.success;
      });

      function run() {
        return _ref2.apply(this, arguments);
      }

      return run;
    }()
  }, {
    key: 'broadcastFeaturesResult',
    value: function () {
      var _ref4 = (0, _bluebird.coroutine)(function* () {
        var event = new _event2.default({ data: this.featuresResult, name: _event2.default.FEATURES_RESULT_EVENT_NAME });
        yield this.eventBroadcaster.broadcastEvent(event);
      });

      function broadcastFeaturesResult() {
        return _ref4.apply(this, arguments);
      }

      return broadcastFeaturesResult;
    }()
  }, {
    key: 'runFeature',
    value: function () {
      var _ref5 = (0, _bluebird.coroutine)(function* (feature) {
        var _this2 = this;

        var event = new _event2.default({ data: feature, name: _event2.default.FEATURE_EVENT_NAME });
        yield this.eventBroadcaster.broadcastAroundEvent(event, (0, _bluebird.coroutine)(function* () {
          yield _bluebird2.default.each(feature.scenarios, function () {
            var _ref7 = (0, _bluebird.coroutine)(function* (scenario) {
              var scenarioResult = yield _this2.runScenario(scenario);
              _this2.featuresResult.witnessScenarioResult(scenarioResult);
            });

            return function (_x2) {
              return _ref7.apply(this, arguments);
            };
          }());
        }));
      });

      function runFeature(_x) {
        return _ref5.apply(this, arguments);
      }

      return runFeature;
    }()
  }, {
    key: 'runScenario',
    value: function () {
      var _ref8 = (0, _bluebird.coroutine)(function* (scenario) {
        if (!this.featuresResult.success && this.options.failFast) {
          return new _scenario_result2.default(scenario, _status2.default.SKIPPED);
        }
        var scenarioRunner = new _scenario_runner2.default({
          eventBroadcaster: this.eventBroadcaster,
          options: this.options,
          scenario: scenario,
          supportCodeLibrary: this.supportCodeLibrary
        });
        return yield scenarioRunner.run();
      });

      function runScenario(_x3) {
        return _ref8.apply(this, arguments);
      }

      return runScenario;
    }()
  }]);
  return FeaturesRunner;
}();

exports.default = FeaturesRunner;