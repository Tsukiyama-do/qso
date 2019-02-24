package env

import (
  "log"
)

var smode bool = true

func S_mode() bool {
   return smode
}

func S_mset(b bool) {
    smode = b
    log.Printf("Mode is changed to %t. \n", b)
}
