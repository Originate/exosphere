// Generated by LiveScript 1.5.0
var AppRunner, ref$, bold, cyan, dim, green, red, fs, yaml, Logger, flatten, util;
AppRunner = require('./app-runner');
ref$ = require('chalk'), bold = ref$.bold, cyan = ref$.cyan, dim = ref$.dim, green = ref$.green, red = ref$.red;
fs = require('fs');
yaml = require('js-yaml');
Logger = require('../../exosphere-shared').Logger;
flatten = require('prelude-ls').flatten;
util = require('util');
module.exports = function(){
  var appConfig, services, silencedServices, type, service, silencedDependencies, res$, i$, ref$, len$, dependency, logger, x$, appRunner, this$ = this;
  if (process.argv[2] === "help") {
    help();
    return;
  }
  appConfig = yaml.safeLoad(fs.readFileSync('application.yml', 'utf8'));
  console.log("Running " + green(appConfig.name) + " " + cyan(appConfig.version) + "\n");
  services = [];
  silencedServices = [];
  for (type in appConfig.services) {
    for (service in appConfig.services[type]) {
      services.push(service);
      if (appConfig.services[type][service].silent) {
        silencedServices.push(service);
      }
    }
  }
  res$ = [];
  for (i$ = 0, len$ = (ref$ = appConfig.dependencies).length; i$ < len$; ++i$) {
    dependency = ref$[i$];
    if (dependency.silent) {
      res$.push(dependency.name);
    }
  }
  silencedDependencies = res$;
  logger = new Logger(services, silencedServices.concat(silencedDependencies));
  x$ = appRunner = new AppRunner({
    appConfig: appConfig,
    logger: logger
  });
  x$.on('routing-setup', function(){
    var command, ref$, routing, text, receivers, res$, i$, ref1$, len$, receiver, results$ = [];
    logger.log({
      role: 'exocom',
      text: 'received routing setup'
    });
    for (command in ref$ = appRunner.exocom.clientRegistry.routes) {
      routing = ref$[command];
      text = "  [ " + bold(command) + " ]  -->  ";
      res$ = [];
      for (i$ = 0, len$ = (ref1$ = routing.receivers).length; i$ < len$; ++i$) {
        receiver = ref1$[i$];
        res$.push(bold(receiver.name) + " (" + receiver.host + ":" + receiver.port + ")");
      }
      receivers = res$;
      text += receivers.join(' & ');
      results$.push(logger.log({
        role: 'exocom',
        text: text
      }));
    }
    return results$;
  });
  x$.on('message', function(arg$){
    var messages, receivers, message, indent, i$, ref$, len$, line, results$ = [];
    messages = arg$.messages, receivers = arg$.receivers;
    message = messages[0];
    if (message.name !== message.originalName) {
      logger.log({
        role: 'exocom',
        text: bold(message.sender) + "  --[ " + bold(message.originalName) + " ]-[ " + bold(message.name) + " ]->  " + bold(receivers.join(' and '))
      });
    } else {
      logger.log({
        role: 'exocom',
        text: bold(message.sender) + "  --[ " + bold(message.name) + " ]->  " + bold(receivers.join(' and '))
      });
    }
    indent = repeatString$(' ', message.sender.length + 2);
    if (message.payload != null) {
      for (i$ = 0, len$ = (ref$ = util.inspect(message.payload, {
        showHidden: false,
        depth: null
      }).split('\n')).length; i$ < len$; ++i$) {
        line = ref$[i$];
        results$.push(logger.log({
          role: 'exocom',
          text: indent + "" + dim(line),
          trim: false
        }));
      }
      return results$;
    } else {
      return logger.log({
        role: 'exocom',
        text: indent + "" + dim('(no payload)'),
        trim: false
      });
    }
  });
  x$.start();
  return process.on('SIGINT', function(){
    return appRunner.shutdown({
      closeMessage: " shutting down ..."
    });
  });
};
function help(){
  var helpMessage;
  helpMessage = "Usage: " + cyan("exo run") + "\n\nRuns an Exosphere application.\nThis command must be run in the root directory of the application.";
  return console.log(helpMessage);
}
function repeatString$(str, n){
  for (var r = ''; n > 0; (n >>= 1) && (str += str)) if (n & 1) r += str;
  return r;
}