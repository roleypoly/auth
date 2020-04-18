## Frontend Login

covering users getting their authorization grants

### min reqs

- User must receive a session key
- Session key is representative of backend data with grant data
- Grant data is typically just user id, sources, and related metadata

### proto flow 1 (discord oauth)

0. UI stores desired next state IF relevant
0. user clicks "login with discord"
0. hands off to auth api to build the redirect URL w/ state
0. waits for Discord OAuth
0. discord oauth returns to api with state and code
0. auth api asks discord client for grant info
0. auth api returns to ui (probably next API route?) with temp auth code in ui auth flow
0. ui resolves session access key from auth api, stores (cookie? sessionstorage?)
0. UI resolves next desired state from storage

### proto flow 2 (bot dm auth)

0. UI stores desired next state IF login is accessed before DM
0. user initiates flow via DM
0. discord requests an auth challenge for a user from auth rpc
0. challenge is presented to user
0. user continues flow in ui (via magic link or human code)
0. either flow recieves temp auth code and moves to ui auth flow
0. ui resolves session access key from auth api, stores (cookie? sessionstorage?)
0. UI resolves next desired state from storage if exists, otherwise to default

### api breakdown

- RPCs
    - AuthClient.ResolveSessionKey(authToken) from UI
        - flow 1 final
        - flow 2 final
    - AuthInternal.GetNewAuthChallenge(userId) from Discord
        - flow 2 entry
    - AuthClient.AuthorizeChallenge(authChallenge) from UI
        - flow 2 authorize
- HTTP (browsers touch this only)
    - auth: /oauth/redirect?state=...
        - flow 1 entry
    - auth: /oauth/callback?code=...&state=...
        - flow 1 authorize
    - ui: /auth/machinery/authenticate?token=...&state=...
        - flow 1 authenticate
        - flow 2 authenticate