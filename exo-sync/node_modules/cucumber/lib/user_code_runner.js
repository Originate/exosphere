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

var _time = require('./time');

var _time2 = _interopRequireDefault(_time);

var _uncaught_exception_manager = require('./uncaught_exception_manager');

var _uncaught_exception_manager2 = _interopRequireDefault(_uncaught_exception_manager);

var _util = require('util');

var _util2 = _interopRequireDefault(_util);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var UserCodeRunner = function () {
  function UserCodeRunner() {
    (0, _classCallCheck3.default)(this, UserCodeRunner);
  }

  (0, _createClass3.default)(UserCodeRunner, null, [{
    key: 'run',
    value: function () {
      var _ref2 = (0, _bluebird.coroutine)(function* (_ref) {
        var argsArray = _ref.argsArray,
            thisArg = _ref.thisArg,
            fn = _ref.fn,
            timeoutInMilliseconds = _ref.timeoutInMilliseconds;

        var callbackPromise = new _bluebird2.default(function (resolve, reject) {
          argsArray.push(function (error, result) {
            if (error) {
              reject(error);
            } else {
              resolve(result);
            }
          });
        });

        var fnReturn = void 0;
        try {
          fnReturn = fn.apply(thisArg, argsArray);
        } catch (e) {
          var _error = e instanceof Error ? e : new Error(_util2.default.format(e));
          return { error: _error };
        }

        var racingPromises = [];
        var callbackInterface = fn.length === argsArray.length;
        var promiseInterface = fnReturn && typeof fnReturn.then === 'function';

        if (callbackInterface && promiseInterface) {
          return { error: new Error('function uses multiple asynchronous interfaces: callback and promise') };
        } else if (callbackInterface) {
          racingPromises.push(callbackPromise);
        } else if (promiseInterface) {
          racingPromises.push(fnReturn);
        } else {
          return { result: fnReturn };
        }

        var exceptionHandler = void 0;
        var uncaughtExceptionPromise = new _bluebird2.default(function (resolve, reject) {
          exceptionHandler = reject;
          _uncaught_exception_manager2.default.registerHandler(exceptionHandler);
        });
        racingPromises.push(uncaughtExceptionPromise);

        var timeoutId = void 0;
        if (timeoutInMilliseconds >= 0) {
          var timeoutPromise = new _bluebird2.default(function (resolve, reject) {
            timeoutId = _time2.default.setTimeout(function () {
              var timeoutMessage = 'function timed out after ' + timeoutInMilliseconds + ' milliseconds';
              reject(new Error(timeoutMessage));
            }, timeoutInMilliseconds);
          });
          racingPromises.push(timeoutPromise);
        }

        var error = void 0,
            result = void 0;
        try {
          result = yield _bluebird2.default.race(racingPromises);
        } catch (e) {
          if (e instanceof Error) {
            error = e;
          } else if (e) {
            error = new Error(_util2.default.format(e));
          } else {
            error = new Error('Promise rejected without a reason');
          }
        }

        _time2.default.clearTimeout(timeoutId);
        _uncaught_exception_manager2.default.unregisterHandler(exceptionHandler);

        return { error: error, result: result };
      });

      function run(_x) {
        return _ref2.apply(this, arguments);
      }

      return run;
    }()
  }]);
  return UserCodeRunner;
}();

exports.default = UserCodeRunner;