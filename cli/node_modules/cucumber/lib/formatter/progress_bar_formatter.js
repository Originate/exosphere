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

var _hook = require('../models/hook');

var _hook2 = _interopRequireDefault(_hook);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

var _2 = require('./');

var _3 = _interopRequireDefault(_2);

var _progress = require('progress');

var _progress2 = _interopRequireDefault(_progress);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var statusToReport = [_status2.default.AMBIGUOUS, _status2.default.FAILED, _status2.default.PENDING, _status2.default.UNDEFINED];

// Inspired by https://github.com/thekompanee/fuubar and https://github.com/martinciu/fuubar-cucumber

var ProgressBarFormatter = function (_Formatter) {
  (0, _inherits3.default)(ProgressBarFormatter, _Formatter);

  function ProgressBarFormatter(options) {
    (0, _classCallCheck3.default)(this, ProgressBarFormatter);

    var _this = (0, _possibleConstructorReturn3.default)(this, (ProgressBarFormatter.__proto__ || Object.getPrototypeOf(ProgressBarFormatter)).call(this, options));

    _this.issueCount = 0;
    return _this;
  }

  (0, _createClass3.default)(ProgressBarFormatter, [{
    key: 'handleBeforeFeatures',
    value: function handleBeforeFeatures(features) {
      var numberOfSteps = _lodash2.default.sumBy(features, function (feature) {
        return _lodash2.default.sumBy(feature.scenarios, function (scenario) {
          return scenario.steps.length;
        });
      });
      this.progressBar = new _progress2.default(':current/:total steps [:bar] ', {
        clear: true,
        incomplete: ' ',
        stream: this.stream,
        total: numberOfSteps,
        width: this.stream.columns || 80
      });
    }
  }, {
    key: 'handleStepResult',
    value: function handleStepResult(stepResult) {
      if (!(stepResult.step instanceof _hook2.default)) {
        this.progressBar.tick();
      }
      if (_lodash2.default.includes(statusToReport, stepResult.status)) {
        this.issueCount += 1;
        this.progressBar.interrupt((0, _helpers.formatIssue)({
          colorFns: this.colorFns,
          cwd: this.cwd,
          number: this.issueCount,
          snippetBuilder: this.snippetBuilder,
          stepResult: stepResult
        }));
      }
    }
  }, {
    key: 'handleFeaturesResult',
    value: function handleFeaturesResult(featuresResult) {
      this.log((0, _helpers.formatSummary)({
        colorFns: this.colorFns,
        featuresResult: featuresResult
      }));
    }
  }]);
  return ProgressBarFormatter;
}(_3.default);

exports.default = ProgressBarFormatter;