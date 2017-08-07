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

var _cucumberExpressions = require('cucumber-expressions');

var _data_table = require('./step_arguments/data_table');

var _data_table2 = _interopRequireDefault(_data_table);

var _doc_string = require('./step_arguments/doc_string');

var _doc_string2 = _interopRequireDefault(_doc_string);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var StepDefinition = function () {
  function StepDefinition(_ref) {
    var code = _ref.code,
        line = _ref.line,
        options = _ref.options,
        pattern = _ref.pattern,
        uri = _ref.uri;
    (0, _classCallCheck3.default)(this, StepDefinition);

    this.code = code;
    this.line = line;
    this.options = options;
    this.pattern = pattern;
    this.uri = uri;
  }

  (0, _createClass3.default)(StepDefinition, [{
    key: 'buildInvalidCodeLengthMessage',
    value: function buildInvalidCodeLengthMessage(syncOrPromiseLength, callbackLength) {
      return 'function has ' + this.code.length + ' arguments' + ', should have ' + syncOrPromiseLength + ' (if synchronous or returning a promise)' + ' or ' + callbackLength + ' (if accepting a callback)';
    }
  }, {
    key: 'getInvalidCodeLengthMessage',
    value: function getInvalidCodeLengthMessage(parameters) {
      return this.buildInvalidCodeLengthMessage(parameters.length, parameters.length + 1);
    }
  }, {
    key: 'getInvocationParameters',
    value: function getInvocationParameters(_ref2) {
      var step = _ref2.step,
          parameterTypeRegistry = _ref2.parameterTypeRegistry;

      var cucumberExpression = this.getCucumberExpression(parameterTypeRegistry);
      var stepNameParameters = _lodash2.default.map(cucumberExpression.match(step.name), 'transformedValue');
      var stepArgumentParameters = step.arguments.map(function (arg) {
        if (arg instanceof _data_table2.default) {
          return arg;
        } else if (arg instanceof _doc_string2.default) {
          return arg.content;
        } else {
          throw new Error('Unknown argument type:' + arg);
        }
      });
      return stepNameParameters.concat(stepArgumentParameters);
    }
  }, {
    key: 'getCucumberExpression',
    value: function getCucumberExpression(parameterTypeRegistry) {
      if (typeof this.pattern === 'string') {
        return new _cucumberExpressions.CucumberExpression(this.pattern, [], parameterTypeRegistry);
      } else {
        return new _cucumberExpressions.RegularExpression(this.pattern, [], parameterTypeRegistry);
      }
    }
  }, {
    key: 'getValidCodeLengths',
    value: function getValidCodeLengths(parameters) {
      return [parameters.length, parameters.length + 1];
    }
  }, {
    key: 'matchesStepName',
    value: function matchesStepName(_ref3) {
      var stepName = _ref3.stepName,
          parameterTypeRegistry = _ref3.parameterTypeRegistry;

      var cucumberExpression = this.getCucumberExpression(parameterTypeRegistry);
      return Boolean(cucumberExpression.match(stepName));
    }
  }]);
  return StepDefinition;
}();

exports.default = StepDefinition;