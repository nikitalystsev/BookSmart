package app

import "BookSmart/pkg/logging"

func Run() {

	logger, err := logging.NewLogger()
	if err != nil {
		panic(err)
	}

	logger.Info("Starting server")
}

/*
docker run \
 --name influxdb2 \
 --publish 8086:8086 \
 --env DOCKER_INFLUXDB_INIT_MODE=setup \
 --env DOCKER_INFLUXDB_INIT_USERNAME=nikitalystsev \
 --env DOCKER_INFLUXDB_INIT_PASSWORD=4CD-sr9-x4N-SjH \
 --env DOCKER_INFLUXDB_INIT_ORG=first_org \
 --env DOCKER_INFLUXDB_INIT_BUCKET=first_bucket \
 influxdb:2
*/
