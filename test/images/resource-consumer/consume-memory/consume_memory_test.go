/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"testing"
)

func TestConsumeMemoryHelpers(t *testing.T) {
	parsedSize, _ := parseSize("100M")
	expectedSize := int64(100 * sizeMB)
	if parsedSize != expectedSize {
		t.Errorf("Size parsing: got: %d, want %d", parsedSize, expectedSize)
	}

	_, err := parseSize("FooM")
	if err == nil {
		t.Error("Expected error while parsing 'FooM'")
	}

	byteCount := 100
	bytes := bigBytes(int64(byteCount))
	if len(*bytes) != byteCount {
		t.Errorf("Byte allocation size: got: %d, want: %d", len(*bytes), byteCount)
	}
}

func TestConsumeMemoryMain(t *testing.T) {
	*workers = 2
	*sizestr = "1M"
	errors := consume()
	if len(errors) != 0 {
		t.Errorf("Encountered unexpected errors: %v", errors)
	}
}
