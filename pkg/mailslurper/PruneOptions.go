// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

/*
PruneOptions is the master list of options given to a user
when pruning mail from the database
*/
var PruneOptions = []PruneOption{
	{PruneCode("60plus"), "Older than 60 days"},
	{PruneCode("30plus"), "Older than 30 days"},
	{PruneCode("2wksplus"), "Older than 2 weeks"},
	{PruneCode("all"), "All emails"},
}
