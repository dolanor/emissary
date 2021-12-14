package service

import (
	"html/template"

	"github.com/benpate/derp"
	"github.com/benpate/ghost/model"
	"github.com/benpate/list"
	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
)

// Layout service manages the global site layout that is stored in a particular path of the
// filesystem.
type Layout struct {
	path    string
	funcMap template.FuncMap
	domain  model.Layout
	global  model.Layout
	group   model.Layout
	user    model.Layout
}

// NewLayout returns a fully initialized Layout service.
func NewLayout(path string, funcMap template.FuncMap) *Layout {

	service := &Layout{
		path:    path,
		funcMap: funcMap,
	}

	files := service.fileNames()

	// Use a separate counter because not all files will be included in the result
	for _, filename := range files {

		// Add all other directories into the Template service as Templates
		if err := service.loadFromFilesystem(filename); err != nil {
			derp.Report(derp.Wrap(err, "ghost.service.layout.NewLayout", "Error loading Layout from Filesystem"))
			panic("Error loading Layout from Filesystem")
		}
	}

	go service.watch()

	return service
}

/*******************************************
 * LAYOUT ACCESSORS
 *******************************************/

func (service *Layout) Global() model.Layout {
	return service.global
}

func (service *Layout) Group() model.Layout {
	return service.group
}

func (service *Layout) Domain() model.Layout {
	return service.domain
}

func (service *Layout) User() model.Layout {
	return service.user
}

/*******************************************
 * FILE WATCHER
 *******************************************/

// fileNames returns a list of directories that are owned by the Layout service.
func (service *Layout) fileNames() []string {
	return []string{"global", "user"}
}

// watch must be run as a goroutine, and constantly monitors the
// "Updates" channel for news that a template has been updated.
func (service *Layout) watch() {

	// Create a new directory watcher
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		panic(err)
	}

	files := service.fileNames()

	// Use a separate counter because not all files will be included in the result
	for _, filename := range files {

		if err := watcher.Add(service.path + "/" + filename); err != nil {
			derp.Report(derp.Wrap(err, "ghost.service.Layout.watch", "Error adding file watcher to file", filename))
		}
	}

	// Repeat indefinitely, listen and process file updates
	for {

		select {

		case event, ok := <-watcher.Events:

			if !ok {
				continue
			}

			filename := list.Last(list.RemoveLast(event.Name, "/"), "/")

			if err := service.loadFromFilesystem(filename); err != nil {
				derp.Report(derp.Wrap(err, "ghost.service.Layout.watch", "Error loading changes to layout", event, filename))
				continue
			}

		case err, ok := <-watcher.Errors:

			if ok {
				derp.Report(derp.Wrap(err, "ghost.service.Layout.watch", "Error watching filesystem"))
			}
		}
	}
}

// loadFromFilesystem retrieves the template from the disk and parses it into
func (service *Layout) loadFromFilesystem(filename string) error {

	path := service.path + "/" + filename
	layout := model.NewLayout(filename, service.funcMap)

	spew.Dump("loadFromFilesystem ---", filename, path)

	// System folders (except for "static" and "global") have a schema.json file
	if (filename != "static") && (filename != "global") {
		if err := loadModelFromFilesystem(path, &layout); err != nil {
			return derp.Wrap(err, "ghost.service.layout.getTemplateFromFilesystem", "Error loading Schema", filename)
		}
	}

	if err := loadHTMLTemplateFromFilesystem(path, layout.HTMLTemplate, service.funcMap); err != nil {
		return derp.Wrap(err, "ghost.service.layout.getTemplateFromFilesystem", "Error loading Schema", filename)
	}

	switch filename {

	case "global":
		spew.Dump("updated global")
		service.global = layout
	case "group":
		spew.Dump("updated group")
		service.group = layout
	case "user":
		spew.Dump("updated user")
		service.user = layout
	}

	return nil
}
