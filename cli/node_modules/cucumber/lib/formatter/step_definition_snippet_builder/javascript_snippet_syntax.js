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

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var JavaScriptSnippetSyntax = function () {
  function JavaScriptSnippetSyntax(snippetInterface) {
    (0, _classCallCheck3.default)(this, JavaScriptSnippetSyntax);

    this.snippetInterface = snippetInterface;
  }

  (0, _createClass3.default)(JavaScriptSnippetSyntax, [{
    key: 'build',
    value: function build(functionName, pattern, parameters, comment) {
      var functionKeyword = 'function ';
      if (this.snippetInterface === 'generator') {
        functionKeyword += '*';
      }

      var implementation = void 0;
      if (this.snippetInterface === 'callback') {
        var callbackName = _lodash2.default.last(parameters);
        implementation = callbackName + '(null, \'pending\');';
      } else {
        parameters.pop();
        implementation = 'return \'pending\';';
      }

      var snippet = functionName + '(\'' + pattern.replace(/'/g, '\\\'') + '\', ' + functionKeyword + '(' + parameters.join(', ') + ') {' + '\n' + '  // ' + comment + '\n' + '  ' + implementation + '\n' + '});';
      return snippet;
    }
  }]);
  return JavaScriptSnippetSyntax;
}();

exports.default = JavaScriptSnippetSyntax;