// Generated by LiveScript 1.5.0
var Asynchronizer, red, EventEmitter, ref$, ApplicationDependency, DockerCompose, fs, path, ServiceRestarter, yaml, AppRunner;
Asynchronizer = require('asynchronizer');
red = require('chalk').red;
EventEmitter = require('events').EventEmitter;
ref$ = require('../../exosphere-shared'), ApplicationDependency = ref$.ApplicationDependency, DockerCompose = ref$.DockerCompose;
fs = require('fs');
path = require('path');
ServiceRestarter = require('./service-restarter');
yaml = require('js-yaml');
AppRunner = (function(superclass){
  var prototype = extend$((import$(AppRunner, superclass).displayName = 'AppRunner', AppRunner), superclass).prototype, constructor = AppRunner;
  function AppRunner(arg$){
    var i$, ref$, len$, dependencyConfig, dependency, ref1$;
    this.appConfig = arg$.appConfig, this.logger = arg$.logger;
    this.write = bind$(this, 'write', prototype);
    this.shutdown = bind$(this, 'shutdown', prototype);
    this.env = {};
    for (i$ = 0, len$ = (ref$ = this.appConfig.dependencies).length; i$ < len$; ++i$) {
      dependencyConfig = ref$[i$];
      dependency = ApplicationDependency.build(dependencyConfig);
      this.env = (ref1$ = {}, import$(ref1$, this.env), import$(ref1$, dependency.getEnvVariables()));
    }
    this.dockerConfigLocation = path.join(process.cwd(), 'tmp');
  }
  AppRunner.prototype.start = function(){
    var onlineTexts, asynchronizer, role, onlineText, this$ = this;
    this.watchServices();
    this.process = DockerCompose.runAllImages({
      env: this.env,
      cwd: this.dockerConfigLocation,
      write: this.write
    }, function(exitCode){
      switch (false) {
      case !exitCode:
        return this$.shutdown({
          errorMessage: 'Failed to run images'
        });
      }
    });
    onlineTexts = this._compileOnlineText();
    asynchronizer = new Asynchronizer(Object.keys(onlineTexts));
    for (role in onlineTexts) {
      onlineText = onlineTexts[role];
      (fn$.call(this, role, onlineText));
    }
    return asynchronizer.then(function(){
      return this$.write('all services online');
    });
    function fn$(role, onlineText){
      var this$ = this;
      this.process.wait(new RegExp(role + ".*" + onlineText), function(){
        this$.logger.log({
          role: role,
          text: "'" + role + "' is running"
        });
        return asynchronizer.check(role);
      });
    }
  };
  AppRunner.prototype.watchServices = function(){
    var protectionLevel, lresult$, role, ref$, serviceData, x$, results$ = [], this$ = this;
    this.services = [];
    for (protectionLevel in this.appConfig.services) {
      lresult$ = [];
      for (role in ref$ = this.appConfig.services[protectionLevel]) {
        serviceData = ref$[role];
        if (serviceData.location) {
          x$ = new ServiceRestarter({
            role: role,
            serviceLocation: path.join(process.cwd(), serviceData.location),
            env: this.env,
            logger: this.logger
          });
          x$.watch();
          x$.on('error', fn$);
          lresult$.push(x$);
        }
      }
      results$.push(lresult$);
    }
    return results$;
    function fn$(message){
      return this$.shutdown({
        errorMessage: message
      });
    }
  };
  AppRunner.prototype.shutdown = function(arg$){
    var closeMessage, errorMessage, exitCode;
    closeMessage = arg$.closeMessage, errorMessage = arg$.errorMessage;
    switch (false) {
    case !errorMessage:
      console.log(red(errorMessage));
      exitCode = 1;
      break;
    default:
      console.log("\n\n " + closeMessage);
      exitCode = 0;
    }
    return DockerCompose.killAllContainers({
      cwd: this.dockerConfigLocation,
      write: this.write
    }, function(){
      return process.exit(exitCode);
    });
  };
  AppRunner.prototype._compileOnlineText = function(){
    var onlineTexts, protectionLevel, role, ref$, serviceData, serviceConfig;
    onlineTexts = {};
    for (protectionLevel in this.appConfig.services) {
      for (role in ref$ = this.appConfig.services[protectionLevel]) {
        serviceData = ref$[role];
        if (serviceData.location) {
          serviceConfig = yaml.safeLoad(fs.readFileSync(path.join(process.cwd(), serviceData.location, 'service.yml')));
          onlineTexts[role] = serviceConfig.startup['online-text'];
        }
      }
    }
    return onlineTexts;
  };
  AppRunner.prototype.write = function(text){
    return this.logger.log({
      role: 'exo-run',
      text: text,
      trim: true
    });
  };
  return AppRunner;
}(EventEmitter));
module.exports = AppRunner;
function bind$(obj, key, target){
  return function(){ return (target || obj)[key].apply(obj, arguments) };
}
function extend$(sub, sup){
  function fun(){} fun.prototype = (sub.superclass = sup).prototype;
  (sub.prototype = new fun).constructor = sub;
  if (typeof sup.extended == 'function') sup.extended(sub);
  return sub;
}
function import$(obj, src){
  var own = {}.hasOwnProperty;
  for (var key in src) if (own.call(src, key)) obj[key] = src[key];
  return obj;
}