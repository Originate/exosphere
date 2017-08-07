'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _keyword_type = require('../keyword_type');

var _step_arguments = require('./step_arguments');

var _step_arguments2 = _interopRequireDefault(_step_arguments);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Step = function Step(options) {
  (0, _classCallCheck3.default)(this, Step);
  var backgroundLines = options.backgroundLines,
      gherkinData = options.gherkinData,
      language = options.language,
      lineToKeywordMapping = options.lineToKeywordMapping,
      previousStep = options.previousStep,
      scenario = options.scenario;


  this.arguments = _lodash2.default.map(gherkinData.arguments, _step_arguments2.default.build);
  this.line = _lodash2.default.last(_lodash2.default.map(gherkinData.locations, 'line'));
  this.name = gherkinData.text;
  this.scenario = scenario;
  this.uri = scenario.uri;

  this.isBackground = _lodash2.default.some(gherkinData.locations, function (_ref) {
    var line = _ref.line;

    return _lodash2.default.includes(backgroundLines, line);
  });

  this.keyword = _lodash2.default.chain(gherkinData.locations).map(function (_ref2) {
    var line = _ref2.line;
    return lineToKeywordMapping[line];
  }).compact().first().value();

  this.keywordType = (0, _keyword_type.getStepKeywordType)({ language: language, previousStep: previousStep, step: this });
};

exports.default = Step;