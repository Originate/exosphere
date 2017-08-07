'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _typeof2 = require('babel-runtime/helpers/typeof');

var _typeof3 = _interopRequireDefault(_typeof2);

var _bluebird = require('bluebird');

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _fs = require('mz/fs');

var _fs2 = _interopRequireDefault(_fs);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _stringArgv = require('string-argv');

var _stringArgv2 = _interopRequireDefault(_stringArgv);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var ProfileLoader = function () {
  function ProfileLoader(directory) {
    (0, _classCallCheck3.default)(this, ProfileLoader);

    this.directory = directory;
  }

  (0, _createClass3.default)(ProfileLoader, [{
    key: 'getDefinitions',
    value: function () {
      var _ref = (0, _bluebird.coroutine)(function* () {
        var definitionsFilePath = _path2.default.join(this.directory, 'cucumber.js');
        var exists = yield _fs2.default.exists(definitionsFilePath);
        if (!exists) {
          return {};
        }
        var definitions = require(definitionsFilePath);
        if ((typeof definitions === 'undefined' ? 'undefined' : (0, _typeof3.default)(definitions)) !== 'object') {
          throw new Error(definitionsFilePath + ' does not export an object');
        }
        return definitions;
      });

      function getDefinitions() {
        return _ref.apply(this, arguments);
      }

      return getDefinitions;
    }()
  }, {
    key: 'getArgv',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* (profiles) {
        var definitions = yield this.getDefinitions();
        if (profiles.length === 0 && definitions['default']) {
          profiles = ['default'];
        }
        var argvs = profiles.map(function (profile) {
          if (!definitions[profile]) {
            throw new Error('Undefined profile: ' + profile);
          }
          return (0, _stringArgv2.default)(definitions[profile]);
        });
        return _lodash2.default.flatten(argvs);
      });

      function getArgv(_x) {
        return _ref2.apply(this, arguments);
      }

      return getArgv;
    }()
  }]);
  return ProfileLoader;
}();

exports.default = ProfileLoader;