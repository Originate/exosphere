'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _gherkin = require('gherkin');

var _gherkin2 = _interopRequireDefault(_gherkin);

var _step = require('./step');

var _step2 = _interopRequireDefault(_step);

var _tag = require('./tag');

var _tag2 = _interopRequireDefault(_tag);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Scenario = function Scenario(options) {
  var _this = this;

  (0, _classCallCheck3.default)(this, Scenario);
  var backgroundStepLines = options.backgroundStepLines,
      feature = options.feature,
      gherkinData = options.gherkinData,
      language = options.language,
      lineToDescriptionMapping = options.lineToDescriptionMapping,
      stepLineToKeywordMapping = options.stepLineToKeywordMapping;


  this.feature = feature;
  this.keyword = _lodash2.default.first(_gherkin2.default.DIALECTS[language].scenario);
  this.lines = _lodash2.default.map(gherkinData.locations, 'line');
  this.name = gherkinData.name;
  this.tags = _lodash2.default.map(gherkinData.tags, _tag2.default.build);
  this.uri = feature.uri;

  this.line = _lodash2.default.first(this.lines);
  this.description = _lodash2.default.chain(this.lines).map(function (line) {
    return lineToDescriptionMapping[line];
  }).compact().first().value();

  var previousStep = void 0;
  this.steps = _lodash2.default.map(gherkinData.steps, function (gherkinStepData) {
    var step = new _step2.default({
      backgroundLines: backgroundStepLines,
      gherkinData: gherkinStepData,
      language: language,
      lineToKeywordMapping: stepLineToKeywordMapping,
      previousStep: previousStep,
      scenario: _this
    });
    previousStep = step;
    return step;
  });
};

exports.default = Scenario;