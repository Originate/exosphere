'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.validateInstall = undefined;

var _bluebird = require('bluebird');

var validateInstall = exports.validateInstall = function () {
  var _ref = (0, _bluebird.coroutine)(function* (cwd) {
    var projectPath = _path2.default.join(__dirname, '..', '..');
    if (projectPath === cwd) {
      return; // cucumber testing itself
    }
    var currentCucumberPath = require.resolve(projectPath);
    var localCucumberPath = yield (0, _bluebird.promisify)(_resolve2.default)('cucumber', { basedir: cwd });
    localCucumberPath = yield _fs2.default.realpath(localCucumberPath);
    if (localCucumberPath !== currentCucumberPath) {
      throw new Error('\n      You appear to be executing an install of cucumber (most likely a global install)\n      that is different from your local install (the one required in your support files).\n      For cucumber to work, you need to execute the same install that is required in your support files.\n      Please execute the locally installed version to run your tests.\n\n      Executed Path: ' + currentCucumberPath + '\n      Local Path:    ' + localCucumberPath + '\n      ');
    }
  });

  return function validateInstall(_x) {
    return _ref.apply(this, arguments);
  };
}();

var _fs = require('mz/fs');

var _fs2 = _interopRequireDefault(_fs);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _resolve = require('resolve');

var _resolve2 = _interopRequireDefault(_resolve);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }