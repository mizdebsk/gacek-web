FROM golang:1.22 as build
ENV CGO_ENABLED=0
COPY go.mod go.sum *.go *.html /
RUN ["go","build","-o","/app","/app.go","/data.go","/model.go","/states.go"]

FROM gcr.io/distroless/static:nonroot
COPY --from=build /app /*.html /
EXPOSE 8080
CMD ["/app"]
