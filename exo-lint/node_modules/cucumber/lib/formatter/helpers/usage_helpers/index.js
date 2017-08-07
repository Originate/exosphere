'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.getUsage = getUsage;

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _location_helpers = require('../location_helpers');

var _hook = require('../../../models/hook');

var _hook2 = _interopRequireDefault(_hook);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function buildEmptyMapping(stepDefinitions) {
  var mapping = {};
  stepDefinitions.forEach(function (stepDefinition) {
    var location = (0, _location_helpers.formatLocation)('', stepDefinition);
    mapping[location] = {
      line: stepDefinition.line,
      pattern: stepDefinition.pattern,
      matches: [],
      uri: stepDefinition.uri
    };
  });
  return mapping;
}

function buildMapping(_ref) {
  var stepDefinitions = _ref.stepDefinitions,
      stepResults = _ref.stepResults;

  var mapping = buildEmptyMapping(stepDefinitions);
  stepResults.forEach(function (stepResult) {
    var duration = stepResult.duration,
        step = stepResult.step,
        stepDefinition = stepResult.stepDefinition;

    if (!(step instanceof _hook2.default) && stepDefinition) {
      var location = (0, _location_helpers.formatLocation)('', stepDefinition);
      var match = {
        line: step.line,
        text: step.name,
        uri: step.uri
      };
      if (isFinite(duration)) {
        match.duration = duration;
      }
      if (mapping[location]) {
        mapping[location].matches.push(match);
      }
    }
  });
  return mapping;
}

function invertNumber(key) {
  return function (obj) {
    var value = obj[key];
    if (isFinite(value)) {
      return -1 * value;
    }
    return 1;
  };
}

function buildResult(mapping) {
  return _lodash2.default.chain(mapping).map(function (_ref2) {
    var line = _ref2.line,
        matches = _ref2.matches,
        pattern = _ref2.pattern,
        uri = _ref2.uri;

    var sortedMatches = _lodash2.default.sortBy(matches, [invertNumber('duration'), 'text']);
    var result = { line: line, matches: sortedMatches, pattern: pattern, uri: uri };
    var meanDuration = _lodash2.default.meanBy(matches, 'duration');
    if (isFinite(meanDuration)) {
      result.meanDuration = meanDuration;
    }
    return result;
  }).sortBy(invertNumber('meanDuration')).value();
}

function getUsage(_ref3) {
  var cwd = _ref3.cwd,
      stepDefinitions = _ref3.stepDefinitions,
      stepResults = _ref3.stepResults;

  var mapping = buildMapping({ cwd: cwd, stepDefinitions: stepDefinitions, stepResults: stepResults });
  return buildResult(mapping);
}