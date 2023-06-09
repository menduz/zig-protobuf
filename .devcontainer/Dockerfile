FROM mcr.microsoft.com/devcontainers/go:0-1.19-bullseye
ADD . .
RUN bash ./install-zig.sh
RUN bash ./install-protoc.sh
ENV AGREE=true
