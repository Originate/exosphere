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

var _get2 = require('babel-runtime/helpers/get');

var _get3 = _interopRequireDefault(_get2);

var _inherits2 = require('babel-runtime/helpers/inherits');

var _inherits3 = _interopRequireDefault(_inherits2);

var _defineProperty2 = require('babel-runtime/helpers/defineProperty');

var _defineProperty3 = _interopRequireDefault(_defineProperty2);

var _STATUS_CHARACTER_MAP;

var _hook = require('../models/hook');

var _hook2 = _interopRequireDefault(_hook);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

var _summary_formatter = require('./summary_formatter');

var _summary_formatter2 = _interopRequireDefault(_summary_formatter);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var STATUS_CHARACTER_MAPPING = (_STATUS_CHARACTER_MAP = {}, (0, _defineProperty3.default)(_STATUS_CHARACTER_MAP, _status2.default.AMBIGUOUS, 'A'), (0, _defineProperty3.default)(_STATUS_CHARACTER_MAP, _status2.default.FAILED, 'F'), (0, _defineProperty3.default)(_STATUS_CHARACTER_MAP, _status2.default.PASSED, '.'), (0, _defineProperty3.default)(_STATUS_CHARACTER_MAP, _status2.default.PENDING, 'P'), (0, _defineProperty3.default)(_STATUS_CHARACTER_MAP, _status2.default.SKIPPED, '-'), (0, _defineProperty3.default)(_STATUS_CHARACTER_MAP, _status2.default.UNDEFINED, 'U'), _STATUS_CHARACTER_MAP);

var ProgressFormatter = function (_SummaryFormatter) {
  (0, _inherits3.default)(ProgressFormatter, _SummaryFormatter);

  function ProgressFormatter() {
    (0, _classCallCheck3.default)(this, ProgressFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (ProgressFormatter.__proto__ || Object.getPrototypeOf(ProgressFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(ProgressFormatter, [{
    key: 'handleStepResult',
    value: function handleStepResult(stepResult) {
      var status = stepResult.status;
      if (!(stepResult.step instanceof _hook2.default && status === _status2.default.PASSED)) {
        var character = this.colorFns[status](STATUS_CHARACTER_MAPPING[status]);
        this.log(character);
      }
    }
  }, {
    key: 'handleFeaturesResult',
    value: function handleFeaturesResult(featuresResult) {
      this.log('\n\n');
      (0, _get3.default)(ProgressFormatter.prototype.__proto__ || Object.getPrototypeOf(ProgressFormatter.prototype), 'handleFeaturesResult', this).call(this, featuresResult);
    }
  }]);
  return ProgressFormatter;
}(_summary_formatter2.default);

exports.default = ProgressFormatter;