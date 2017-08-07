'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.addStatusPredicates = addStatusPredicates;
exports.getStatusMapping = getStatusMapping;

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _upperCaseFirst = require('upper-case-first');

var _upperCaseFirst2 = _interopRequireDefault(_upperCaseFirst);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var statuses = {
  AMBIGUOUS: 'ambiguous',
  FAILED: 'failed',
  PASSED: 'passed',
  PENDING: 'pending',
  SKIPPED: 'skipped',
  UNDEFINED: 'undefined'
};

exports.default = statuses;
function addStatusPredicates(protoype) {
  _lodash2.default.each(statuses, function (status) {
    protoype['is' + (0, _upperCaseFirst2.default)(status)] = function () {
      return this.status === status;
    };
  });
}

function getStatusMapping(initialValue) {
  return _lodash2.default.chain(statuses).map(function (status) {
    return [status, initialValue];
  }).fromPairs().value();
}