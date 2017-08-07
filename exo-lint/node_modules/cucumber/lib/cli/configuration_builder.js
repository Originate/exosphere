'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _bluebird = require('bluebird');

var _bluebird2 = _interopRequireDefault(_bluebird);

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _argv_parser = require('./argv_parser');

var _argv_parser2 = _interopRequireDefault(_argv_parser);

var _fs = require('mz/fs');

var _fs2 = _interopRequireDefault(_fs);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _path_expander = require('./path_expander');

var _path_expander2 = _interopRequireDefault(_path_expander);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var ConfigurationBuilder = function () {
  (0, _createClass3.default)(ConfigurationBuilder, null, [{
    key: 'build',
    value: function () {
      var _ref = (0, _bluebird.coroutine)(function* (options) {
        var builder = new ConfigurationBuilder(options);
        return yield builder.build();
      });

      function build(_x) {
        return _ref.apply(this, arguments);
      }

      return build;
    }()
  }]);

  function ConfigurationBuilder(_ref2) {
    var argv = _ref2.argv,
        cwd = _ref2.cwd;
    (0, _classCallCheck3.default)(this, ConfigurationBuilder);

    this.cwd = cwd;
    this.pathExpander = new _path_expander2.default(cwd);

    var parsedArgv = _argv_parser2.default.parse(argv);
    this.args = parsedArgv.args;
    this.options = parsedArgv.options;
  }

  (0, _createClass3.default)(ConfigurationBuilder, [{
    key: 'build',
    value: function () {
      var _ref3 = (0, _bluebird.coroutine)(function* () {
        var unexpandedFeaturePaths = yield this.getUnexpandedFeaturePaths();
        var featurePaths = yield this.expandFeaturePaths(unexpandedFeaturePaths);
        var featureDirectoryPaths = this.getFeatureDirectoryPaths(featurePaths);
        var unexpandedSupportCodePaths = this.options.require.length > 0 ? this.options.require : featureDirectoryPaths;
        var supportCodePaths = yield this.expandSupportCodePaths(unexpandedSupportCodePaths);
        return {
          featurePaths: featurePaths,
          formats: this.getFormats(),
          formatOptions: this.getFormatOptions(),
          profiles: this.options.profile,
          runtimeOptions: {
            dryRun: !!this.options.dryRun,
            failFast: !!this.options.failFast,
            filterStacktraces: !this.options.backtrace,
            strict: !!this.options.strict,
            worldParameters: this.options.worldParameters
          },
          scenarioFilterOptions: {
            cwd: this.cwd,
            featurePaths: unexpandedFeaturePaths,
            names: this.options.name,
            tagExpression: this.options.tags
          },
          supportCodePaths: supportCodePaths
        };
      });

      function build() {
        return _ref3.apply(this, arguments);
      }

      return build;
    }()
  }, {
    key: 'expandFeaturePaths',
    value: function () {
      var _ref4 = (0, _bluebird.coroutine)(function* (featurePaths) {
        featurePaths = featurePaths.map(function (p) {
          return p.replace(/(:\d+)*$/g, '');
        }); // Strip line numbers
        return yield this.pathExpander.expandPathsWithExtensions(featurePaths, ['feature']);
      });

      function expandFeaturePaths(_x2) {
        return _ref4.apply(this, arguments);
      }

      return expandFeaturePaths;
    }()
  }, {
    key: 'getFeatureDirectoryPaths',
    value: function getFeatureDirectoryPaths(featurePaths) {
      var _this = this;

      var featureDirs = featurePaths.map(function (featurePath) {
        var featureDir = _path2.default.dirname(featurePath);
        var childDir = void 0;
        var parentDir = featureDir;
        while (childDir !== parentDir) {
          childDir = parentDir;
          parentDir = _path2.default.dirname(childDir);
          if (_path2.default.basename(parentDir) === 'features') {
            featureDir = parentDir;
            break;
          }
        }
        return _path2.default.relative(_this.cwd, featureDir);
      });
      return _lodash2.default.uniq(featureDirs);
    }
  }, {
    key: 'getFormatOptions',
    value: function getFormatOptions() {
      var formatOptions = _lodash2.default.clone(this.options.formatOptions);
      formatOptions.cwd = this.cwd;
      _lodash2.default.defaults(formatOptions, { colorsEnabled: true });
      return formatOptions;
    }
  }, {
    key: 'getFormats',
    value: function getFormats() {
      var mapping = { '': 'pretty' };
      this.options.format.forEach(function (format) {
        var parts = format.split(':');
        var type = parts[0];
        var outputTo = parts.slice(1).join(':');
        mapping[outputTo] = type;
      });
      return _lodash2.default.map(mapping, function (type, outputTo) {
        return { outputTo: outputTo, type: type };
      });
    }
  }, {
    key: 'getUnexpandedFeaturePaths',
    value: function () {
      var _ref5 = (0, _bluebird.coroutine)(function* () {
        var _this2 = this;

        if (this.args.length > 0) {
          var nestedFeaturePaths = yield _bluebird2.default.map(this.args, function () {
            var _ref6 = (0, _bluebird.coroutine)(function* (arg) {
              var filename = _path2.default.basename(arg);
              if (filename[0] === '@') {
                var filePath = _path2.default.join(_this2.cwd, arg);
                var content = yield _fs2.default.readFile(filePath, 'utf8');
                return _lodash2.default.chain(content).split('\n').map(_lodash2.default.trim).compact().value();
              } else {
                return arg;
              }
            });

            return function (_x3) {
              return _ref6.apply(this, arguments);
            };
          }());
          var featurePaths = _lodash2.default.flatten(nestedFeaturePaths);
          if (featurePaths.length > 0) {
            return featurePaths;
          }
        }
        return ['features'];
      });

      function getUnexpandedFeaturePaths() {
        return _ref5.apply(this, arguments);
      }

      return getUnexpandedFeaturePaths;
    }()
  }, {
    key: 'expandSupportCodePaths',
    value: function () {
      var _ref7 = (0, _bluebird.coroutine)(function* (supportCodePaths) {
        var extensions = ['js'];
        this.options.compiler.forEach(function (compiler) {
          var parts = compiler.split(':');
          extensions.push(parts[0]);
          require(parts[1]);
        });
        return yield this.pathExpander.expandPathsWithExtensions(supportCodePaths, extensions);
      });

      function expandSupportCodePaths(_x4) {
        return _ref7.apply(this, arguments);
      }

      return expandSupportCodePaths;
    }()
  }]);
  return ConfigurationBuilder;
}();

exports.default = ConfigurationBuilder;