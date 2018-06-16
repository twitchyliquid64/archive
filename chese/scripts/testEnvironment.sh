#!/bin/bash

go build -o chese-setup github.com/twitchyliquid64/chese/setup
go build -o chese-client github.com/twitchyliquid64/chese/chese
go build -o chese-server github.com/twitchyliquid64/chese/server

./chese-setup \
  --ca-cert ca-cert.pem\
  --ca-key ca-key.pem\
  --serv-cert serv-cert.pem\
  --serv-key serv-key.pem\
  --client-cert client-cert.pem\
  --client-key client-key.pem\
   mint-server

 ./chese-setup \
   --ca-cert ca-cert.pem\
   --ca-key ca-key.pem\
   --serv-cert serv-cert.pem\
   --serv-key serv-key.pem\
   --client-cert client-cert.pem\
   --client-key client-key.pem\
   --client-cert-validity 96h \
    issue-cert
