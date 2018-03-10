// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
PruneOption represents an option for pruning to be displayed in
an application.
*/
type PruneOption struct {
	PruneCode   PruneCode `json:"pruneCode"`
	Description string    `json:"description"`
}
