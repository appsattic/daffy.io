__DAFFY_NAKED_DOMAIN__ {
  proxy / localhost:__DAFFY_PORT__ {
    transparent
  }
  tls chilts@appsattic.com
  log stdout
  errors stderr
}

www.__DAFFY_NAKED_DOMAIN__ {
  redir http://__DAFFY_NAKED_DOMAIN__{uri} 302
}
