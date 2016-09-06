@echo off
node_modules\o-tools-livescript\bin\build
node_modules\o-tools\bin\lint
node_modules\.bin\cucumber-js --strict
