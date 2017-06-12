require! {
  'chai': {expect}
  './logger': Logger
}

describe 'Logger', ->

  describe '_parse-line', ->
    role = 'exo-run'
    logger = new Logger [role]

    specify 'should parse non-service log message correctly' ->
      line = 'Attaching to exocom0.21.8, web'
      logger._parse-line role, line, (left, right) ->
        expect(left).to.eql(role)
        expect(right).to.eql(line)

    specify 'should parse service log message correctly' ->
      line = '\u001b[33mweb             |\u001b[0m web server running at port 4000'
      result = logger._parse-line role, line, (left, right) ->
      expect(result.left).to.eql('web')
      expect(result.right).to.eql('web server running at port 4000')

    specify 'should strip version from service name' ->
      line = '\u001b[36mexocom0.21.8    |\u001b[0m ExoCom HTTP service online at port 80'
      result = logger._parse-line role, line, (left, right) ->
      expect(result.left).to.eql('exocom')
      expect(result.right).to.eql('ExoCom HTTP service online at port 80')


  describe '_pad', ->
    roles = ['exo-run', 'mongodb-es6']
    logger = new Logger roles

    specify 'should return a padded string of length equal to the length of the longest role' ->
      expect-padded = '    exo-run'
      expect(logger._pad 'exo-run').to.eql(expect-padded)
