'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _defineProperty2 = require('babel-runtime/helpers/defineProperty');

var _defineProperty3 = _interopRequireDefault(_defineProperty2);

exports.defineHook = defineHook;
exports.defineStep = defineStep;
exports.registerHandler = registerHandler;
exports.addTransform = addTransform;
exports.defineParameterType = defineParameterType;

var _util = require('util');

var _util2 = _interopRequireDefault(_util);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _cucumberExpressions = require('cucumber-expressions');

var _helpers = require('../formatter/helpers');

var _hook_definition = require('../models/hook_definition');

var _hook_definition2 = _interopRequireDefault(_hook_definition);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _stacktraceJs = require('stacktrace-js');

var _stacktraceJs2 = _interopRequireDefault(_stacktraceJs);

var _step_definition = require('../models/step_definition');

var _step_definition2 = _interopRequireDefault(_step_definition);

var _validate_arguments = require('./validate_arguments');

var _validate_arguments2 = _interopRequireDefault(_validate_arguments);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function defineHook(cwd, collection) {
  return function (options, code) {
    if (typeof options === 'string') {
      options = { tags: options };
    } else if (typeof options === 'function') {
      code = options;
      options = {};
    }

    var _getDefinitionLineAnd = getDefinitionLineAndUri(),
        line = _getDefinitionLineAnd.line,
        uri = _getDefinitionLineAnd.uri;

    (0, _validate_arguments2.default)({
      args: { code: code, options: options },
      fnName: 'defineHook',
      relativeUri: (0, _helpers.formatLocation)(cwd, { line: line, uri: uri })
    });
    var hookDefinition = new _hook_definition2.default({ code: code, line: line, options: options, uri: uri });
    collection.push(hookDefinition);
  };
}

function defineStep(cwd, collection) {
  return function (pattern, options, code) {
    if (typeof options === 'function') {
      code = options;
      options = {};
    }

    var _getDefinitionLineAnd2 = getDefinitionLineAndUri(),
        line = _getDefinitionLineAnd2.line,
        uri = _getDefinitionLineAnd2.uri;

    (0, _validate_arguments2.default)({
      args: { code: code, pattern: pattern, options: options },
      fnName: 'defineStep',
      relativeUri: (0, _helpers.formatLocation)(cwd, { line: line, uri: uri })
    });
    var stepDefinition = new _step_definition2.default({ code: code, line: line, options: options, pattern: pattern, uri: uri });
    collection.push(stepDefinition);
  };
}

function getDefinitionLineAndUri() {
  var stackframes = _stacktraceJs2.default.getSync();
  var stackframe = stackframes.length > 2 ? stackframes[2] : stackframes[0];
  var line = stackframe.getLineNumber();
  var fileName = stackframe.getFileName();
  var uri = fileName ? fileName.replace(/\//g, _path2.default.sep) : 'unknown';
  return { line: line, uri: uri };
}

function registerHandler(cwd, collection) {
  return function (eventName, options, code) {
    var _$assign;

    if (typeof options === 'function') {
      code = options;
      options = {};
    }

    var _getDefinitionLineAnd3 = getDefinitionLineAndUri(),
        line = _getDefinitionLineAnd3.line,
        uri = _getDefinitionLineAnd3.uri;

    (0, _validate_arguments2.default)({
      args: { code: code, eventName: eventName, options: options },
      fnName: 'registerHandler',
      relativeUri: (0, _helpers.formatLocation)(cwd, { line: line, uri: uri })
    });
    var listener = _lodash2.default.assign((_$assign = {}, (0, _defineProperty3.default)(_$assign, 'handle' + eventName, code), (0, _defineProperty3.default)(_$assign, 'relativeUri', (0, _helpers.formatLocation)(cwd, { line: line, uri: uri })), _$assign), options);
    collection.push(listener);
  };
}

function addTransform(parameterTypeRegistry) {
  return _util2.default.deprecate(function (_ref) {
    var captureGroupRegexps = _ref.captureGroupRegexps,
        transformer = _ref.transformer,
        typeName = _ref.typeName;

    var parameter = new _cucumberExpressions.ParameterType(typeName, null, captureGroupRegexps, transformer);
    parameterTypeRegistry.defineParameterType(parameter);
  }, 'addTransform is deprecated and will be removed in a future version. Please use defineParameterType instead.');
}

function defineParameterType(parameterTypeRegistry) {
  return function (_ref2) {
    var regexp = _ref2.regexp,
        transformer = _ref2.transformer,
        typeName = _ref2.typeName;

    var parameter = new _cucumberExpressions.ParameterType(typeName, null, regexp, transformer);
    parameterTypeRegistry.defineParameterType(parameter);
  };
}