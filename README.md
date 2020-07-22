# common

Minekube's Commons Library

**Get it:**
`go get -u go.minekube.com/common`

Libraries:
- minecraft/component
  - A Minecraft text components library.
  - Marshal/Unmarshal in different formats
    - json components (faster than Go's standard encoding/json,
    thanks to [Gojay](https://github.com/francoispqt/gojay))
    - plain text marshalling
    - legacy colors & formats
    - Minecraft 1.16+ hex colors
    - click/hover events
    - options for encoding/decoding
    - support older Minecraft client versions by default!