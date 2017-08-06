# Copyright 2017 aerth. All rights reserved.
# Use of this source code is governed by a GPL-style
# license that can be found in the LICENSE file.

# build helper
build:
	GOBIN=${PWD} go install -v .

install:
	install helper /usr/local/bin/

package:
	gox && \
	for i in $(ls helper_*); do zip $i.zip $i README.md LICENSE main.go makefile; done

clean:
	rm -vf helper*

