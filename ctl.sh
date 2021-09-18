#!/bin/sh

SELF_DIR=$( (cd $(dirname $0); pwd) )

run_ss() {
  (
    source "${SELF_DIR}/configs/env.sh"
    cd "${SELF_DIR}"
    go run .
  )
}

run_kratos() {
  kratos serve -c "${SELF_DIR}/configs/kratos.yml" --watch-courier --sqa-opt-out
  # kratos courier watch -c "${SELF_DIR}/configs/kratos.yml"
}

run_hydra() {
  hydra serve all -c "${SELF_DIR}/configs/hydra.yml" --dangerous-force-http
}

run_mailslurper() {
  (
    cd "${SELF_DIR}/configs"
    mailslurper
  )
}

create_client() {
  hydra clients create \
    --endpoint http://127.0.0.1:4445 \
    --id auth-code-client \
    --secret secret \
    --response-types code,id_token \
    --grant-types authorization_code,refresh_token \
    --scope openid,offline_access,email \
    --callbacks http://127.0.0.1:5555/callback
}

list_client() {
  hydra clients list \
    --endpoint http://127.0.0.1:4445
}

delete_client() {
  hydra clients delete auth-code-client \
    --endpoint http://127.0.0.1:4445
}

auth() {
  hydra token user \
    --endpoint http://127.0.0.1:4444 \
    --client-id auth-code-client \
    --client-secret secret \
    --scope openid,offline_access,email \
    --port 5555
}

token() {
  hydra token introspect --endpoint http://127.0.0.1:4445 "$@"
}

"$@"
