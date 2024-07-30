package app

import "BookSmart/pkg/logging"

func Run() {

	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}

	logger.Info("Starting server")
}
