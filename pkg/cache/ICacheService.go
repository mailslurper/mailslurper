// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cache

import (
	"time"
)

type ICacheService interface {
	Delete(key string)
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, timeout time.Duration)
}
