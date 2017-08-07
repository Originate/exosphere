"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require("babel-runtime/helpers/classCallCheck");

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require("babel-runtime/helpers/createClass");

var _createClass3 = _interopRequireDefault(_createClass2);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Tag = function () {
  (0, _createClass3.default)(Tag, null, [{
    key: "build",
    value: function build(gherkinData) {
      return new Tag(gherkinData);
    }
  }]);

  function Tag(gherkinData) {
    (0, _classCallCheck3.default)(this, Tag);

    this.line = gherkinData.location.line;
    this.name = gherkinData.name;
  }

  return Tag;
}();

exports.default = Tag;