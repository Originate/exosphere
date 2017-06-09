require! {
  'chai': {expect}
  'chalk' : {yellow}
  './logger': Logger
}

describe 'Logger', ->

  describe '_parse-line', ->
    logger = new Logger [role]
    role = 'exo-run'

    specify 'should parse non-service log message correctly' ->
      line = 'Attaching to exocom0.21.8, web'
      logger._parse-line role, line, (left, right) ->
        expect(left).to.eql(role)
        expect(right).to.eql(line)

    specify 'should parse service log message correctly' ->
      line = '\u001b[33mweb             |\u001b[0m web server running at port 4000'
      logger._parse-line role, line, (left, right) ->
        expect(left).to.eql('\u001b[33mweb')
        expect(right).to.eql('\u001b[0mweb server running at port 4000')

    specify 'should strip version from service name' ->
      line = '\u001b[36mexocom0.21.8    |\u001b[0m ExoCom HTTP service online at port 80'
      logger._parse-line role, line, (left, right) ->
        expect(left).to.eql('\u001b[36mexocom')
        expect(right).to.eql('\u001b[0mExoCom HTTP service online at port 80')


  describe '_pad', ->
    role = 'exo-run'
    logger = new Logger [role]

    specify 'should pad non-styled string correctly' ->
      expect-padded = ' ' * (logger.length - role.length) + role
      expect(logger._pad role).to.eql(expect-padded)

    specify 'should pad styled string correctly and preserve color' ->
      styled-service-name = '\u001b[33mweb'
      reset-ansi = '\u001b[39m'
      expect-padded = (yellow ' ' * (logger.length - 'web'.length) + 'web').replace reset-ansi, ''
      expect(logger._pad styled-service-name).to.eql(expect-padded)

    # Some possible tests:
    #   bold(), magenta() are chalk functions
    #
    # Example 1
    #   input: 'exo-run', 'Creating space-tweet-web-service'
    #   output: '     bold(exo-run)  Creating space-tweet-web-service'
    # Example 2
    #   input: 'exo-run', 'magenta(exosphere-users-service     |) MongoDB space-tweet-users-dev connected'
      # output: 'bold(magenta(exosphere-users-service))  MongoDB space-tweet-users-dev connected'
