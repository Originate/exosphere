"use strict";

const colorpint = require('colorpint')

colorpint.notice('This is NOTICE') //Pipe to stdout with magenta color.
colorpint.info('This is INFO') //Pipe to stdout with green color.
colorpint.debug('This is DEBUG') //Pipe to stdout with  color.
colorpint.trace('This is TRACE') //Pipe to stdout with white color.
colorpint.warn('This is WARN') //Pipe to stdout with yellow color.
colorpint.error('This is ERROR') //Pipe to stderr with red color.
colorpint.fatal('This is FATAL') //Pipe to stderr with bgRed color.
