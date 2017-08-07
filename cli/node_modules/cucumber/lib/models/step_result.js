'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _status = require('../status');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var StepResult = function StepResult(data) {
  (0, _classCallCheck3.default)(this, StepResult);

  _lodash2.default.assign(this, _lodash2.default.pick(data, ['ambiguousStepDefinitions', 'attachments', 'duration', 'failureException', 'step', 'stepDefinition', 'status']));
};

exports.default = StepResult;


(0, _status.addStatusPredicates)(StepResult.prototype);