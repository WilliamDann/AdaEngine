# ada-server

A small HTTP service that exposes the AdaEngine chess search as a stateless API, designed to run on Cloud Run. The frontend (`pf`) lives elsewhere and calls this service cross-origin.

## API

### `POST /api/move`

Request:
```json
{ "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "movetimeMs": 1000 }
```

Response:
```json
{ "move": "e2e4", "score": 30, "depth": 18, "nodes": 51234 }
```

The client owns game state and sends the current FEN on every move. The server parses it, searches, and returns the best move in UCI format (`e2e4`, `e7e8q` for promotions). `movetimeMs` is capped at 5000; out-of-range values fall back to the cap.

### `GET /api/start`
Returns `{"fen": "<start position FEN>"}`.

### `GET /api/health`
Returns `ok`. Use as a Cloud Run liveness/readiness probe.

## Run locally

From the AdaEngine repo root:

```
go run ./ada-server
```

Then:
```
curl -s localhost:8080/api/move -d '{"fen":"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1","movetimeMs":500}'
```

## Deploy to Cloud Run

From the AdaEngine repo root (the Dockerfile expects the full module as build context):

```
gcloud run deploy ada-engine \
  --source . \
  --region us-west1 \
  --port 8080 \
  --cpu 1 \
  --memory 512Mi \
  --min-instances 0 \
  --max-instances 2 \
  --timeout 60 \
  --allow-unauthenticated \
  --set-env-vars CORS_ORIGIN=https://your-frontend.example
```

Notes on the knobs:

- **`--min-instances 0`**: scales to zero between visitors. ~2-3s cold start, fine for a personal site.
- **`--max-instances 2`**: the engine is CPU-bound; this caps surprise cost. Each instance processes one search at a time (server-side semaphore).
- **`--cpu 1`**: one vCPU is plenty. Bump to 2 or 4 for stronger play (Lazy SMP uses all cores).
- **`--memory 512Mi`**: the transposition table is `1<<22` entries, comfortably under this.
- **`CORS_ORIGIN`**: set to your frontend's origin in production. Defaults to `*` for local dev.

## Architecture

Stateless by design. No game state lives on the server, so cold starts and instance churn don't break games in progress — the client just sends the current FEN. This sidesteps session affinity, Redis, and a database entirely.

One search at a time per instance (`searchSemaphore`), since each `search.Search` call already fans out across all cores via Lazy SMP.
