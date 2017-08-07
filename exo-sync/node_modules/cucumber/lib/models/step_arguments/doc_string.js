"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require("babel-runtime/helpers/classCallCheck");

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var DocString = function DocString(gherkinData) {
  (0, _classCallCheck3.default)(this, DocString);

  this.content = gherkinData.content;
  this.contentType = gherkinData.contentType;
  this.line = gherkinData.location.line;
};

exports.default = DocString;