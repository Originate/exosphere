require! {
  'chai': {expect}
}

describe 'Logger', ->
  specify '_parse-line', ->
    expect(1).to.eql 0
    # Some possible tests:
    #   bold(), magenta() are chalk functions
    #
    # Example 1
    #   input: 'exo-run', 'Creating space-tweet-web-service'
    #   output: '     bold(exo-run)  Creating space-tweet-web-service'
    # Example 2
    #   input: 'exo-run', 'magenta(exosphere-users-service     |) MongoDB space-tweet-users-dev connected'
    #   output: 'bold(magenta(exosphere-users-service))  MongoDB space-tweet-users-dev connected'
