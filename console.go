package main

import (
	"fmt"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/sevenc-nanashi/pjsekai-overlay/pkg/pjsekaioverlay"
)

func Title() {
	fmt.Printf(
		strings.TrimSpace(dedent.Dedent(`
      %s------------------------------- pjsekai-overlay -------------------------------%s
        %sold pjsekai-overlay i think%s
        Version: %s%s%s (based on version %s%s%s)
        Original developed by %s名無し｡(@sevenc-nanashi)%s
        old version fixed by %sbnarpbnarp%s
        https://github.com/sevenc-nanashi/pjsekai-overlay
      %s-------------------------------------------------------------------------------%s
    `)) + "\n\n",
		RgbColorEscape(0x00afc7), ResetEscape(),
		RgbColorEscape(0x00afc7), ResetEscape(),
		RgbColorEscape(0x0f6ea3), pjsekaioverlay.Version, ResetEscape(),
                RgbColorEscape(0x614475), pjsekaioverlay.BaseVersion, ResetEscape(),
		RgbColorEscape(0x48b0d5), ResetEscape(),
                RgbColorEscape(0x5875a3), ResetEscape(),
		RgbColorEscape(0xff5a91), ResetEscape(),
	)

}

func RgbColorEscape(rgb int) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", (rgb>>16)&0xff, (rgb>>8)&0xff, rgb&0xff)
}

func AnsiColorEscape(color int) string {
  return fmt.Sprintf("\033[38;5;%dm", color)
}

func ResetEscape() string {
	return "\033[0m"
}

