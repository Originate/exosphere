'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.formatIssue = formatIssue;

var _location_helpers = require('./location_helpers');

var _error_helpers = require('./error_helpers');

var _indentString = require('indent-string');

var _indentString2 = _interopRequireDefault(_indentString);

var _status = require('../../status');

var _status2 = _interopRequireDefault(_status);

var _cliTable = require('cli-table');

var _cliTable2 = _interopRequireDefault(_cliTable);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function formatIssue(_ref) {
  var colorFns = _ref.colorFns,
      cwd = _ref.cwd,
      number = _ref.number,
      snippetBuilder = _ref.snippetBuilder,
      stepResult = _ref.stepResult;

  var message = getStepResultMessage({ colorFns: colorFns, cwd: cwd, snippetBuilder: snippetBuilder, stepResult: stepResult });
  var prefix = number + ') ';
  var step = stepResult.step;
  var scenario = step.scenario;

  var text = prefix;

  if (scenario) {
    var scenarioLocation = (0, _location_helpers.formatLocation)(cwd, scenario);
    text += 'Scenario: ' + colorFns.bold(scenario.name) + ' - ' + colorFns.location(scenarioLocation);
  } else {
    text += 'Background:';
  }
  text += '\n';

  var stepText = 'Step: ' + colorFns.bold(step.keyword + (step.name || ''));
  if (step.uri) {
    var stepLocation = (0, _location_helpers.formatLocation)(cwd, step);
    stepText += ' - ' + colorFns.location(stepLocation);
  }
  text += (0, _indentString2.default)(stepText, prefix.length) + '\n';

  var stepDefinition = stepResult.stepDefinition;

  if (stepDefinition) {
    var stepDefinitionLocation = (0, _location_helpers.formatLocation)(cwd, stepDefinition);
    var stepDefinitionLine = 'Step Definition: ' + colorFns.location(stepDefinitionLocation);
    text += (0, _indentString2.default)(stepDefinitionLine, prefix.length) + '\n';
  }

  text += (0, _indentString2.default)('Message:', prefix.length) + '\n';
  text += (0, _indentString2.default)(message, prefix.length + 2) + '\n\n';
  return text;
}

function getAmbiguousStepResultMessage(_ref2) {
  var colorFns = _ref2.colorFns,
      cwd = _ref2.cwd,
      stepResult = _ref2.stepResult;
  var ambiguousStepDefinitions = stepResult.ambiguousStepDefinitions;

  var table = new _cliTable2.default({
    chars: {
      'bottom': '', 'bottom-left': '', 'bottom-mid': '', 'bottom-right': '',
      'left': '', 'left-mid': '',
      'mid': '', 'mid-mid': '', 'middle': ' - ',
      'right': '', 'right-mid': '',
      'top': '', 'top-left': '', 'top-mid': '', 'top-right': ''
    },
    style: {
      border: [], 'padding-left': 0, 'padding-right': 0
    }
  });
  table.push.apply(table, ambiguousStepDefinitions.map(function (stepDefinition) {
    var pattern = stepDefinition.pattern.toString();
    return [pattern, (0, _location_helpers.formatLocation)(cwd, stepDefinition)];
  }));
  var message = 'Multiple step definitions match:' + '\n' + (0, _indentString2.default)(table.toString(), 2);
  return colorFns.ambiguous(message);
}

function getFailedStepResultMessage(_ref3) {
  var colorFns = _ref3.colorFns,
      stepResult = _ref3.stepResult;
  var failureException = stepResult.failureException;

  return (0, _error_helpers.formatError)(failureException, colorFns);
}

function getPendingStepResultMessage(_ref4) {
  var colorFns = _ref4.colorFns;

  return colorFns.pending('Pending');
}

function getStepResultMessage(_ref5) {
  var colorFns = _ref5.colorFns,
      cwd = _ref5.cwd,
      snippetBuilder = _ref5.snippetBuilder,
      stepResult = _ref5.stepResult;

  switch (stepResult.status) {
    case _status2.default.AMBIGUOUS:
      return getAmbiguousStepResultMessage({ colorFns: colorFns, cwd: cwd, stepResult: stepResult });
    case _status2.default.FAILED:
      return getFailedStepResultMessage({ colorFns: colorFns, stepResult: stepResult });
    case _status2.default.UNDEFINED:
      return getUndefinedStepResultMessage({ colorFns: colorFns, snippetBuilder: snippetBuilder, stepResult: stepResult });
    case _status2.default.PENDING:
      return getPendingStepResultMessage({ colorFns: colorFns });
  }
}

function getUndefinedStepResultMessage(_ref6) {
  var colorFns = _ref6.colorFns,
      snippetBuilder = _ref6.snippetBuilder,
      stepResult = _ref6.stepResult;
  var step = stepResult.step;

  var snippet = snippetBuilder.build(step);
  var message = 'Undefined. Implement with the following snippet:' + '\n\n' + (0, _indentString2.default)(snippet, 2);
  return colorFns.undefined(message);
}