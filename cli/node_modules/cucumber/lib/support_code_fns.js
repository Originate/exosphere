"use strict";

Object.defineProperty(exports, "__esModule", {
  value: true
});
var fns = [];

exports.default = {
  add: function add(fn) {
    fns.push(fn);
  },
  get: function get() {
    return fns;
  },
  reset: function reset() {
    fns = [];
  }
};