'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _cucumberExpressions = require('cucumber-expressions');

function build() {
  var parameterTypeRegistry = new _cucumberExpressions.ParameterTypeRegistry();
  var stringInDoubleQuotesParameterType = new _cucumberExpressions.ParameterType('stringInDoubleQuotes', null, /"[^"]+"/, JSON.parse);
  parameterTypeRegistry.defineParameterType(stringInDoubleQuotesParameterType);
  return parameterTypeRegistry;
}

exports.default = { build: build };