'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.getStepKeywordType = getStepKeywordType;

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _gherkin = require('gherkin');

var _gherkin2 = _interopRequireDefault(_gherkin);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var types = {
  EVENT: 'event',
  OUTCOME: 'outcome',
  PRECONDITION: 'precondition'
};

exports.default = types;
function getStepKeywordType(_ref) {
  var language = _ref.language,
      previousStep = _ref.previousStep,
      step = _ref.step;

  var dialect = _gherkin2.default.DIALECTS[language];
  var type = _lodash2.default.find(['given', 'when', 'then', 'and', 'but'], function (keyword) {
    return _lodash2.default.includes(dialect[keyword], step.keyword);
  });
  switch (type) {
    case 'when':
      return types.EVENT;
    case 'then':
      return types.OUTCOME;
    case 'and':
    case 'but':
      if (previousStep) {
        return previousStep.keywordType;
      }
    // fallthrough
    default:
      return types.PRECONDITION;
  }
}