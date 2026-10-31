# Private panel profile

This fork keeps the upstream 3x-ui behavior by default and adds opt-in controls
for a private, operator-focused deployment.

## Runtime profile

Set the service environment variables without rebuilding the panel:

```env
XUI_BRAND_NAME=Private Proxy
XUI_ENABLED_PROTOCOLS=vless,hysteria,tuic,masque
```

`XUI_BRAND_NAME` changes the login and navigation identity. An empty value keeps
the upstream `3X-UI` name.

`XUI_ENABLED_PROTOCOLS` is a comma-separated inbound creation allowlist. An
empty value or `all` keeps every upstream protocol available. Existing inbounds
remain visible and editable when their protocol is removed from the allowlist;
the restriction applies to new inbounds and is enforced by the backend as well
as the form.

## Upstream workflow

```bash
git fetch upstream
git switch custom/main
git rebase upstream/main
```

The repository remotes are intentionally separated:

- `origin`: `https://github.com/jememouse/3x-ui.git`
- `upstream`: `https://github.com/MHSanaei/3x-ui.git`

Production releases should be cut from an explicit stable upstream tag. Do not
deploy a moving upstream branch directly to a VPS.

## Deployment boundary

This branch is not deployed automatically. Build and test it in parallel with
the current panel, using a separate database, panel port, and service name. A
production migration requires a verified backup and rollback path.
