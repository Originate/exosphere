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

var _cliTable = require('cli-table');

var _cliTable2 = _interopRequireDefault(_cliTable);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var UsageFormatter = function (_Formatter) {
  (0, _inherits3.default)(UsageFormatter, _Formatter);

  function UsageFormatter() {
    (0, _classCallCheck3.default)(this, UsageFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (UsageFormatter.__proto__ || Object.getPrototypeOf(UsageFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(UsageFormatter, [{
    key: 'handleFeaturesResult',
    value: function handleFeaturesResult(featuresResult) {
      var _this2 = this;

      var usage = (0, _helpers.getUsage)({
        stepDefinitions: this.supportCodeLibrary.stepDefinitions,
        stepResults: featuresResult.stepResults
      });
      if (usage.length === 0) {
        this.log('No step definitions');
        return;
      }
      var table = new _cliTable2.default({
        head: ['Pattern / Text', 'Duration', 'Location'],
        style: {
          border: [],
          head: []
        }
      });
      usage.forEach(function (_ref) {
        var line = _ref.line,
            matches = _ref.matches,
            meanDuration = _ref.meanDuration,
            pattern = _ref.pattern,
            uri = _ref.uri;

        var col1 = [pattern.toString()];
        var col2 = [];
        if (matches.length > 0) {
          if (isFinite(meanDuration)) {
            col2.push(parseFloat(meanDuration.toFixed(2)) + 'ms');
          } else {
            col2.push('-');
          }
        } else {
          col2.push('UNUSED');
        }
        var col3 = [(0, _helpers.formatLocation)(_this2.cwd, { line: line, uri: uri })];
        _lodash2.default.take(matches, 5).forEach(function (match) {
          col1.push('  ' + match.text);
          if (isFinite(match.duration)) {
            col2.push(match.duration + 'ms');
          } else {
            col2.push('-');
          }
          col3.push((0, _helpers.formatLocation)(_this2.cwd, match));
        });
        if (matches.length > 5) {
          col1.push('  ' + (matches.length - 5) + ' more');
        }
        table.push([col1.join('\n'), col2.join('\n'), col3.join('\n')]);
      });
      this.log(table.toString() + '\n');
    }
  }]);
  return UsageFormatter;
}(_3.default);

exports.default = UsageFormatter;