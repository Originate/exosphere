require! {
  'cucumber': {defineSupportCode}
  'path'
}


defineSupportCode ({Given}) ->

  Given /^a running "([^"]*)" application$/ timeout: 600_000, (@app-name, done) ->
    @checkout-and-run-app {online-text: 'all services online'}, done 


