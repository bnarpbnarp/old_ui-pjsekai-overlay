package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/sevenc-nanashi/pjsekai-overlay/pkg/pjsekaioverlay"
	"github.com/srinathh/gokilo/rawmode"
	"golang.org/x/sys/windows"
)

func origMain(isOptionSpecified bool) {
	Title()

	var skipAviutlInstall bool
	flag.BoolVar(&skipAviutlInstall, "no-aviutl-install", false, "Skip installation of AviUtl object.")

	var outDir string
	flag.StringVar(&outDir, "out-dir", "./chart folder/_chartId_", "Specify the output directory. _chartId_ will be replaced with the chart ID.")

	var teamPower int
	flag.IntVar(&teamPower, "team-power", 250000, "Specify the team power.")

	var apCombo bool
	flag.BoolVar(&apCombo, "ap-combo", true, "Enables AP for combos.")

	flag.Usage = func() {
		fmt.Println("Usage: pjsekai-overlay [ChartID] [Option]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if !skipAviutlInstall {
		success := pjsekaioverlay.TryInstallObject()
		if success {
			fmt.Println("AviUtl objects installed yaaay")
		}
	}

	var chartId string
	if flag.Arg(0) != "" {
		chartId = flag.Arg(0)
		fmt.Printf("chartID: %s\n", color.GreenString(chartId))
	} else {
		fmt.Print("Chart ID (chcy- only) \ntype it here -> ")
		fmt.Scanln(&chartId)
		fmt.Printf("\033[A\033[2K\rok %s\n", color.GreenString(chartId))
	}

	chartSource, err := pjsekaioverlay.DetectChartSource(chartId)
	if err != nil {
		fmt.Println(color.RedString("server not found"))
		return
	}
	fmt.Printf("%s%s%s, getting chart.. ", RgbColorEscape(chartSource.Color), chartSource.Name, ResetEscape())
	chart, err := pjsekaioverlay.FetchChart(chartSource, chartId)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}
	if chart.Engine.Version != 12 {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: This engine is not supported. （version: %d）", chart.Engine.Version)))
		return
	}

	fmt.Println(color.GreenString("Done"))
	fmt.Printf("  %s / %s - %s (Lv. %s)\n",
		color.CyanString(chart.Title),
		color.CyanString(chart.Artists),
		color.CyanString(chart.Author),
		color.MagentaString(strconv.Itoa(chart.Rating)),
	)

	fmt.Printf("Finding exe file path.. ")
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("Done"))

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}

	formattedOutDir := filepath.Join(cwd, strings.Replace(outDir, "_chartId_", chartId, -1))
	fmt.Printf("chart path here idit: %s\n", color.CyanString(filepath.Dir(formattedOutDir)))

	fmt.Print("Downloading song image.. ")
	err = pjsekaioverlay.DownloadCover(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("Done"))

	fmt.Print("Downloading background.. (sadly you are going to have to get the v1 background yourself)")
	err = pjsekaioverlay.DownloadBackground(chartSource, chart, formattedOutDir)
	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("Done"))

	fmt.Print("Analyzing chart.. ")
	levelData, err := pjsekaioverlay.FetchLevelData(chartSource, chart)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("failed!: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("Done"))

	if !isOptionSpecified {
		fmt.Print("Team power: (e.g.: 300500; 393939; 324153) \ntype it here -> ")
		var tmpTeamPower string
		fmt.Scanln(&tmpTeamPower)
		teamPower, err = strconv.Atoi(tmpTeamPower)
		if err != nil {
			fmt.Println(color.RedString(fmt.Sprintf("failed!: %s", err.Error())))
			return
		}
		fmt.Printf("\033[A\033[2K\rok %s\n", color.GreenString(tmpTeamPower))

	}

	fmt.Print("Calculating score.. ")
	scoreData := pjsekaioverlay.CalculateScore(chart, levelData, teamPower)

	fmt.Println(color.GreenString("Done"))

	if !isOptionSpecified {
		fmt.Print("Enable ap combo? (Y/n)\n-> ")
		before, _ := rawmode.Enable()
		tmpEnableComboApByte, _ := bufio.NewReader(os.Stdin).ReadByte()
		tmpEnableComboAp := string(tmpEnableComboApByte)
		rawmode.Restore(before)
		fmt.Printf("\n\033[A\033[2K\rok %s\n", color.GreenString(tmpEnableComboAp))
		if tmpEnableComboAp == "Y" || tmpEnableComboAp == "y" || tmpEnableComboAp == "" {
			apCombo = true
		} else {
			apCombo = false
		}
	}
	executableDir := filepath.Dir(executablePath)
	assets := filepath.Join(executableDir, "assets")

	fmt.Print("Generating ped data file.. ")

	err = pjsekaioverlay.WritePedFile(scoreData, assets, apCombo, filepath.Join(formattedOutDir, "data.ped"))

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("Done"))

	fmt.Print("Generating exo file.. ")

	err = pjsekaioverlay.WriteExoFiles(assets, formattedOutDir)

	if err != nil {
		fmt.Println(color.RedString(fmt.Sprintf("Failed!: %s", err.Error())))
		return
	}

	fmt.Println(color.GreenString("Done"))

	fmt.Println(color.GreenString("\nYour overlay exo file was generated, try finding a tutorial if you don't know how to import it into AviUtl!"))
}

func main() {
	isOptionSpecified := len(os.Args) > 1
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	origMain(isOptionSpecified)

	if !isOptionSpecified {
		fmt.Print(color.CyanString("\nPress any key to exit.."))

		before, _ := rawmode.Enable()
		bufio.NewReader(os.Stdin).ReadByte()
		rawmode.Restore(before)
	}
}
