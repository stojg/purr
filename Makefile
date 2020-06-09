run: build
	docker run -it --rm stojg/purr

build:
	docker build . -t stojg/purr

push: build
	docker build . -t stojg/purr:latest -t stojg/purr:$(shell git rev-parse --verify HEAD)
	docker push stojg/purr:latest
	docker push stojg/purr:$(shell git rev-parse --verify HEAD)
