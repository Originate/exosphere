normalize-path = (file-path) ->
  | process.platform is 'win32'  =>  file-path.replace /\//g, '\\'
  | otherwise                    =>  file-path.replace /\\/g, '/'
