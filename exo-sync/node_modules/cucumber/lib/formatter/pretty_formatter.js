'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _defineProperty2 = require('babel-runtime/helpers/defineProperty');

var _defineProperty3 = _interopRequireDefault(_defineProperty2);

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _possibleConstructorReturn2 = require('babel-runtime/helpers/possibleConstructorReturn');

var _possibleConstructorReturn3 = _interopRequireDefault(_possibleConstructorReturn2);

var _inherits2 = require('babel-runtime/helpers/inherits');

var _inherits3 = _interopRequireDefault(_inherits2);

var _PrettyFormatter$CHAR;

var _data_table = require('../models/step_arguments/data_table');

var _data_table2 = _interopRequireDefault(_data_table);

var _doc_string = require('../models/step_arguments/doc_string');

var _doc_string2 = _interopRequireDefault(_doc_string);

var _figures = require('figures');

var _figures2 = _interopRequireDefault(_figures);

var _hook = require('../models/hook');

var _hook2 = _interopRequireDefault(_hook);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

var _summary_formatter = require('./summary_formatter');

var _summary_formatter2 = _interopRequireDefault(_summary_formatter);

var _cliTable = require('cli-table');

var _cliTable2 = _interopRequireDefault(_cliTable);

var _indentString = require('indent-string');

var _indentString2 = _interopRequireDefault(_indentString);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var PrettyFormatter = function (_SummaryFormatter) {
  (0, _inherits3.default)(PrettyFormatter, _SummaryFormatter);

  function PrettyFormatter() {
    (0, _classCallCheck3.default)(this, PrettyFormatter);
    return (0, _possibleConstructorReturn3.default)(this, (PrettyFormatter.__proto__ || Object.getPrototypeOf(PrettyFormatter)).apply(this, arguments));
  }

  (0, _createClass3.default)(PrettyFormatter, [{
    key: 'applyColor',
    value: function applyColor(stepResult, text) {
      var status = stepResult.status;
      return this.colorFns[status](text);
    }
  }, {
    key: 'formatDataTable',
    value: function formatDataTable(dataTable) {
      var rows = dataTable.raw().map(function (row) {
        return row.map(function (cell) {
          return cell.replace(/\\/g, '\\\\').replace(/\n/g, '\\n');
        });
      });
      var table = new _cliTable2.default({
        chars: {
          'bottom': '', 'bottom-left': '', 'bottom-mid': '', 'bottom-right': '',
          'left': '|', 'left-mid': '',
          'mid': '', 'mid-mid': '', 'middle': '|',
          'right': '|', 'right-mid': '',
          'top': '', 'top-left': '', 'top-mid': '', 'top-right': ''
        },
        style: {
          border: [], 'padding-left': 1, 'padding-right': 1
        }
      });
      table.push.apply(table, rows);
      return table.toString();
    }
  }, {
    key: 'formatDocString',
    value: function formatDocString(docString) {
      return '"""\n' + docString.content + '\n"""';
    }
  }, {
    key: 'formatTags',
    value: function formatTags(tags) {
      if (tags.length === 0) {
        return '';
      }
      var tagNames = tags.map(function (tag) {
        return tag.name;
      });
      return this.colorFns.tag(tagNames.join(' '));
    }
  }, {
    key: 'handleAfterScenario',
    value: function handleAfterScenario() {
      this.log('\n');
    }
  }, {
    key: 'handleBeforeFeature',
    value: function handleBeforeFeature(feature) {
      var text = '';
      var tagsText = this.formatTags(feature.tags);
      if (tagsText) {
        text = tagsText + '\n';
      }
      text += feature.keyword + ': ' + feature.name;
      var description = feature.description;
      if (description) {
        text += '\n\n' + (0, _indentString2.default)(description, 2);
      }
      this.log(text + '\n\n');
    }
  }, {
    key: 'handleBeforeScenario',
    value: function handleBeforeScenario(scenario) {
      var text = '';
      var tagsText = this.formatTags(scenario.tags);
      if (tagsText) {
        text = tagsText + '\n';
      }
      text += scenario.keyword + ': ' + scenario.name;
      this.logIndented(text + '\n', 1);
    }
  }, {
    key: 'handleStepResult',
    value: function handleStepResult(stepResult) {
      if (!(stepResult.step instanceof _hook2.default)) {
        this.logStepResult(stepResult);
      }
    }
  }, {
    key: 'logIndented',
    value: function logIndented(text, level) {
      this.log((0, _indentString2.default)(text, level * 2));
    }
  }, {
    key: 'logStepResult',
    value: function logStepResult(stepResult) {
      var _this2 = this;

      var status = stepResult.status,
          step = stepResult.step;

      var colorFn = this.colorFns[status];

      var symbol = PrettyFormatter.CHARACTERS[stepResult.status];
      var identifier = colorFn(symbol + ' ' + step.keyword + (step.name || ''));
      this.logIndented(identifier + '\n', 1);

      step.arguments.forEach(function (arg) {
        var str = void 0;
        if (arg instanceof _data_table2.default) {
          str = _this2.formatDataTable(arg);
        } else if (arg instanceof _doc_string2.default) {
          str = _this2.formatDocString(arg);
        } else {
          throw new Error('Unknown argument type: ' + arg);
        }
        _this2.logIndented(colorFn(str) + '\n', 3);
      });
    }
  }]);
  return PrettyFormatter;
}(_summary_formatter2.default);

exports.default = PrettyFormatter;


PrettyFormatter.CHARACTERS = (_PrettyFormatter$CHAR = {}, (0, _defineProperty3.default)(_PrettyFormatter$CHAR, _status2.default.AMBIGUOUS, _figures2.default.cross), (0, _defineProperty3.default)(_PrettyFormatter$CHAR, _status2.default.FAILED, _figures2.default.cross), (0, _defineProperty3.default)(_PrettyFormatter$CHAR, _status2.default.PASSED, _figures2.default.tick), (0, _defineProperty3.default)(_PrettyFormatter$CHAR, _status2.default.PENDING, '?'), (0, _defineProperty3.default)(_PrettyFormatter$CHAR, _status2.default.SKIPPED, '-'), (0, _defineProperty3.default)(_PrettyFormatter$CHAR, _status2.default.UNDEFINED, '?'), _PrettyFormatter$CHAR);