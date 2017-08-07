# Node.js Text Stream Search
> Searches for occurrences of a given search term in a Node.js text stream

[![Circle CI](https://circleci.com/gh/Originate/node-text-stream-search.svg?style=shield)](https://circleci.com/gh/Originate/node-text-stream-search)
[![Dependency Status](https://david-dm.org/originate/node-text-stream-search.svg)](https://david-dm.org/originate/node-text-stream-search)
[![devDependency Status](https://david-dm.org/originate/node-text-stream-search/dev-status.svg)](https://david-dm.org/originate/node-text-stream-search#info=devDependencies)
<a href="https://yarnpkg.com">
  <img src="https://img.shields.io/badge/yarn-compatible-brightgreen.svg">
</a>


Super simple and robust way to reliably recognize text in text streams.


```javascript

streamSearch = require('text-stream-search');

new TextStreamSearch(myStream).wait('hello', function() { console.log('found hello') });

// Now if the given stream emits someting including the word "hello",
// the callback above is called exactly once.
```

More details about the wait command is [here](features/wait.feature).

TextStreamSearch uses the [debug](https://github.com/visionmedia/debug)
library.
To see the parsed text, run your code that uses TextStreamSearch
with the environment variable `DEBUG` set to `text-stream-search` or `*`.


## Related Projects

* [StreamSnitch](https://github.com/dmotz/stream-snitch): does the same thing with regular expressions,
  but is buggy and blocks the event queue


## Development

See our [developer guidelines](CONTRIBUTING.md)
