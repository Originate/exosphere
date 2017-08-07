'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var ScenarioResult = function () {
  function ScenarioResult(scenario, status) {
    (0, _classCallCheck3.default)(this, ScenarioResult);

    this.duration = 0;
    this.failureException = null;
    this.scenario = scenario;
    this.status = status || _status2.default.PASSED;
    this.stepResults = [];
  }

  (0, _createClass3.default)(ScenarioResult, [{
    key: 'shouldUpdateStatus',
    value: function shouldUpdateStatus(stepResultStatus) {
      switch (stepResultStatus) {
        case _status2.default.FAILED:
          return true;
        case _status2.default.AMBIGUOUS:
        case _status2.default.PENDING:
        case _status2.default.SKIPPED:
        case _status2.default.UNDEFINED:
          return this.status === _status2.default.PASSED;
        default:
          return false;
      }
    }
  }, {
    key: 'witnessStepResult',
    value: function witnessStepResult(stepResult) {
      var duration = stepResult.duration,
          failureException = stepResult.failureException,
          status = stepResult.status;

      if (duration) {
        this.duration += duration;
      }
      if (status === _status2.default.FAILED) {
        this.failureException = failureException;
      }
      if (this.shouldUpdateStatus(status)) {
        this.status = status;
      }
      this.stepResults.push(stepResult);
    }
  }]);
  return ScenarioResult;
}();

exports.default = ScenarioResult;


(0, _status.addStatusPredicates)(ScenarioResult.prototype);