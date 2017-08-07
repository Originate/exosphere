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

var _path = require('path');

var _path2 = _interopRequireDefault(_path);

var _tag_expression_parser = require('cucumber-tag-expressions/lib/tag_expression_parser');

var _tag_expression_parser2 = _interopRequireDefault(_tag_expression_parser);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var FEATURE_LINENUM_REGEXP = /^(.*?)((?::[\d]+)+)?$/;
var tagExpressionParser = new _tag_expression_parser2.default();

var ScenarioFilter = function () {
  function ScenarioFilter(_ref) {
    var cwd = _ref.cwd,
        featurePaths = _ref.featurePaths,
        names = _ref.names,
        tagExpression = _ref.tagExpression;
    (0, _classCallCheck3.default)(this, ScenarioFilter);

    this.cwd = cwd;
    this.featureUriToLinesMapping = this.getFeatureUriToLinesMapping(featurePaths || []);
    this.names = names || [];
    if (tagExpression) {
      this.tagExpressionNode = tagExpressionParser.parse(tagExpression || '');
    }
  }

  (0, _createClass3.default)(ScenarioFilter, [{
    key: 'getFeatureUriToLinesMapping',
    value: function getFeatureUriToLinesMapping(featurePaths) {
      var _this = this;

      var mapping = {};
      featurePaths.forEach(function (featurePath) {
        var match = FEATURE_LINENUM_REGEXP.exec(featurePath);
        if (match) {
          var uri = _path2.default.resolve(_this.cwd, match[1]);
          var linesExpression = match[2];
          if (linesExpression) {
            if (!mapping[uri]) {
              mapping[uri] = [];
            }
            linesExpression.slice(1).split(':').forEach(function (line) {
              mapping[uri].push(parseInt(line));
            });
          }
        }
      });
      return mapping;
    }
  }, {
    key: 'matches',
    value: function matches(scenario) {
      return this.matchesAnyLine(scenario) && this.matchesAnyName(scenario) && this.matchesAllTagExpressions(scenario);
    }
  }, {
    key: 'matchesAnyLine',
    value: function matchesAnyLine(scenario) {
      var lines = this.featureUriToLinesMapping[scenario.uri];
      if (lines) {
        return _lodash2.default.size(_lodash2.default.intersection(lines, scenario.lines)) > 0;
      } else {
        return true;
      }
    }
  }, {
    key: 'matchesAnyName',
    value: function matchesAnyName(scenario) {
      if (this.names.length === 0) {
        return true;
      }
      var scenarioName = scenario.name;
      return _lodash2.default.some(this.names, function (name) {
        return scenarioName.match(name);
      });
    }
  }, {
    key: 'matchesAllTagExpressions',
    value: function matchesAllTagExpressions(scenario) {
      if (!this.tagExpressionNode) {
        return true;
      }
      var scenarioTags = scenario.tags.map(function (t) {
        return t.name;
      });
      return this.tagExpressionNode.evaluate(scenarioTags);
    }
  }]);
  return ScenarioFilter;
}();

exports.default = ScenarioFilter;