// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
DeleteMailRequest is used when requesting to delete mail
items.
*/
type DeleteMailRequest struct {
	PruneCode PruneCode `json:"pruneCode" form:"pruneCoe"`
}
