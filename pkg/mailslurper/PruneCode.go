// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package mailslurper

import "time"

/*
A PruneCode is a text code that represents a set of date ranges
*/
type PruneCode string

/*
ConvertToDate converts a prune code to a start
date based on the current date (in UTC). This function is a bit hard-coded
and weak based on the fact that the actual valid values are
defined in PruneOptions. It will have to do for now.
*/
func (pc PruneCode) ConvertToDate() string {
	now := time.Now().UTC()
	startDate := ""
	dateFormat := "2006-01-02"

	if pc.String() == "60plus" {
		startDate = now.AddDate(0, 0, -60).Format(dateFormat)
	} else if pc.String() == "30plus" {
		startDate = now.AddDate(0, 0, -30).Format(dateFormat)
	} else if pc.String() == "2wksplus" {
		startDate = now.AddDate(0, 0, -14).Format(dateFormat)
	}

	return startDate
}

/*
IsValid returns true/false if the prune code is a valid
option.
*/
func (pc PruneCode) IsValid() bool {
	result := false

	for _, option := range PruneOptions {
		if option.PruneCode == pc {
			result = true
			break
		}
	}

	return result
}

/*
String convers a PruneCode to string
*/
func (pc PruneCode) String() string {
	return string(pc)
}
