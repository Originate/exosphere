'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Event = function () {
  function Event(_ref) {
    var data = _ref.data,
        name = _ref.name;
    (0, _classCallCheck3.default)(this, Event);

    this.data = data;
    this.name = name;
  }

  (0, _createClass3.default)(Event, [{
    key: 'buildBeforeEvent',
    value: function buildBeforeEvent() {
      return new Event({
        data: this.data,
        name: 'Before' + this.name
      });
    }
  }, {
    key: 'buildAfterEvent',
    value: function buildAfterEvent() {
      return new Event({
        data: this.data,
        name: 'After' + this.name
      });
    }
  }]);
  return Event;
}();

exports.default = Event;


_lodash2.default.assign(Event, {
  FEATURES_EVENT_NAME: 'Features',
  FEATURES_RESULT_EVENT_NAME: 'FeaturesResult',
  FEATURE_EVENT_NAME: 'Feature',
  SCENARIO_EVENT_NAME: 'Scenario',
  SCENARIO_RESULT_EVENT_NAME: 'ScenarioResult',
  STEP_EVENT_NAME: 'Step',
  STEP_RESULT_EVENT_NAME: 'StepResult'
});