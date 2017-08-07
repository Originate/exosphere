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

var _commander = require('commander');

var _package = require('../../package.json');

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var ArgvParser = function () {
  function ArgvParser() {
    (0, _classCallCheck3.default)(this, ArgvParser);
  }

  (0, _createClass3.default)(ArgvParser, null, [{
    key: 'collect',
    value: function collect(val, memo) {
      memo.push(val);
      return memo;
    }
  }, {
    key: 'mergeJson',
    value: function mergeJson(option) {
      return function (str, memo) {
        var val = void 0;
        try {
          val = JSON.parse(str);
        } catch (error) {
          throw new Error(option + ' passed invalid JSON: ' + error.message + ': ' + str);
        }
        if (!_lodash2.default.isPlainObject(val)) {
          throw new Error(option + ' must be passed JSON of an object: ' + str);
        }
        return _lodash2.default.merge(memo, val);
      };
    }
  }, {
    key: 'parse',
    value: function parse(argv) {
      var program = new _commander.Command(_path2.default.basename(argv[1]));

      program.usage('[options] [<DIR|FILE[:LINE]>...]').version(_package.version, '-v, --version').option('-b, --backtrace', 'show full backtrace for errors').option('--compiler <EXTENSION:MODULE>', 'require files with the given EXTENSION after requiring MODULE (repeatable)', ArgvParser.collect, []).option('-d, --dry-run', 'invoke formatters without executing steps').option('--fail-fast', 'abort the run on first failure').option('-f, --format <TYPE[:PATH]>', 'specify the output format, optionally supply PATH to redirect formatter output (repeatable)', ArgvParser.collect, []).option('--format-options <JSON>', 'provide options for formatters (repeatable)', ArgvParser.mergeJson('--format-options'), {}).option('--name <REGEXP>', 'only execute the scenarios with name matching the expression (repeatable)', ArgvParser.collect, []).option('--no-strict', 'succeed even if there are pending or undefined steps').option('-p, --profile <NAME>', 'specify the profile to use (repeatable)', ArgvParser.collect, []).option('-r, --require <FILE|DIR>', 'require files before executing features (repeatable)', ArgvParser.collect, []).option('-t, --tags <EXPRESSION>', 'only execute the features or scenarios with tags matching the expression', '').option('--world-parameters <JSON>', 'provide parameters that will be passed to the world constructor (repeatable)', ArgvParser.mergeJson('--world-parameters'), {});

      program.on('--help', function () {
        /* eslint-disable no-console */
        console.log('  For more details please visit https://github.com/cucumber/cucumber-js#cli\n');
        /* eslint-enable no-console */
      });

      program.parse(argv);

      return {
        options: program.opts(),
        args: program.args
      };
    }
  }]);
  return ArgvParser;
}();

exports.default = ArgvParser;