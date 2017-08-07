'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});
var previousTimestamp = void 0;

var methods = {
  beginTiming: function beginTiming() {
    previousTimestamp = getTimestamp();
  },

  clearInterval: clearInterval.bind(global),
  clearTimeout: clearTimeout.bind(global),
  Date: Date,
  endTiming: function endTiming() {
    return getTimestamp() - previousTimestamp;
  },

  setInterval: setInterval.bind(global),
  setTimeout: setTimeout.bind(global)
};

if (typeof setImmediate !== 'undefined') {
  methods.setImmediate = setImmediate.bind(global);
  methods.clearImmediate = clearImmediate.bind(global);
}

function getTimestamp() {
  return new methods.Date().getTime();
}

exports.default = methods;