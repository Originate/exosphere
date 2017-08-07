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

var _data_table = require('../../models/step_arguments/data_table');

var _data_table2 = _interopRequireDefault(_data_table);

var _doc_string = require('../../models/step_arguments/doc_string');

var _doc_string2 = _interopRequireDefault(_doc_string);

var _keyword_type = require('../../keyword_type');

var _keyword_type2 = _interopRequireDefault(_keyword_type);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var StepDefinitionSnippetBuilder = function () {
  function StepDefinitionSnippetBuilder(_ref) {
    var snippetSyntax = _ref.snippetSyntax,
        parameterTypeRegistry = _ref.parameterTypeRegistry;
    (0, _classCallCheck3.default)(this, StepDefinitionSnippetBuilder);

    this.snippetSyntax = snippetSyntax;
    this.cucumberExpressionGenerator = new _cucumberExpressions.CucumberExpressionGenerator(parameterTypeRegistry);
  }

  (0, _createClass3.default)(StepDefinitionSnippetBuilder, [{
    key: 'build',
    value: function build(step) {
      var functionName = this.getFunctionName(step);
      var generatedExpression = this.cucumberExpressionGenerator.generateExpression(step.name, true);
      var pattern = generatedExpression.source;
      var parameters = this.getParameters(step, generatedExpression.parameterNames);
      var comment = 'Write code here that turns the phrase above into concrete actions';
      return this.snippetSyntax.build(functionName, pattern, parameters, comment);
    }
  }, {
    key: 'getFunctionName',
    value: function getFunctionName(step) {
      switch (step.keywordType) {
        case _keyword_type2.default.EVENT:
          return 'When';
        case _keyword_type2.default.OUTCOME:
          return 'Then';
        case _keyword_type2.default.PRECONDITION:
          return 'Given';
      }
    }
  }, {
    key: 'getParameters',
    value: function getParameters(step, expressionParameterNames) {
      return _lodash2.default.concat(expressionParameterNames, this.getStepArgumentParameters(step), 'callback');
    }
  }, {
    key: 'getStepArgumentParameters',
    value: function getStepArgumentParameters(step) {
      return step.arguments.map(function (arg) {
        if (arg instanceof _data_table2.default) {
          return 'table';
        } else if (arg instanceof _doc_string2.default) {
          return 'string';
        } else {
          throw new Error('Unknown argument type: ' + arg);
        }
      });
    }
  }]);
  return StepDefinitionSnippetBuilder;
}();

exports.default = StepDefinitionSnippetBuilder;