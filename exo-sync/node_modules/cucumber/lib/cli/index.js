'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _slicedToArray2 = require('babel-runtime/helpers/slicedToArray');

var _slicedToArray3 = _interopRequireDefault(_slicedToArray2);

var _bluebird = require('bluebird');

var _bluebird2 = _interopRequireDefault(_bluebird);

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _helpers = require('./helpers');

var _configuration_builder = require('./configuration_builder');

var _configuration_builder2 = _interopRequireDefault(_configuration_builder);

var _builder = require('../formatter/builder');

var _builder2 = _interopRequireDefault(_builder);

var _fs = require('mz/fs');

var _fs2 = _interopRequireDefault(_fs);

var _runtime = require('../runtime');

var _runtime2 = _interopRequireDefault(_runtime);

var _scenario_filter = require('../scenario_filter');

var _scenario_filter2 = _interopRequireDefault(_scenario_filter);

var _support_code_fns = require('../support_code_fns');

var _support_code_fns2 = _interopRequireDefault(_support_code_fns);

var _builder3 = require('../support_code_library/builder');

var _builder4 = _interopRequireDefault(_builder3);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var Cli = function () {
  function Cli(_ref) {
    var argv = _ref.argv,
        cwd = _ref.cwd,
        stdout = _ref.stdout;
    (0, _classCallCheck3.default)(this, Cli);

    this.argv = argv;
    this.cwd = cwd;
    this.stdout = stdout;
  }

  (0, _createClass3.default)(Cli, [{
    key: 'getConfiguration',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* () {
        var fullArgv = yield (0, _helpers.getExpandedArgv)({ argv: this.argv, cwd: this.cwd });
        return yield _configuration_builder2.default.build({ argv: fullArgv, cwd: this.cwd });
      });

      function getConfiguration() {
        return _ref2.apply(this, arguments);
      }

      return getConfiguration;
    }()
  }, {
    key: 'getFormatters',
    value: function () {
      var _ref4 = (0, _bluebird.coroutine)(function* (_ref3) {
        var _this = this;

        var formatOptions = _ref3.formatOptions,
            formats = _ref3.formats,
            supportCodeLibrary = _ref3.supportCodeLibrary;

        var streamsToClose = [];
        var formatters = yield _bluebird2.default.map(formats, function () {
          var _ref6 = (0, _bluebird.coroutine)(function* (_ref5) {
            var _context;

            var type = _ref5.type,
                outputTo = _ref5.outputTo;

            var stream = _this.stdout;
            if (outputTo) {
              var fd = yield _fs2.default.open(outputTo, 'w');
              stream = _fs2.default.createWriteStream(null, { fd: fd });
              streamsToClose.push(stream);
            }
            var typeOptions = _lodash2.default.assign({ log: (_context = stream).write.bind(_context), stream: stream, supportCodeLibrary: supportCodeLibrary }, formatOptions);
            return _builder2.default.build(type, typeOptions);
          });

          return function (_x2) {
            return _ref6.apply(this, arguments);
          };
        }());
        var cleanup = function cleanup() {
          return _bluebird2.default.each(streamsToClose, function (stream) {
            return _bluebird2.default.promisify(stream.end.bind(stream))();
          });
        };
        return { cleanup: cleanup, formatters: formatters };
      });

      function getFormatters(_x) {
        return _ref4.apply(this, arguments);
      }

      return getFormatters;
    }()
  }, {
    key: 'getSupportCodeLibrary',
    value: function getSupportCodeLibrary(supportCodePaths) {
      supportCodePaths.forEach(function (codePath) {
        return require(codePath);
      });
      return _builder4.default.build({ cwd: this.cwd, fns: _support_code_fns2.default.get() });
    }
  }, {
    key: 'run',
    value: function () {
      var _ref7 = (0, _bluebird.coroutine)(function* () {
        var configuration = yield this.getConfiguration();
        var supportCodeLibrary = this.getSupportCodeLibrary(configuration.supportCodePaths);
        var scenarioFilter = new _scenario_filter2.default(configuration.scenarioFilterOptions);

        var _ref8 = yield _bluebird2.default.all([(0, _helpers.getFeatures)({ featurePaths: configuration.featurePaths, scenarioFilter: scenarioFilter }), this.getFormatters({
          formatOptions: configuration.formatOptions,
          formats: configuration.formats,
          supportCodeLibrary: supportCodeLibrary
        })]),
            _ref9 = (0, _slicedToArray3.default)(_ref8, 2),
            features = _ref9[0],
            _ref9$ = _ref9[1],
            cleanup = _ref9$.cleanup,
            formatters = _ref9$.formatters;

        var runtime = new _runtime2.default({
          features: features,
          listeners: formatters,
          options: configuration.runtimeOptions,
          supportCodeLibrary: supportCodeLibrary
        });
        var result = yield runtime.start();
        yield cleanup();
        return result;
      });

      function run() {
        return _ref7.apply(this, arguments);
      }

      return run;
    }()
  }]);
  return Cli;
}();

exports.default = Cli;