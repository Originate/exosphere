'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _possibleConstructorReturn2 = require('babel-runtime/helpers/possibleConstructorReturn');

var _possibleConstructorReturn3 = _interopRequireDefault(_possibleConstructorReturn2);

var _inherits2 = require('babel-runtime/helpers/inherits');

var _inherits3 = _interopRequireDefault(_inherits2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _helpers = require('./helpers');

var _2 = require('./');

var _3 = _interopRequireDefault(_2);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var SummaryFormatter = function (_Formatter) {
  (0, _inherits3.default)(SummaryFormatter, _Formatter);

  function SummaryFormatter() {
    (0, _classCallCheck3.default)(this, SummaryFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (SummaryFormatter.__proto__ || Object.getPrototypeOf(SummaryFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(SummaryFormatter, [{
    key: 'handleFeaturesResult',
    value: function handleFeaturesResult(featuresResult) {
      var failures = featuresResult.stepResults.filter(function (stepResult) {
        return _lodash2.default.includes([_status2.default.AMBIGUOUS, _status2.default.FAILED], stepResult.status);
      });
      if (failures.length > 0) {
        this.logIssues({ stepResults: failures, title: 'Failures' });
      }
      var warnings = featuresResult.stepResults.filter(function (stepResult) {
        return _lodash2.default.includes([_status2.default.PENDING, _status2.default.UNDEFINED], stepResult.status);
      });
      if (warnings.length > 0) {
        this.logIssues({ stepResults: warnings, title: 'Warnings' });
      }
      this.log((0, _helpers.formatSummary)({
        colorFns: this.colorFns,
        featuresResult: featuresResult
      }));
    }
  }, {
    key: 'logIssues',
    value: function logIssues(_ref) {
      var _this2 = this;

      var stepResults = _ref.stepResults,
          title = _ref.title;

      this.log(title + ':\n\n');
      stepResults.forEach(function (stepResult, index) {
        _this2.log((0, _helpers.formatIssue)({
          colorFns: _this2.colorFns,
          cwd: _this2.cwd,
          number: index + 1,
          snippetBuilder: _this2.snippetBuilder,
          stepResult: stepResult
        }));
      });
    }
  }]);
  return SummaryFormatter;
}(_3.default);

exports.default = SummaryFormatter;