"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require("babel-runtime/helpers/classCallCheck");

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Attachment = function Attachment(_ref) {
  var data = _ref.data,
      mimeType = _ref.mimeType;
  (0, _classCallCheck3.default)(this, Attachment);

  this.data = data;
  this.mimeType = mimeType;
};

exports.default = Attachment;