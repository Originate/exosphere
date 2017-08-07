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

var _scenario_filter = require('../scenario_filter');

var _scenario_filter2 = _interopRequireDefault(_scenario_filter);

var _step_definition = require('./step_definition');

var _step_definition2 = _interopRequireDefault(_step_definition);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var HookDefinition = function (_StepDefinition) {
  (0, _inherits3.default)(HookDefinition, _StepDefinition);

  function HookDefinition(data) {
    (0, _classCallCheck3.default)(this, HookDefinition);

    var _this = (0, _possibleConstructorReturn3.default)(this, (HookDefinition.__proto__ || Object.getPrototypeOf(HookDefinition)).call(this, data));

    _this.scenarioFilter = new _scenario_filter2.default({ tagExpression: _this.options.tags });
    return _this;
  }

  (0, _createClass3.default)(HookDefinition, [{
    key: 'appliesToScenario',
    value: function appliesToScenario(scenario) {
      return this.scenarioFilter.matches(scenario);
    }
  }, {
    key: 'getInvalidCodeLengthMessage',
    value: function getInvalidCodeLengthMessage() {
      return this.buildInvalidCodeLengthMessage('0 or 1', '2');
    }
  }, {
    key: 'getInvocationParameters',
    value: function getInvocationParameters(_ref) {
      var scenarioResult = _ref.scenarioResult;

      return [scenarioResult];
    }
  }, {
    key: 'getValidCodeLengths',
    value: function getValidCodeLengths() {
      return [0, 1, 2];
    }
  }]);
  return HookDefinition;
}(_step_definition2.default);

exports.default = HookDefinition;