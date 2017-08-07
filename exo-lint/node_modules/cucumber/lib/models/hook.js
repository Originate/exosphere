'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Hook = function Hook(_ref) {
  var keyword = _ref.keyword,
      scenario = _ref.scenario;
  (0, _classCallCheck3.default)(this, Hook);

  this.keyword = keyword;
  this.scenario = scenario;
};

exports.default = Hook;


Hook.BEFORE_STEP_KEYWORD = 'Before';
Hook.AFTER_STEP_KEYWORD = 'After';