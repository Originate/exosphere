'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.formatLocation = formatLocation;

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function formatLocation(cwd, obj) {
  return _path2.default.relative(cwd, obj.uri) + ':' + obj.line;
}