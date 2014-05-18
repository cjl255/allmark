// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonhandler

import (
	"encoding/json"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/orchestrator"
	"github.com/andreaskoch/allmark2/ui/web/server/handler/errorhandler"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func New(logger logger.Logger, config *config.Config, itemIndex *index.Index, patherFactory paths.PatherFactory, converter conversion.Converter) *JsonHandler {

	// tags
	itemPathProvider := patherFactory.Absolute("/")
	tagPathProvider := patherFactory.Absolute("/tags.html#")
	tagsOrchestrator := orchestrator.NewTagsOrchestrator(itemIndex, tagPathProvider, itemPathProvider)

	// navigation
	navigationPathProvider := patherFactory.Absolute("/")
	navigationOrchestrator := orchestrator.NewNavigationOrchestrator(itemIndex, navigationPathProvider)

	// error
	error404Handler := errorhandler.New(logger, config, itemIndex, patherFactory)

	// viewmodel
	viewModelOrchestrator := orchestrator.NewViewModelOrchestrator(itemIndex, converter, &navigationOrchestrator, &tagsOrchestrator)

	return &JsonHandler{
		logger:                logger,
		itemIndex:             itemIndex,
		config:                config,
		patherFactory:         patherFactory,
		error404Handler:       error404Handler,
		viewModelOrchestrator: viewModelOrchestrator,
	}
}

type JsonHandler struct {
	logger                logger.Logger
	itemIndex             *index.Index
	config                *config.Config
	patherFactory         paths.PatherFactory
	error404Handler       *errorhandler.ErrorHandler
	viewModelOrchestrator orchestrator.ViewModelOrchestrator
}

func (handler *JsonHandler) Func() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the path from the request variables
		vars := mux.Vars(r)
		path := vars["path"]

		// get the request route
		requestRoute, err := route.NewFromRequest(path)
		if err != nil {
			handler.logger.Error("Unable to get route from request. Error: %s", err)
			return
		}

		// make sure the request body is closed
		defer r.Body.Close()

		// stage 1: check if there is a item for the request
		if item, found := handler.itemIndex.IsMatch(*requestRoute); found {

			// create the view model
			pathProvider := handler.patherFactory.Relative(item.Route())
			viewModel := handler.viewModelOrchestrator.GetViewModel(pathProvider, item)

			// return the json
			writeJson(w, viewModel)
			return
		}

		// display a 404 error page
		error404Handler := handler.error404Handler.Func()
		error404Handler(w, r)
	}
}

func writeJson(writer io.Writer, viewModel viewmodel.Model) error {
	bytes, err := json.MarshalIndent(viewModel, "", "\t")
	if err != nil {
		return err
	}

	writer.Write(bytes)
	return nil
}
