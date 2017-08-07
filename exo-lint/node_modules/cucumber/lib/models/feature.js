'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _scenario = require('./scenario');

var _scenario2 = _interopRequireDefault(_scenario);

var _tag = require('./tag');

var _tag2 = _interopRequireDefault(_tag);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Feature = function Feature(_ref) {
  var _this = this;

  var gherkinData = _ref.gherkinData,
      gherkinPickles = _ref.gherkinPickles,
      scenarioFilter = _ref.scenarioFilter,
      uri = _ref.uri;
  (0, _classCallCheck3.default)(this, Feature);

  this.description = gherkinData.description;
  this.keyword = gherkinData.keyword;
  this.line = gherkinData.location.line;
  this.name = gherkinData.name;
  this.tags = _lodash2.default.map(gherkinData.tags, _tag2.default.build);
  this.uri = uri;

  var backgroundStepLines = _lodash2.default.chain(gherkinData.children).filter(['type', 'Background']).map('steps').flatten().map(function (step) {
    return step.location.line;
  }).value();

  var scenarioLineToDescriptionMapping = _lodash2.default.chain(gherkinData.children).map(function (element) {
    return [element.location.line, element.description];
  }).fromPairs().value();

  var stepLineToKeywordMapping = _lodash2.default.chain(gherkinData.children).map('steps').flatten().map(function (step) {
    return [step.location.line, step.keyword];
  }).fromPairs().value();

  this.scenarios = _lodash2.default.chain(gherkinPickles).map(function (gherkinPickle) {
    return new _scenario2.default({
      backgroundStepLines: backgroundStepLines,
      feature: _this,
      gherkinData: gherkinPickle,
      language: gherkinData.language,
      lineToDescriptionMapping: scenarioLineToDescriptionMapping,
      stepLineToKeywordMapping: stepLineToKeywordMapping
    });
  }).filter(function (scenario) {
    return scenarioFilter.matches(scenario);
  }).value();
};

exports.default = Feature;