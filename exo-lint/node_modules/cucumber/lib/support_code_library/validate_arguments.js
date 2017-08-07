'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.default = validateArguments;

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var optionsValidation = {
  expectedType: 'object or function',
  predicate: function predicate(_ref) {
    var options = _ref.options;
    return _lodash2.default.isPlainObject(options);
  }
};

var optionsTimeoutValidation = {
  identifier: '"options.timeout"',
  expectedType: 'integer',
  predicate: function predicate(_ref2) {
    var options = _ref2.options;
    return !options.timeout || _lodash2.default.isInteger(options.timeout);
  }
};

var fnValidation = {
  expectedType: 'function',
  predicate: function predicate(_ref3) {
    var code = _ref3.code;
    return _lodash2.default.isFunction(code);
  }
};

var validations = {
  defineHook: [_lodash2.default.assign({ identifier: 'first argument' }, optionsValidation), {
    identifier: '"options.tags"',
    expectedType: 'string',
    predicate: function predicate(_ref4) {
      var options = _ref4.options;
      return !options.tags || _lodash2.default.isString(options.tags);
    }
  }, optionsTimeoutValidation, _lodash2.default.assign({ identifier: 'second argument' }, fnValidation)],
  defineStep: [{
    identifier: 'first argument',
    expectedType: 'string or regular expression',
    predicate: function predicate(_ref5) {
      var pattern = _ref5.pattern;
      return _lodash2.default.isRegExp(pattern) || _lodash2.default.isString(pattern);
    }
  }, _lodash2.default.assign({ identifier: 'second argument' }, optionsValidation), optionsTimeoutValidation, _lodash2.default.assign({ identifier: 'third argument' }, fnValidation)],
  registerHandler: [{
    identifier: 'first argument',
    expectedType: 'string',
    predicate: function predicate(_ref6) {
      var eventName = _ref6.eventName;
      return _lodash2.default.isString(eventName);
    }
  }, _lodash2.default.assign({ identifier: 'second argument' }, optionsValidation), optionsTimeoutValidation, _lodash2.default.assign({ identifier: 'third argument' }, fnValidation)]
};

function validateArguments(_ref7) {
  var args = _ref7.args,
      fnName = _ref7.fnName,
      relativeUri = _ref7.relativeUri;

  validations[fnName].forEach(function (_ref8) {
    var identifier = _ref8.identifier,
        expectedType = _ref8.expectedType,
        predicate = _ref8.predicate;

    if (!predicate(args)) {
      throw new Error(relativeUri + ': Invalid ' + identifier + ': should be a ' + expectedType);
    }
  });
}