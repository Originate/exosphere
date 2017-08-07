/**
 * Test case for colorprint.
 * Runs with mocha.
 */
'use strict'

const Colorprint = require('../shim/node/colorprint.js')
const assert = require('assert')

describe('colorpint', () => {
  it('Colorprint', (done) => {
    let colorprint = new Colorprint({})
    colorprint.point('This is point')
    colorprint.notice('This is notice')
    colorprint.info('This is info')
    colorprint.debug('This is debug')
    colorprint.trace('This is trace')
    colorprint.error('This is error')
    colorprint.warn('This is warn')
    colorprint.fatal('This is fatal')

    colorprint.INFO('This is INFO')
    colorprint.DEBUG('This is DEBUG')
    colorprint.TRACE('This is TRACE')
    colorprint.ERROR('This is ERROR')
    colorprint.WARN('This is WARN')
    colorprint.FATAL('This is FATAL')
    done()
  })

  it('Colorprint with indent', (done) => {
    let colorprint = new Colorprint({
      indent: 2
    })
    colorprint.notice('This is indented notice')
    colorprint.info('This is indented info')
    colorprint.debug('This is indented debug')
    colorprint.trace('This is indented trace')
    colorprint.error('This is indented error')
    colorprint.warn('This is indented warn')
    colorprint.fatal('This is indented fatal')

    colorprint.INFO('This is indented INFO')
    colorprint.DEBUG('This is indented DEBUG')
    colorprint.TRACE('This is indented TRACE')
    colorprint.WARN('This is indented WARN')
    colorprint.ERROR('This is indented ERROR')
    colorprint.FATAL('This is indented FATAL')
    done()
  })

  it('Customize color print.', (done) => {
    let colorprint = new Colorprint({
      PREFIX: 'Yeah!',
      SUFFIX: 'That\'s it!',
      INFO_COLOR: 'blue'
    })
    colorprint.info('This will be blue with prefix.')
    done()
  })
})

/* global describe, it */
