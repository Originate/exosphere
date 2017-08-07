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

var _attachment_manager = require('./attachment_manager');

var _attachment_manager2 = _interopRequireDefault(_attachment_manager);

var _event = require('./event');

var _event2 = _interopRequireDefault(_event);

var _hook = require('../models/hook');

var _hook2 = _interopRequireDefault(_hook);

var _scenario_result = require('../models/scenario_result');

var _scenario_result2 = _interopRequireDefault(_scenario_result);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

var _step_result = require('../models/step_result');

var _step_result2 = _interopRequireDefault(_step_result);

var _step_runner = require('./step_runner');

var _step_runner2 = _interopRequireDefault(_step_runner);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var ScenarioRunner = function () {
  function ScenarioRunner(_ref) {
    var _context;

    var eventBroadcaster = _ref.eventBroadcaster,
        options = _ref.options,
        scenario = _ref.scenario,
        supportCodeLibrary = _ref.supportCodeLibrary;
    (0, _classCallCheck3.default)(this, ScenarioRunner);

    this.attachmentManager = new _attachment_manager2.default();
    this.eventBroadcaster = eventBroadcaster;
    this.options = options;
    this.scenario = scenario;
    this.supportCodeLibrary = supportCodeLibrary;
    this.scenarioResult = new _scenario_result2.default(scenario);
    this.world = new supportCodeLibrary.World({
      attach: (_context = this.attachmentManager).create.bind(_context),
      parameters: options.worldParameters
    });
  }

  (0, _createClass3.default)(ScenarioRunner, [{
    key: 'broadcastScenarioResult',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* () {
        var event = new _event2.default({ data: this.scenarioResult, name: _event2.default.SCENARIO_RESULT_EVENT_NAME });
        yield this.eventBroadcaster.broadcastEvent(event);
      });

      function broadcastScenarioResult() {
        return _ref2.apply(this, arguments);
      }

      return broadcastScenarioResult;
    }()
  }, {
    key: 'broadcastStepResult',
    value: function () {
      var _ref3 = (0, _bluebird.coroutine)(function* (stepResult) {
        this.scenarioResult.witnessStepResult(stepResult);
        var event = new _event2.default({ data: stepResult, name: _event2.default.STEP_RESULT_EVENT_NAME });
        yield this.eventBroadcaster.broadcastEvent(event);
      });

      function broadcastStepResult(_x) {
        return _ref3.apply(this, arguments);
      }

      return broadcastStepResult;
    }()
  }, {
    key: 'invokeStep',
    value: function invokeStep(step, stepDefinition) {
      return _step_runner2.default.run({
        attachmentManager: this.attachmentManager,
        defaultTimeout: this.supportCodeLibrary.defaultTimeout,
        scenarioResult: this.scenarioResult,
        step: step,
        stepDefinition: stepDefinition,
        parameterTypeRegistry: this.supportCodeLibrary.parameterTypeRegistry,
        world: this.world
      });
    }
  }, {
    key: 'isSkippingSteps',
    value: function isSkippingSteps() {
      return this.scenarioResult.status !== _status2.default.PASSED;
    }
  }, {
    key: 'run',
    value: function () {
      var _ref4 = (0, _bluebird.coroutine)(function* () {
        var _this = this;

        var event = new _event2.default({ data: this.scenario, name: _event2.default.SCENARIO_EVENT_NAME });
        yield this.eventBroadcaster.broadcastAroundEvent(event, (0, _bluebird.coroutine)(function* () {
          yield _this.runBeforeHooks();
          yield _this.runSteps();
          yield _this.runAfterHooks();
          yield _this.broadcastScenarioResult();
        }));
        return this.scenarioResult;
      });

      function run() {
        return _ref4.apply(this, arguments);
      }

      return run;
    }()
  }, {
    key: 'runAfterHooks',
    value: function () {
      var _ref6 = (0, _bluebird.coroutine)(function* () {
        yield this.runHooks({
          hookDefinitions: this.supportCodeLibrary.afterHookDefinitions,
          hookKeyword: _hook2.default.AFTER_STEP_KEYWORD
        });
      });

      function runAfterHooks() {
        return _ref6.apply(this, arguments);
      }

      return runAfterHooks;
    }()
  }, {
    key: 'runBeforeHooks',
    value: function () {
      var _ref7 = (0, _bluebird.coroutine)(function* () {
        yield this.runHooks({
          hookDefinitions: this.supportCodeLibrary.beforeHookDefinitions,
          hookKeyword: _hook2.default.BEFORE_STEP_KEYWORD
        });
      });

      function runBeforeHooks() {
        return _ref7.apply(this, arguments);
      }

      return runBeforeHooks;
    }()
  }, {
    key: 'runHook',
    value: function () {
      var _ref8 = (0, _bluebird.coroutine)(function* (hook, hookDefinition) {
        if (this.options.dryRun) {
          return new _step_result2.default({
            step: hook,
            stepDefinition: hookDefinition,
            status: _status2.default.SKIPPED
          });
        } else {
          return yield this.invokeStep(hook, hookDefinition);
        }
      });

      function runHook(_x2, _x3) {
        return _ref8.apply(this, arguments);
      }

      return runHook;
    }()
  }, {
    key: 'runHooks',
    value: function () {
      var _ref10 = (0, _bluebird.coroutine)(function* (_ref9) {
        var _this2 = this;

        var hookDefinitions = _ref9.hookDefinitions,
            hookKeyword = _ref9.hookKeyword;

        yield _bluebird2.default.each(hookDefinitions, function () {
          var _ref11 = (0, _bluebird.coroutine)(function* (hookDefinition) {
            if (!hookDefinition.appliesToScenario(_this2.scenario)) {
              return;
            }
            var hook = new _hook2.default({ keyword: hookKeyword, scenario: _this2.scenario });
            var event = new _event2.default({ data: hook, name: _event2.default.STEP_EVENT_NAME });
            yield _this2.eventBroadcaster.broadcastAroundEvent(event, (0, _bluebird.coroutine)(function* () {
              var stepResult = yield _this2.runHook(hook, hookDefinition);
              yield _this2.broadcastStepResult(stepResult);
            }));
          });

          return function (_x5) {
            return _ref11.apply(this, arguments);
          };
        }());
      });

      function runHooks(_x4) {
        return _ref10.apply(this, arguments);
      }

      return runHooks;
    }()
  }, {
    key: 'runStep',
    value: function () {
      var _ref13 = (0, _bluebird.coroutine)(function* (step) {
        var _this3 = this;

        var stepDefinitions = this.supportCodeLibrary.stepDefinitions.filter(function (stepDefinition) {
          return stepDefinition.matchesStepName({
            stepName: step.name,
            parameterTypeRegistry: _this3.supportCodeLibrary.parameterTypeRegistry
          });
        });
        if (stepDefinitions.length === 0) {
          return new _step_result2.default({
            step: step,
            status: _status2.default.UNDEFINED
          });
        } else if (stepDefinitions.length > 1) {
          return new _step_result2.default({
            ambiguousStepDefinitions: stepDefinitions,
            step: step,
            status: _status2.default.AMBIGUOUS
          });
        } else if (this.options.dryRun || this.isSkippingSteps()) {
          return new _step_result2.default({
            step: step,
            stepDefinition: stepDefinitions[0],
            status: _status2.default.SKIPPED
          });
        } else {
          return yield this.invokeStep(step, stepDefinitions[0]);
        }
      });

      function runStep(_x6) {
        return _ref13.apply(this, arguments);
      }

      return runStep;
    }()
  }, {
    key: 'runSteps',
    value: function () {
      var _ref14 = (0, _bluebird.coroutine)(function* () {
        var _this4 = this;

        yield _bluebird2.default.each(this.scenario.steps, function () {
          var _ref15 = (0, _bluebird.coroutine)(function* (step) {
            var event = new _event2.default({ data: step, name: _event2.default.STEP_EVENT_NAME });
            yield _this4.eventBroadcaster.broadcastAroundEvent(event, (0, _bluebird.coroutine)(function* () {
              var stepResult = yield _this4.runStep(step);
              yield _this4.broadcastStepResult(stepResult);
            }));
          });

          return function (_x7) {
            return _ref15.apply(this, arguments);
          };
        }());
      });

      function runSteps() {
        return _ref14.apply(this, arguments);
      }

      return runSteps;
    }()
  }]);
  return ScenarioRunner;
}();

exports.default = ScenarioRunner;