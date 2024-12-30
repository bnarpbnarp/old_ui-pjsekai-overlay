package pjsekaioverlay

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"strings"

	wapi "github.com/iamacarpet/go-win64api"
	so "github.com/iamacarpet/go-win64api/shared"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

//go:embed sekai.obj
var sekaiObj []byte

//go:embed sekai-en.obj
var sekaiObjEn []byte

func TryInstallObject() bool {
	processes, _ := wapi.ProcessList()
	var aviutlProcess *so.Process
	for _, process := range processes {
		if process.Executable == "aviutl.exe" {
			aviutlProcess = &process
			break
		}
	}
	if aviutlProcess == nil {
		return false
	}
	var aviutlPath string
	aviutlPath = filepath.Dir(aviutlProcess.Fullpath)
	var exeditRoot string
	if _, err := os.Stat(filepath.Join(aviutlPath, "exedit.auf")); err == nil {
		exeditRoot = filepath.Join(aviutlPath)
	} else if _, err := os.Stat(filepath.Join(aviutlPath, "Plugins", "exedit.auf")); err == nil {
		exeditRoot = filepath.Join(aviutlPath, "Plugins")
	} else {
		return false
	}

	os.MkdirAll(filepath.Join(exeditRoot, "script"), 0755)

	var sekaiObjPath = filepath.Join(exeditRoot, "script", "@pjsekai_ui.obj")
	if _, err := os.Stat(sekaiObjPath); err == nil {
		var sekaiObjFile, _ = os.OpenFile(sekaiObjPath, os.O_RDONLY, 0755)
		defer sekaiObjFile.Close()
		var sekaiObjDecoder = japanese.ShiftJIS.NewDecoder()
		var existingSekaiObj, _ = io.ReadAll(transform.NewReader(sekaiObjFile, sekaiObjDecoder))
		if strings.Contains(string(existingSekaiObj), "--version: "+Version) && Version != "0.0.0" {
			return false
		}
	}
	var sekaiObjPathEn = filepath.Join(exeditRoot, "script", "@pjsekai_ui_en.obj")
	if _, err := os.Stat(sekaiObjPathEn); err == nil {
		var sekaiObjFileEn, _ = os.OpenFile(sekaiObjPathEn, os.O_RDONLY, 0755)
		defer sekaiObjFileEn.Close()
		var sekaiObjDecoderEn = japanese.ShiftJIS.NewDecoder()
		var existingSekaiObjEn, _ = io.ReadAll(transform.NewReader(sekaiObjFileEn, sekaiObjDecoderEn))
		if strings.Contains(string(existingSekaiObjEn), "--version: "+Version) && Version != "0.0.0" {
			return false
		}
	}
	err := os.MkdirAll(filepath.Join(exeditRoot, "script"), 0755)
	if err != nil {
		return false
	}
	sekaiObjFile, err := os.Create(sekaiObjPath)
	if err != nil {
		return false
	}
	defer sekaiObjFile.Close()

	sekaiObjFileEn, err := os.Create(sekaiObjPathEn)
	if err != nil {
		return false
	}
	defer sekaiObjFileEn.Close()

	var sekaiObjWriter = transform.NewWriter(sekaiObjFile, japanese.ShiftJIS.NewEncoder())
	var sekaiObjWriterEn = transform.NewWriter(sekaiObjFileEn, japanese.ShiftJIS.NewEncoder())

	strings.NewReader(strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
		"{version}", Version,
	).Replace(string(sekaiObj))).WriteTo(sekaiObjWriter)

	strings.NewReader(strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
		"{version}", Version,
	).Replace(string(sekaiObjEn))).WriteTo(sekaiObjWriterEn)
	return true
}
