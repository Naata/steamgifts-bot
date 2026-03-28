FROM ubuntu:24.04

RUN apt-get update -y && apt-get install -y ca-certificates

WORKDIR /workspace

COPY ./sg_bot /workspace/sg_bot
RUN chmod +x /workspace/sg_bot

CMD ["/workspace/sg_bot"]