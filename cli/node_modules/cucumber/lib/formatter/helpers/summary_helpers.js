'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.formatSummary = formatSummary;

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _duration = require('duration');

var _duration2 = _interopRequireDefault(_duration);

var _hook = require('../../models/hook');

var _hook2 = _interopRequireDefault(_hook);

var _status = require('../../status');

var _status2 = _interopRequireDefault(_status);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var STATUS_REPORT_ORDER = [_status2.default.FAILED, _status2.default.AMBIGUOUS, _status2.default.UNDEFINED, _status2.default.PENDING, _status2.default.SKIPPED, _status2.default.PASSED];

function formatSummary(_ref) {
  var colorFns = _ref.colorFns,
      featuresResult = _ref.featuresResult;

  var scenarioSummary = getCountSummary({
    colorFns: colorFns,
    objects: featuresResult.scenarioResults,
    type: 'scenario'
  });
  var stepSummary = getCountSummary({
    colorFns: colorFns,
    objects: featuresResult.stepResults.filter(function (_ref2) {
      var step = _ref2.step;
      return !(step instanceof _hook2.default);
    }),
    type: 'step'
  });
  var durationSummary = getDuration(featuresResult);
  return [scenarioSummary, stepSummary, durationSummary].join('\n');
}

function getCountSummary(_ref3) {
  var colorFns = _ref3.colorFns,
      objects = _ref3.objects,
      type = _ref3.type;

  var counts = _lodash2.default.chain(objects).groupBy('status').mapValues('length').value();
  var total = _lodash2.default.reduce(counts, function (memo, value) {
    return memo + value;
  }) || 0;
  var text = total + ' ' + type + (total === 1 ? '' : 's');
  if (total > 0) {
    var details = [];
    STATUS_REPORT_ORDER.forEach(function (status) {
      if (counts[status] > 0) {
        details.push(colorFns[status](counts[status] + ' ' + status));
      }
    });
    text += ' (' + details.join(', ') + ')';
  }
  return text;
}

function getDuration(featuresResult) {
  var milliseconds = featuresResult.duration;
  var start = new Date(0);
  var end = new Date(milliseconds);
  var duration = new _duration2.default(start, end);

  return duration.minutes + 'm' + duration.toString('%S') + '.' + duration.toString('%L') + 's' + '\n';
}