'use strict';

Object.defineProperty(exports, "__esModule", {
  value: true
});

var _classCallCheck2 = require('babel-runtime/helpers/classCallCheck');

var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);

var _createClass2 = require('babel-runtime/helpers/createClass');

var _createClass3 = _interopRequireDefault(_createClass2);

var _attachment = require('./attachment');

var _attachment2 = _interopRequireDefault(_attachment);

var _isStream = require('is-stream');

var _isStream2 = _interopRequireDefault(_isStream);

var _bluebird = require('bluebird');

var _bluebird2 = _interopRequireDefault(_bluebird);

function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { default: obj }; }

var AttachmentManager = function () {
  function AttachmentManager() {
    (0, _classCallCheck3.default)(this, AttachmentManager);

    this.attachments = [];
  }

  (0, _createClass3.default)(AttachmentManager, [{
    key: 'create',
    value: function create(data, mimeType, callback) {
      if (Buffer.isBuffer(data)) {
        if (!mimeType) {
          throw Error('Buffer attachments must specify a mimeType');
        }
        this.createBufferAttachment(data, mimeType);
      } else if (_isStream2.default.readable(data)) {
        if (!mimeType) {
          throw Error('Stream attachments must specify a mimeType');
        }
        return this.createStreamAttachment(data, mimeType, callback);
      } else if (typeof data === 'string') {
        if (!mimeType) {
          mimeType = 'text/plain';
        }
        this.createStringAttachment(data, mimeType);
      } else {
        throw Error('Invalid attachment data: must be a buffer, readable stream, or string');
      }
    }
  }, {
    key: 'createBufferAttachment',
    value: function createBufferAttachment(data, mimeType) {
      this.createStringAttachment(data.toString('base64'), mimeType);
    }
  }, {
    key: 'createStreamAttachment',
    value: function createStreamAttachment(data, mimeType, callback) {
      var _this = this;

      var promise = new _bluebird2.default(function (resolve, reject) {
        var buffers = [];
        data.on('data', function (chunk) {
          buffers.push(chunk);
        });
        data.on('end', function () {
          _this.createBufferAttachment(Buffer.concat(buffers), mimeType);
          resolve();
        });
        data.on('error', reject);
      });
      if (callback) {
        promise.then(callback, callback);
      } else {
        return promise;
      }
    }
  }, {
    key: 'createStringAttachment',
    value: function createStringAttachment(data, mimeType) {
      var attachment = new _attachment2.default({ data: data, mimeType: mimeType });
      this.attachments.push(attachment);
    }
  }, {
    key: 'getAll',
    value: function getAll() {
      return this.attachments;
    }
  }, {
    key: 'reset',
    value: function reset() {
      this.attachments = [];
    }
  }]);
  return AttachmentManager;
}();

exports.default = AttachmentManager;