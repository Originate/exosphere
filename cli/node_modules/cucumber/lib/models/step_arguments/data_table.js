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

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var DataTable = function () {
  function DataTable(gherkinData) {
    (0, _classCallCheck3.default)(this, DataTable);

    this.rawTable = gherkinData.rows.map(function (row) {
      return row.cells.map(function (cell) {
        return cell.value;
      });
    });
  }

  (0, _createClass3.default)(DataTable, [{
    key: 'hashes',
    value: function hashes() {
      var copy = this.raw();
      var keys = copy[0];
      var valuesArray = copy.slice(1);
      return valuesArray.map(function (values) {
        return _lodash2.default.zipObject(keys, values);
      });
    }
  }, {
    key: 'raw',
    value: function raw() {
      return this.rawTable.slice(0);
    }
  }, {
    key: 'rows',
    value: function rows() {
      var copy = this.raw();
      copy.shift();
      return copy;
    }
  }, {
    key: 'rowsHash',
    value: function rowsHash() {
      var rows = this.raw();
      var everyRowHasTwoColumns = _lodash2.default.every(rows, function (row) {
        return row.length === 2;
      });
      if (!everyRowHasTwoColumns) {
        throw new Error('rowsHash can only be called on a data table where all rows have exactly two columns');
      }
      return _lodash2.default.fromPairs(rows);
    }
  }]);
  return DataTable;
}();

exports.default = DataTable;