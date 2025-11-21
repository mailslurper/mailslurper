// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
package mailslurper

import (
	"github.com/microcosm-cc/bluemonday"
)

type XSSService struct {
	sanitizer *bluemonday.Policy
}

/*
NewXSSService creates a new cross-site scripting service with default policy and with additional
allowed attributes on elements from config file.
*/
func NewXSSService(config *Configuration) *XSSService {
	policy := bluemonday.UGCPolicy()
	for _, attributeSets := range config.AllowedHTMLTags {
		policy.AllowAttrs(attributeSets["attributes"]...).OnElements(attributeSets["elements"]...)
	}

	return &XSSService{
		sanitizer: policy,
	}
}

/*
SanitizeString attempts to sanitize a string by removing potentially dangerous
HTML/JS markup.
*/
func (service *XSSService) SanitizeString(input string) string {
	return service.sanitizer.Sanitize(input)
}
