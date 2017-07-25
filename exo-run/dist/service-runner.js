// Generated by LiveScript 1.5.0
var child_process, watch, dashify, DockerRunner, EventEmitter, ref$, callArgs, DockerHelper, fs, Handlebars, yaml, N, os, ObservableProcess, path, portReservation, last, mkdir, ServiceRunner;
child_process = require('child_process');
watch = require('chokidar').watch;
dashify = require('dashify');
DockerRunner = require('./docker-runner');
EventEmitter = require('events').EventEmitter;
ref$ = require('../../exosphere-shared'), callArgs = ref$.callArgs, DockerHelper = ref$.DockerHelper;
fs = require('fs');
Handlebars = require('handlebars');
yaml = require('js-yaml');
N = require('nitroglycerin');
os = require('os');
ObservableProcess = require('observable-process');
path = require('path');
portReservation = require('port-reservation');
last = require('prelude-ls').last;
mkdir = require('shelljs').mkdir;
ServiceRunner = (function(superclass){
  var prototype = extend$((import$(ServiceRunner, superclass).displayName = 'ServiceRunner', ServiceRunner), superclass).prototype, constructor = ServiceRunner;
  function ServiceRunner(arg$){
    this.role = arg$.role, this.config = arg$.config, this.logger = arg$.logger;
    this.write = bind$(this, 'write', prototype);
    this.serviceConfigurationContent = bind$(this, 'serviceConfigurationContent', prototype);
    this.compileServiceDependencies = bind$(this, 'compileServiceDependencies', prototype);
    this.start = bind$(this, 'start', prototype);
    this.serviceConfig = yaml.safeLoad(this.serviceConfigurationContent());
  }
  ServiceRunner.prototype.start = function(done){
    var ref$, x$, y$, this$ = this;
    this.dockerConfig = {
      author: this.serviceConfig.author,
      image: dashify(this.serviceConfig.type),
      appName: dashify(this.config.appName),
      startCommand: this.serviceConfig.startup.command,
      startText: this.serviceConfig.startup['online-text'],
      cwd: this.config.root,
      env: {
        EXOCOM_PORT: this.config.EXOCOM_PORT,
        ROLE: this.role
      },
      publish: (ref$ = this.serviceConfig.docker) != null ? ref$.publish : void 8,
      dependencies: this.compileServiceDependencies()
    };
    x$ = this.dockerRunner = new DockerRunner({
      role: this.role,
      dockerConfig: this.dockerConfig,
      logger: this.logger
    });
    x$.startService();
    x$.on('online', function(){
      return typeof done == 'function' ? done() : void 8;
    });
    x$.on('error', function(message){
      return this$.emit('error', {
        errorMessage: message
      });
    });
    /* Ignores any sub-path including dotfiles.
    '[\/\\]' accounts for both windows and unix systems, the '\.' matches a single '.', and the final '.' matches any character. */
    y$ = this.watcher = watch(this.config.root, {
      ignoreInitial: true,
      ignored: [/.*\/node_modules\/.*/, /(^|[\/\\])\../]
    });
    y$.on('add', function(addedPath){
      this$.logger.log({
        role: 'exo-run',
        text: "Restarting service '" + this$.role + "' because " + addedPath + " was created"
      });
      return this$.restart();
    });
    y$.on('change', function(changedPath){
      this$.logger.log({
        role: 'exo-run',
        text: "Restarting service '" + this$.role + "' because " + changedPath + " was changed"
      });
      return this$.restart();
    });
    y$.on('unlink', function(removedPath){
      this$.logger.log({
        role: 'exo-run',
        text: "Restarting service '" + this$.role + "' because " + removedPath + " was deleted"
      });
      return this$.restart();
    });
    return y$;
  };
  ServiceRunner.prototype.restart = function(){
    var x$, this$ = this;
    this.dockerRunner.dockerContainer.kill();
    this.watcher.close();
    x$ = new ObservableProcess(callArgs(DockerHelper.getBuildCommand({
      author: this.dockerConfig.author,
      name: this.dockerConfig.image
    })), {
      cwd: this.config.root,
      stdout: {
        write: this.write
      },
      stderr: {
        write: this.write
      }
    });
    x$.on('ended', function(exitCode, killed){
      switch (false) {
      case exitCode !== 0:
        this$.logger.log({
          role: this$.role,
          text: "Docker image rebuilt"
        });
        return this$.start(function(){
          return this$.logger.log({
            role: 'exo-run',
            text: "'" + this$.role + "' restarted successfully"
          });
        });
      default:
        this$.logger.log({
          role: this$.role,
          text: "Docker image failed to rebuild"
        });
        return process.exit(exitCode);
      }
    });
    return x$;
  };
  ServiceRunner.prototype.compileServiceDependencies = function(){
    var dependencies, dependencyName, ref$, dependencyConfig, containerName, that, dataPath, volume, onlineText, port;
    dependencies = [];
    for (dependencyName in ref$ = this.serviceConfig.dependencies || {}) {
      dependencyConfig = ref$[dependencyName];
      containerName = this.config.appName + "-" + dependencyName;
      if ((that = dependencyConfig != null ? dependencyConfig.docker_flags : void 8) != null) {
        dataPath = path.join(os.homedir(), '.exosphere', this.config.appName, dependencyName, 'data');
        mkdir('-p', dataPath);
        volume = Handlebars.compile(that.volume)({
          "EXO_DATA_PATH": dataPath
        });
        onlineText = that.online_text;
        port = that.port;
      }
      dependencies.push({
        containerName: containerName,
        dependencyName: dependencyName,
        version: dependencyConfig != null ? dependencyConfig.version : void 8,
        volume: volume,
        onlineText: onlineText,
        port: port
      });
    }
    return dependencies;
  };
  ServiceRunner.prototype.serviceConfigurationContent = function(){
    switch (false) {
    case !this.config.image:
      return DockerHelper.getConfig(this.config.image);
    default:
      return fs.readFileSync(path.join(this.config.root, 'service.yml'));
    }
  };
  ServiceRunner.prototype.shutdownDependencies = function(){
    var i$, ref$, len$, dependency, results$ = [];
    for (i$ = 0, len$ = (ref$ = this.dockerConfig.dependencies).length; i$ < len$; ++i$) {
      dependency = ref$[i$];
      results$.push(DockerHelper.removeContainer(dependency.containerName + ""));
    }
    return results$;
  };
  ServiceRunner.prototype.write = function(text){
    return this.logger.log({
      role: this.role,
      text: text,
      trim: true
    });
  };
  return ServiceRunner;
}(EventEmitter));
module.exports = ServiceRunner;
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