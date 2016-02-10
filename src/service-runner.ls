require! {
  'chalk' : {white}
  'js-yaml' : yaml
  'fs'
  'observable-process' : ObservableProcess
  'path'
}


class ServiceRunner

  (@name, @config, @fg-color, @bg-color) ->


  run: ->
    @service-config = yaml.safeLoad fs.readFileSync(path.join(process.cwd!, @name, 'config.yml'), 'utf8')
    @process = new ObservableProcess(@_create-start-command(@service-config['start-command'])
                                     cwd: path.join(process.cwd!, @name),
                                     verbose: yes,
                                     console: log: @log, error: @log)


  _create-start-command: (template) ->
    for key, value of @config
      template = template.replace "{{#{key}}}", value
    template


  log: (text) ~>
    console.log @bg-color(white(" #{@name.to-upper-case!} ")), @fg-color(text.trim!)



module.exports = ServiceRunner
