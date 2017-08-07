'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var FeaturesResult = function () {
  function FeaturesResult(strict) {
    (0, _classCallCheck3.default)(this, FeaturesResult);

    this.duration = 0;
    this.scenarioResults = [];
    this.success = true;
    this.stepResults = [];
    this.strict = strict;
  }

  (0, _createClass3.default)(FeaturesResult, [{
    key: 'witnessScenarioResult',
    value: function witnessScenarioResult(scenarioResult) {
      var duration = scenarioResult.duration,
          status = scenarioResult.status,
          stepResults = scenarioResult.stepResults;

      this.duration += duration;
      this.scenarioResults.push(scenarioResult);
      this.stepResults = this.stepResults.concat(stepResults);
      if (_lodash2.default.includes([_status2.default.AMBIGUOUS, _status2.default.FAILED], status)) {
        this.success = false;
      }
      if (this.strict && _lodash2.default.includes([_status2.default.PENDING, _status2.default.UNDEFINED], status)) {
        this.success = false;
      }
    }
  }]);
  return FeaturesResult;
}();

exports.default = FeaturesResult;