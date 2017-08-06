# helper

Command Line Assistant for people who [bang](https://duckduckgo.com/bang)[!](https://duckduckgo.com/lite?q=!g+"aerth"+helper+github)[!](https://duckduckgo.com/bang?c=Tech&sc=Languages+(Go))

### Usage

```
helper v0.0.3
https://github.com/aerth/helper

Available Resources:
         g https://encrypted.google.com/search?q=%s
         ddg https://duckduckgo.com/lite?q=%s

Available Bookmarks:
         gonew https://github.com/search?o=desc&q=language%3Ago+stars%3A1&repo=&s=updated&start_value=1&type=Repositories
         gotrend https://github.com/trending/go
         news https://news.ycombinator.com/
         r https://reddit.com/r/golang/

Flags:
  -x    dont open links, just print them
```

### Copyright

Copyright 2017 aerth. All rights reserved.

Use of this source code is governed by a GPL-style

license that can be found in the [LICENSE](LICENSE.md) file.

### Contributing

Bugs/Issues addressed at [Github](https://github.com/aerth/helper/issues)

Contributions welcome. Please run `gofmt -w -l -s .` before your commit, Thanks!

### Known Issues

With bash, it is **highly recommended** to `set +H` to disable `!`-style history expansion.
If you actually use bash `!-3` to repeat the third most recent command, just single quote your bangs:

```
helper '!gh language:go stars:<5'
```

Other shells don't seem to have this issue.

You can add this to your `.bash_profile` or `.bashrc`:

```
# disable !-style history expansion for helper bangs (https://github.com/aerth/helper)
set +H
```

