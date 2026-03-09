# Error Diagnostics Design

## Goal

Improve CLI error visibility so LLM callers can identify failure category, retry strategy, and probable remediation from stderr output.

## Scope

- Keep stdout behavior unchanged: only Markdown table on success.
- Keep process exit semantics unchanged: non-zero on failure.
- Enhance stderr only, using one-line key=value format.

## Chosen Approach

Use a centralized application error model in `internal/apperr`:

- Stable error `code`
- `retryable` boolean
- Human-readable `message`
- Action-oriented `hint`
- Raw `cause`
- Optional metadata (`server`, `query`, `timeout`, `status`)

All command-facing failures are normalized through this model before printing.

## Error Codes

- `NETWORK_DNS`
- `NETWORK_TIMEOUT`
- `NETWORK_RESET`
- `TLS_HANDSHAKE`
- `HTTP_NON_2XX`
- `RESPONSE_DECODE`
- `CONFIG_INVALID`
- `TEMPLATE_INVALID`
- `INTERNAL_UNKNOWN`

## Output Format

`ERROR code=<CODE> retryable=<true|false> message="..." hint="..." cause="..." [meta fields]`

Meta fields are appended only when available, for example `server`, `query`, `timeout`, `status`.

## Mapping Rules

Priority order:

1. timeout/context deadline -> `NETWORK_TIMEOUT`
2. TLS/x509/handshake errors -> `TLS_HANDSHAKE`
3. DNS resolution errors -> `NETWORK_DNS`
4. reset/EOF/broken pipe -> `NETWORK_RESET`
5. HTTP non-2xx -> `HTTP_NON_2XX`
6. JSON decode -> `RESPONSE_DECODE`
7. config/template failures -> `CONFIG_INVALID` / `TEMPLATE_INVALID`
8. fallback -> `INTERNAL_UNKNOWN`

## Verification

- Unit tests for classification and formatter output.
- Command tests for help text and rendered stderr line content.
- End-to-end smoke checks for invalid host and timeout-like errors.
