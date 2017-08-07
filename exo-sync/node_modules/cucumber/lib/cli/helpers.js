'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
exports.getFeatures = exports.getExpandedArgv = undefined;

var _bluebird = require('bluebird');

var _bluebird2 = _interopRequireDefault(_bluebird);

var getExpandedArgv = exports.getExpandedArgv = function () {
  var _ref2 = (0, _bluebird.coroutine)(function* (_ref) {
    var argv = _ref.argv,
        cwd = _ref.cwd;

    var _ArgvParser$parse = _argv_parser2.default.parse(argv),
        options = _ArgvParser$parse.options;

    var fullArgv = argv;
    var profileArgv = yield new _profile_loader2.default(cwd).getArgv(options.profile);
    if (profileArgv.length > 0) {
      fullArgv = _lodash2.default.concat(argv.slice(0, 2), profileArgv, argv.slice(2));
    }
    return fullArgv;
  });

  return function getExpandedArgv(_x) {
    return _ref2.apply(this, arguments);
  };
}();

var getFeatures = exports.getFeatures = function () {
  var _ref4 = (0, _bluebird.coroutine)(function* (_ref3) {
    var featurePaths = _ref3.featurePaths,
        scenarioFilter = _ref3.scenarioFilter;

    var features = yield _bluebird2.default.map(featurePaths, function () {
      var _ref5 = (0, _bluebird.coroutine)(function* (featurePath) {
        var source = yield _fs2.default.readFile(featurePath, 'utf8');
        return _feature_parser2.default.parse({ scenarioFilter: scenarioFilter, source: source, uri: featurePath });
      });

      return function (_x3) {
        return _ref5.apply(this, arguments);
      };
    }());
    return _lodash2.default.chain(features).compact().filter(function (feature) {
      return feature.scenarios.length > 0;
    }).value();
  });

  return function getFeatures(_x2) {
    return _ref4.apply(this, arguments);
  };
}();

var _lodash = require('lodash');

var _lodash2 = _interopRequireDefault(_lodash);

var _argv_parser = require('./argv_parser');

var _argv_parser2 = _interopRequireDefault(_argv_parser);

var _fs = require('mz/fs');

var _fs2 = _interopRequireDefault(_fs);

var _feature_parser = require('./feature_parser');

var _feature_parser2 = _interopRequireDefault(_feature_parser);

var _profile_loader = require('./profile_loader');

var _profile_loader2 = _interopRequireDefault(_profile_loader);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }