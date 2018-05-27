FROM denismakogon/gocv-alpine:3.4.1-buildstage as build-stage

ENV GOPATH /go

RUN go get -u -d gocv.io/x/gocv

COPY main.go $GOPATH/src/github.com/justmiles/slack-webcam/

RUN go install github.com/justmiles/slack-webcam

FROM denismakogon/gocv-alpine:3.4.1-runtime

COPY --from=build-stage /go/bin/slack-webcam /slack-webcam

ENTRYPOINT ["/slack-webcam"]