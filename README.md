# ory-selfservice-ui

[![ci status](https://github.com/twihike/ory-selfservice-ui/workflows/ci/badge.svg)](https://github.com/twihike/ory-selfservice-ui/actions) [![license](https://img.shields.io/github/license/twihike/ory-selfservice-ui)](LICENSE)

Self-service UI for Ory Kratos and Ory Hydra.

## Usage

Prerequisites:

- Kratos: <https://www.ory.sh/kratos/docs/install>
- Hydra: <https://www.ory.sh/hydra/docs/install>
- MailSlurper: <https://github.com/mailslurper/mailslurper/releases>

Start the processes.

```shell
./ctl.sh run_kratos
./ctl.sh run_hydra
./ctl.sh run_ss
./ctl.sh run_mailslurper
```

Access the UI.

- <http://127.0.0.1:4455/auth/dashboard>
- <http://127.0.0.1:4455/auth/registration>
- <http://127.0.0.1:4455/auth/login>

Try OAuth 2.0 flow.

```shell
./ctl.sh create_client
./ctl.sh auth
```

## License

Copyright (c) 2021 twihike. All rights reserved.

This project is licensed under the terms of the MIT license.
