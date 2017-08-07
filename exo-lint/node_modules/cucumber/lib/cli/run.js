'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _bluebird = require('bluebird');

var _ = require('./');

var _2 = _interopRequireDefault(_);

var _verror = require('verror');

var _verror2 = _interopRequireDefault(_verror);

var _install_validator = require('./install_validator');

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

function exitWithError(error) {
  console.error(_verror2.default.fullStack(error)); // eslint-disable-line no-console
  process.exit(1);
}

exports.default = function () {
  var _ref = (0, _bluebird.coroutine)(function* () {
    var cwd = process.cwd();

    try {
      yield (0, _install_validator.validateInstall)(cwd);
    } catch (error) {
      exitWithError(error);
    }

    var cli = new _2.default({
      argv: process.argv,
      cwd: cwd,
      stdout: process.stdout
    });

    var success = void 0;
    try {
      success = yield cli.run();
    } catch (error) {
      exitWithError(error);
    }

    var exitCode = success ? 0 : 1;
    function exitNow() {
      process.exit(exitCode);
    }

    // If stdout.write() returned false, kernel buffer is not empty yet
    if (process.stdout.write('')) {
      exitNow();
    } else {
      process.stdout.on('drain', exitNow);
    }
  });

  function run() {
    return _ref.apply(this, arguments);
  }

  return run;
}();