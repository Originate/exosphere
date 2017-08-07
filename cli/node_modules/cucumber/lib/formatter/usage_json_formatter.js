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

var _helpers = require('./helpers');

var _ = require('./');

var _2 = _interopRequireDefault(_);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var UsageJsonFormatter = function (_Formatter) {
  (0, _inherits3.default)(UsageJsonFormatter, _Formatter);

  function UsageJsonFormatter() {
    (0, _classCallCheck3.default)(this, UsageJsonFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (UsageJsonFormatter.__proto__ || Object.getPrototypeOf(UsageJsonFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(UsageJsonFormatter, [{
    key: 'handleFeaturesResult',
    value: function handleFeaturesResult(featuresResult) {
      var usage = (0, _helpers.getUsage)({
        cwd: this.cwd,
        stepDefinitions: this.supportCodeLibrary.stepDefinitions,
        stepResults: featuresResult.stepResults
      });
      this.log(JSON.stringify(usage, null, 2));
    }
  }]);
  return UsageJsonFormatter;
}(_2.default);

exports.default = UsageJsonFormatter;