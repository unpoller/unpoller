# Saving UniFi API Output

Ways to save API responses and explorer output for discovery or debugging.

## Single endpoint → file

Redirect `unpoller -j "other <path>"` to a file:

```bash
unpoller -c up.conf -j "other /api/s/default/stat/device" > device.json
unpoller -c up.conf -j "other /api/s/default/stat/sta"     > clients.json
```

Use `jq` to inspect: `jq . device.json`

## Bulk dump → directory

Use the dump script to request many known endpoints and save each to a JSON file:

```bash
./scripts/dump_unifi_api.sh -c up.conf -s default -o ./api_dump
```

Output goes to `./api_dump` by default. See `./scripts/dump_unifi_api.sh -h` for options.

Note: some endpoints (e.g. `sitedpi`, `stadpi`) require POST with a body; the script only issues GETs, so those may fail or return errors. You can still inspect the saved responses.

## Saving the API explorer UI

If you're using the developer UI (e.g. [developer.ui.com](https://developer.ui.com) or another API explorer) and want to save the **list of endpoints and their details**:

1. **OpenAPI / Swagger spec**  
   Open DevTools → **Network**, (re)load the explorer, and look for requests to `openapi.json`, `swagger.json`, or similar. Right‑click the response → **Copy** → **Save as**, or use **Save all as HAR** and extract the spec from the HAR.

2. **Save page**  
   Use **File → Save As** (HTML) or **Print → Save as PDF** to capture the visible explorer structure. This won’t persist dynamically loaded data unless the page embeds it.

3. **Export**  
   If the explorer has an **Export** or **Download** button (e.g. for OpenAPI YAML/JSON), use that to save the full spec.

4. **Community specs**  
   Community OpenAPI specs for the UniFi API exist (e.g. [ubiquiti-community/unifi-api](https://github.com/ubiquiti-community/unifi-api), [ringods/unifi-api-spec](https://github.com/ringods/unifi-api-spec)). Clone or download those repos to get machine‑readable API definitions.
