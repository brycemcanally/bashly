package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/bryce/bashly/boxes"
	"github.com/bryce/bashly/config"
	"github.com/jroimartin/gocui"
)

func main() {
	cfgFile := flag.String("config", "", "Custom configuration file")
	scriptFile := flag.String("script", "", "Script to read/modify")
	flag.Parse()

	var cfg *config.Config
	if *cfgFile != "" {
		var err error
		cfg, err = config.Load(*cfgFile)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		cfg = config.Default()
	}

	if cfg.LogsDirectory != "" {
		file := setupLogging(cfg.LogsDirectory)
		defer file.Close()
	}

	boxs, err := boxes.New(cfg.Boxes)
	if err != nil {
		log.Panicln(err)
	}

	gui := setupGocui(boxs)
	defer gui.Close()

	if *scriptFile != "" {
		setupScript(gui, boxs, *scriptFile)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func setupLogging(logsDir string) *os.File {
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err := os.Mkdir(logsDir, os.ModePerm); err != nil {
			log.Panicln(err)
		}
	}

	file, err := os.OpenFile(logsDir+"/log", os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Panicln(err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(file)

	return file
}

func setupGocui(boxs *boxes.Boxes) *gocui.Gui {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	gui.Highlight = true
	gui.Cursor = true
	gui.SelFgColor = gocui.ColorGreen
	gui.Mouse = true

	layoutFunc := layout(boxs)
	gui.SetManagerFunc(layoutFunc)
	setupGlobalKeybindings(gui, boxs)
	if err := boxs.Setup(gui); err != nil {
		log.Panicln(err)
	}
	/* Setting the views manually to preload
	   the boxes so that I can load data outside of
	   the main loop */
	if err := boxs.SetViews(gui); err != nil {
		log.Panicln(err)
	}

	return gui
}

func setupGlobalKeybindings(gui *gocui.Gui, boxs *boxes.Boxes) {
	if err := gui.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, nextBox(boxs)); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, previousBox(boxs)); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
}

func setupScript(gui *gocui.Gui, boxs *boxes.Boxes, scriptFile string) {
	file, err := os.OpenFile(scriptFile, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicln(err)
	}

	scriptView, err := gui.View(boxs.Current().Name())
	if err != nil {
		log.Panicln(err)
	}
	if _, err := scriptView.Write(bytes); err != nil {
		log.Panicln(err)
	}

	// Function that will save changes to the script
	save := func(gui *gocui.Gui, view *gocui.View) error {
		if err := ioutil.WriteFile(scriptFile, []byte(scriptView.Buffer()), os.ModePerm); err != nil {
			log.Panicln(err)
		}
		return quit(gui, view)
	}

	if err := gui.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, save); err != nil {
		log.Panicln(err)
	}
}

func layout(boxs *boxes.Boxes) func(gui *gocui.Gui) error {
	return func(gui *gocui.Gui) error {
		if err := boxs.SetViews(gui); err != nil {
			return err
		}
		if err := boxs.Update(gui); err != nil {
			return err
		}

		return nil
	}
}

func nextBox(boxs *boxes.Boxes) func(gui *gocui.Gui, _ *gocui.View) error {
	return func(gui *gocui.Gui, _ *gocui.View) error {
		boxs.Next(gui)
		return nil
	}
}

func previousBox(boxs *boxes.Boxes) func(gui *gocui.Gui, _ *gocui.View) error {
	return func(gui *gocui.Gui, _ *gocui.View) error {
		boxs.Previous(gui)
		return nil
	}
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
