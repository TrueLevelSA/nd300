<!--
  Copyright 2022 TrueLevel SA

  This Source Code Form is subject to the terms of the Mozilla Public
  License, v. 2.0. If a copy of the MPL was not distributed with this
  file, You can obtain one at https://mozilla.org/MPL/2.0/.

  SPDX-License-Identifier: MPL-2.0
-->

# go.swapbox.cash/nd300 (ND-300KM Driver)

[![License](https://img.shields.io/badge/license-MPL--2.0-orange)](https://gitlab.com/TrueLevel/swapbox/nd300/-/blob/main/LICENSE)

Driver to communicate and control a [ND-300KM Note dispencer from ICT Corp.][product] 
over Serial, written in Go.

> ⚠️ **This library is in active development and not in a working state.** ⚠️
> 
> This software is distributed as is, without warranty.
> The authors are not liable for any claim, damage, or financial loss 
> related to the use of this library.

## Installation

```shell
go get -u go.swapbox.cash/nd300
```

## Development

### Code Generation 
For development purpose, install the following tools to be able to regenerate
the code:

- [stringer][]: `go install golang.org/x/tools/cmd/stringer@latest`

Then run: `go generate`.

### Tests

TODO 😞

[product]: http://www.ictgroup.tw/pro_cen.php?prod_id=70
[stringer]: https://pkg.go.dev/golang.org/x/tools/cmd/stringer
