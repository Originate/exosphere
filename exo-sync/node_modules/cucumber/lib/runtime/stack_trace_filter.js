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

var _stackChain = require('stack-chain');

var _stackChain2 = _interopRequireDefault(_stackChain);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var StackTraceFilter = function () {
  function StackTraceFilter() {
    (0, _classCallCheck3.default)(this, StackTraceFilter);

    this.cucumberPath = _path2.default.join(__dirname, '..', '..');
  }

  (0, _createClass3.default)(StackTraceFilter, [{
    key: 'filter',
    value: function filter() {
      var _this = this;

      this.currentFilter = _stackChain2.default.filter.attach(function (error, frames) {
        if (_this.isErrorInCucumber(frames)) {
          return frames;
        }
        var index = _lodash2.default.findIndex(frames, _this.isFrameInCucumber.bind(_this));
        if (index === -1) {
          return frames;
        } else {
          return frames.slice(0, index);
        }
      });
    }
  }, {
    key: 'isErrorInCucumber',
    value: function isErrorInCucumber(frames) {
      var filteredFrames = _lodash2.default.reject(frames, this.isFrameInNode.bind(this));
      return filteredFrames.length > 0 && this.isFrameInCucumber(filteredFrames[0]);
    }
  }, {
    key: 'isFrameInCucumber',
    value: function isFrameInCucumber(frame) {
      var fileName = frame.getFileName() || '';
      return _lodash2.default.startsWith(fileName, this.cucumberPath);
    }
  }, {
    key: 'isFrameInNode',
    value: function isFrameInNode(frame) {
      var fileName = frame.getFileName() || '';
      return !_lodash2.default.includes(fileName, _path2.default.sep);
    }
  }, {
    key: 'unfilter',
    value: function unfilter() {
      _stackChain2.default.filter.deattach(this.currentFilter);
    }
  }]);
  return StackTraceFilter;
}();

exports.default = StackTraceFilter;