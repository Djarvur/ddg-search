## Open Questions

1. Should the `search` tool have a `source` parameter to explicitly choose between DDG and Perplexity, or should it always auto-select?
   - **Decision:** Support explicit source selection via `source` parameter with "auto" as default

2. What should be the default output format for the `search` tool?
   - **Decision:** Configurable via config file, with per-request override support

3. Should the `fetch` tool support any additional parameters from the CLI?
   - **Decision:** Support `timeout` and `user_agent` parameters, matching CLI options

4. Should the config file support multiple API keys?
   - **Decision:** Single API key for now, can be extended later if needed

5. Should there be a health check endpoint for TCP transport?
   - **Decision:** Yes, `/health` endpoint returning `{"status": "ok"}`

6. How should the server handle SIGINT (ctrl-c)?
   - **Decision:** Immediate shutdown - close all connections immediately

7. What should happen when SIGHUP is received?
   - **Decision:** Reload all configuration from the config file
