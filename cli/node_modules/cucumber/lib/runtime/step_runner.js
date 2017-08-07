'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _bluebird = require('bluebird');

var _bluebird2 = _interopRequireDefault(_bluebird);

var run = function () {
  var _ref2 = (0, _bluebird.coroutine)(function* (_ref) {
    var attachmentManager = _ref.attachmentManager,
        defaultTimeout = _ref.defaultTimeout,
        scenarioResult = _ref.scenarioResult,
        step = _ref.step,
        stepDefinition = _ref.stepDefinition,
        parameterTypeRegistry = _ref.parameterTypeRegistry,
        world = _ref.world;

    beginTiming();
    var error = void 0,
        result = void 0,
        parameters = void 0;

    try {
      parameters = yield _bluebird2.default.all(stepDefinition.getInvocationParameters({ scenarioResult: scenarioResult, step: step, parameterTypeRegistry: parameterTypeRegistry }));
    } catch (err) {
      error = err;
    }

    if (!error) {
      var timeoutInMilliseconds = stepDefinition.options.timeout || defaultTimeout;

      var validCodeLengths = stepDefinition.getValidCodeLengths(parameters);
      if (_lodash2.default.includes(validCodeLengths, stepDefinition.code.length)) {
        var data = yield _user_code_runner2.default.run({
          argsArray: parameters,
          fn: stepDefinition.code,
          thisArg: world,
          timeoutInMilliseconds: timeoutInMilliseconds
        });
        error = data.error;
        result = data.result;
      } else {
        error = stepDefinition.getInvalidCodeLengthMessage(parameters);
      }
    }

    var attachments = attachmentManager.getAll();
    attachmentManager.reset();

    var stepResultData = {
      attachments: attachments,
      duration: endTiming(),
      step: step,
      stepDefinition: stepDefinition
    };

    if (result === 'pending') {
      stepResultData.status = _status2.default.PENDING;
    } else if (error) {
      stepResultData.failureException = error;
      stepResultData.status = _status2.default.FAILED;
    } else {
      stepResultData.status = _status2.default.PASSED;
    }

    return new _step_result2.default(stepResultData);
  });

  return function run(_x) {
    return _ref2.apply(this, arguments);
  };
}();

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

var _step_result = require('../models/step_result');

var _step_result2 = _interopRequireDefault(_step_result);

var _time = require('../time');

var _time2 = _interopRequireDefault(_time);

var _user_code_runner = require('../user_code_runner');

var _user_code_runner2 = _interopRequireDefault(_user_code_runner);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var beginTiming = _time2.default.beginTiming,
    endTiming = _time2.default.endTiming;
exports.default = { run: run };