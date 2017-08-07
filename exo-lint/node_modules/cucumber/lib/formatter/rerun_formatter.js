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

var _2 = require('./');

var _3 = _interopRequireDefault(_2);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var RerunFormatter = function (_Formatter) {
  (0, _inherits3.default)(RerunFormatter, _Formatter);

  function RerunFormatter() {
    (0, _classCallCheck3.default)(this, RerunFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (RerunFormatter.__proto__ || Object.getPrototypeOf(RerunFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(RerunFormatter, [{
    key: 'handleFeaturesResult',
    value: function handleFeaturesResult(featuresResult) {
      var _this2 = this;

      var mapping = {};
      featuresResult.scenarioResults.forEach(function (scenarioResult) {
        if (scenarioResult.status !== _status2.default.PASSED) {
          var scenario = scenarioResult.scenario;
          var relativeUri = _path2.default.relative(_this2.cwd, scenario.uri);
          if (!mapping[relativeUri]) {
            mapping[relativeUri] = [];
          }
          mapping[relativeUri].push(scenario.line);
        }
      });
      var text = _lodash2.default.map(mapping, function (lines, relativeUri) {
        return relativeUri + ':' + lines.join(':');
      }).join('\n');
      this.log(text);
    }
  }]);
  return RerunFormatter;
}(_3.default);

exports.default = RerunFormatter;