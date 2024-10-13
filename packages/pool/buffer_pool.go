package pool

import (
	"bytes"
	"sync"
)

// Global buffer pool for request handling
var BufPool = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
