FROM mcr.microsoft.com/devcontainers/go:0-1.19-bullseye
ADD . .
RUN bash ./install-zig.sh
RUN bash ./install-protoc.sh
RUN go install github.com/protocolbuffers/protoscope/cmd/protoscope...@latest
ENV AGREE=true
