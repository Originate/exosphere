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

var _fs = require('mz/fs');

var _fs2 = _interopRequireDefault(_fs);

var _glob = require('glob');

var _glob2 = _interopRequireDefault(_glob);

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var PathExpander = function () {
  function PathExpander(directory) {
    (0, _classCallCheck3.default)(this, PathExpander);

    this.directory = directory;
  }

  (0, _createClass3.default)(PathExpander, [{
    key: 'expandPathsWithExtensions',
    value: function () {
      var _ref = (0, _bluebird.coroutine)(function* (paths, extensions) {
        var _this = this;

        var expandedPaths = yield _bluebird2.default.map(paths, function () {
          var _ref2 = (0, _bluebird.coroutine)(function* (p) {
            return yield _this.expandPathWithExtensions(p, extensions);
          });

          return function (_x3) {
            return _ref2.apply(this, arguments);
          };
        }());
        return _lodash2.default.uniq(_lodash2.default.flatten(expandedPaths));
      });

      function expandPathsWithExtensions(_x, _x2) {
        return _ref.apply(this, arguments);
      }

      return expandPathsWithExtensions;
    }()
  }, {
    key: 'expandPathWithExtensions',
    value: function () {
      var _ref3 = (0, _bluebird.coroutine)(function* (p, extensions) {
        var fullPath = _path2.default.resolve(this.directory, p);
        var stats = yield _fs2.default.stat(fullPath);
        if (stats.isDirectory()) {
          return yield this.expandDirectoryWithExtensions(fullPath, extensions);
        } else {
          return [fullPath];
        }
      });

      function expandPathWithExtensions(_x4, _x5) {
        return _ref3.apply(this, arguments);
      }

      return expandPathWithExtensions;
    }()
  }, {
    key: 'expandDirectoryWithExtensions',
    value: function () {
      var _ref4 = (0, _bluebird.coroutine)(function* (realPath, extensions) {
        var pattern = realPath + '/**/*.';
        if (extensions.length > 1) {
          pattern += '{' + extensions.join(',') + '}';
        } else {
          pattern += extensions[0];
        }
        var results = yield _bluebird2.default.promisify(_glob2.default)(pattern);
        return results.map(function (filePath) {
          return filePath.replace(/\//g, _path2.default.sep);
        });
      });

      function expandDirectoryWithExtensions(_x6, _x7) {
        return _ref4.apply(this, arguments);
      }

      return expandDirectoryWithExtensions;
    }()
  }]);
  return PathExpander;
}();

exports.default = PathExpander;