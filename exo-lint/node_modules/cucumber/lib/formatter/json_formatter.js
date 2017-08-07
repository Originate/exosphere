'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _possibleConstructorReturn2 = require('babel-runtime/helpers/possibleConstructorReturn');

var _possibleConstructorReturn3 = _interopRequireDefault(_possibleConstructorReturn2);

var _inherits2 = require('babel-runtime/helpers/inherits');

var _inherits3 = _interopRequireDefault(_inherits2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _data_table = require('../models/step_arguments/data_table');

var _data_table2 = _interopRequireDefault(_data_table);

var _doc_string = require('../models/step_arguments/doc_string');

var _doc_string2 = _interopRequireDefault(_doc_string);

var _2 = require('./');

var _3 = _interopRequireDefault(_2);

var _status = require('../status');

var _status2 = _interopRequireDefault(_status);

var _util = require('util');

var _util2 = _interopRequireDefault(_util);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var JsonFormatter = function (_Formatter) {
  (0, _inherits3.default)(JsonFormatter, _Formatter);

  function JsonFormatter(options) {
    (0, _classCallCheck3.default)(this, JsonFormatter);

    var _this = (0, _possibleConstructorReturn3.default)(this, (JsonFormatter.__proto__ || Object.getPrototypeOf(JsonFormatter)).call(this, options));

    _this.features = [];
    return _this;
  }

  (0, _createClass3.default)(JsonFormatter, [{
    key: 'convertNameToId',
    value: function convertNameToId(obj) {
      return obj.name.replace(/ /g, '-').toLowerCase();
    }
  }, {
    key: 'formatAttachments',
    value: function formatAttachments(attachments) {
      return attachments.map(function (attachment) {
        return {
          data: attachment.data,
          mime_type: attachment.mimeType
        };
      });
    }
  }, {
    key: 'formatDataTable',
    value: function formatDataTable(dataTable) {
      return {
        rows: dataTable.raw().map(function (row) {
          return { cells: row };
        })
      };
    }
  }, {
    key: 'formatDocString',
    value: function formatDocString(docString) {
      return _lodash2.default.pick(docString, ['content', 'contentType', 'line']);
    }
  }, {
    key: 'formatStepArguments',
    value: function formatStepArguments(stepArguments) {
      var _this2 = this;

      return _lodash2.default.map(stepArguments, function (arg) {
        if (arg instanceof _data_table2.default) {
          return _this2.formatDataTable(arg);
        } else if (arg instanceof _doc_string2.default) {
          return _this2.formatDocString(arg);
        } else {
          throw new Error('Unknown argument type: ' + _util2.default.inspect(arg));
        }
      });
    }
  }, {
    key: 'handleAfterFeatures',
    value: function handleAfterFeatures() {
      this.log(JSON.stringify(this.features, null, 2));
    }
  }, {
    key: 'handleBeforeFeature',
    value: function handleBeforeFeature(feature) {
      this.currentFeature = _lodash2.default.pick(feature, ['description', 'keyword', 'line', 'name', 'tags', 'uri']);
      _lodash2.default.assign(this.currentFeature, {
        elements: [],
        id: this.convertNameToId(feature)
      });
      this.features.push(this.currentFeature);
    }
  }, {
    key: 'handleBeforeScenario',
    value: function handleBeforeScenario(scenario) {
      this.currentScenario = _lodash2.default.pick(scenario, ['description', 'keyword', 'line', 'name', 'tags']);
      _lodash2.default.assign(this.currentScenario, {
        id: this.currentFeature.id + ';' + this.convertNameToId(scenario),
        steps: []
      });
      this.currentFeature.elements.push(this.currentScenario);
    }
  }, {
    key: 'handleStepResult',
    value: function handleStepResult(stepResult) {
      var step = stepResult.step;
      var status = stepResult.status;

      var currentStep = {
        arguments: this.formatStepArguments(step.arguments),
        keyword: step.keyword,
        name: step.name,
        result: { status: status }
      };

      if (step.isBackground) {
        currentStep.isBackground = true;
      }

      if (step.constructor.name === 'Hook') {
        currentStep.hidden = true;
      } else {
        currentStep.line = step.line;
      }

      if (status === _status2.default.PASSED || status === _status2.default.FAILED) {
        currentStep.result.duration = stepResult.duration;
      }

      if (_lodash2.default.size(stepResult.attachments) > 0) {
        currentStep.embeddings = this.formatAttachments(stepResult.attachments);
      }

      if (status === _status2.default.FAILED && stepResult.failureException) {
        currentStep.result.error_message = stepResult.failureException.stack || stepResult.failureException;
      }

      if (stepResult.stepDefinition) {
        var location = stepResult.stepDefinition.uri + ':' + stepResult.stepDefinition.line;
        currentStep.match = { location: location };
      }

      this.currentScenario.steps.push(currentStep);
    }
  }]);
  return JsonFormatter;
}(_3.default);

exports.default = JsonFormatter;