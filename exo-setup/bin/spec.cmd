call node_modules\o-tools-livescript\bin\build
call node_modules\o-tools\bin\lint
node_modules\.bin\cucumber-js %*
