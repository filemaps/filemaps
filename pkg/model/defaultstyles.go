// Copyright (c) 2017, CodeBoy. All rights reserved.
//
// This Source Code Form is subject to the terms of the
// license that can be found in the LICENSE file.

package model

func NewDefaultStyles() []Style {
	var styles []Style
	styles = append(styles, Style{
		SClass: "go",
		Rules: map[string]string{
			"color": "#375eab",
		},
	})
	styles = append(styles, Style{
		SClass: "html",
		Rules: map[string]string{
			"color": "#ff0000",
		},
	})
	styles = append(styles, Style{
		SClass: "md",
		Rules: map[string]string{
			"color": "#00ff00",
		},
	})
	styles = append(styles, Style{
		SClass: "ts",
		Rules: map[string]string{
			"color": "#0000ff",
		},
	})
	return styles
}
