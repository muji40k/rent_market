#! /bin/bash

vault policy write runner policy.hcl
vault write auth/jwt/role/runner @roles.json

