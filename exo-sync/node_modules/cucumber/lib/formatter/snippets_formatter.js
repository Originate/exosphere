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

var _ = require('./');

var _2 = _interopRequireDefault(_);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var SnippetsFormatter = function (_Formatter) {
  (0, _inherits3.default)(SnippetsFormatter, _Formatter);

  function SnippetsFormatter() {
    (0, _classCallCheck3.default)(this, SnippetsFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (SnippetsFormatter.__proto__ || Object.getPrototypeOf(SnippetsFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(SnippetsFormatter, [{
    key: 'handleStepResult',
    value: function handleStepResult(stepResult) {
      if (stepResult.status === _status2.default.UNDEFINED) {
        var snippet = this.snippetBuilder.build(stepResult.step);
        this.log(snippet + '\n\n');
      }
    }
  }]);
  return SnippetsFormatter;
}(_2.default);

exports.default = SnippetsFormatter;