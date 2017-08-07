'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _feature = require('../models/feature');

var _feature2 = _interopRequireDefault(_feature);

var _gherkin = require('gherkin');

var _gherkin2 = _interopRequireDefault(_gherkin);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var gherkinCompiler = new _gherkin2.default.Compiler();
var gherkinParser = new _gherkin2.default.Parser();

var Parser = function () {
  function Parser() {
    (0, _classCallCheck3.default)(this, Parser);
  }

  (0, _createClass3.default)(Parser, null, [{
    key: 'parse',
    value: function parse(_ref) {
      var scenarioFilter = _ref.scenarioFilter,
          source = _ref.source,
          uri = _ref.uri;

      var gherkinDocument = void 0;
      try {
        gherkinDocument = gherkinParser.parse(source);
      } catch (error) {
        error.message += '\npath: ' + uri;
        throw error;
      }

      if (gherkinDocument.feature) {
        return new _feature2.default({
          gherkinData: gherkinDocument.feature,
          gherkinPickles: gherkinCompiler.compile(gherkinDocument, uri),
          scenarioFilter: scenarioFilter,
          uri: uri
        });
      }
    }
  }]);
  return Parser;
}();

exports.default = Parser;