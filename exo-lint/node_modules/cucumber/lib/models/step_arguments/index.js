'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _data_table = require('./data_table');

var _data_table2 = _interopRequireDefault(_data_table);

var _doc_string = require('./doc_string');

var _doc_string2 = _interopRequireDefault(_doc_string);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var StepArguments = function () {
  function StepArguments() {
    (0, _classCallCheck3.default)(this, StepArguments);
  }

  (0, _createClass3.default)(StepArguments, null, [{
    key: 'build',
    value: function build(gherkinData) {
      if (gherkinData.hasOwnProperty('content')) {
        return new _doc_string2.default(gherkinData);
      } else if (gherkinData.hasOwnProperty('rows')) {
        return new _data_table2.default(gherkinData);
      } else {
        throw new Error('Unknown step argument type: ' + JSON.stringify(gherkinData));
      }
    }
  }]);
  return StepArguments;
}();

exports.default = StepArguments;