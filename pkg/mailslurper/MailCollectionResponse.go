// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
A MailCollectionResponse is sent in response to getting a collection
of mail items.
*/
type MailCollectionResponse struct {
	MailItems    []*MailItem `json:"mailItems"`
	TotalPages   int         `json:"totalPages"`
	TotalRecords int         `json:"totalRecords"`
}
