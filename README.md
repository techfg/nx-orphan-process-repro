# Nx Issue #32438: Orphaned Processes Reproduction

Minimal reproduction for orphaned child processes when stopping continuous tasks with Ctrl+C.

## Prerequisites

1. **Go v1.26.1** - <https://go.dev/doc/install>
2. macOS - Problem does not happen on Windows/WSL/Ubuntu (unable to test pure linux environment)

## Reproduction Steps

```bash
# 1. Clone repo
git clone https://github.com/techfg/nx-orphan-process-repro

# 2. Change directory
cd nx-orphan-process-repro

# 3. Install dependencies
npm install

# 4. Run the server
npm run start

# 5. Wait for the server to start, you'll see the following messages:
2026/03/11 02:31:06 INFO Starting platform server
2026/03/11 02:31:06 INFO Service started on port: 3000

# 6. Press Ctrl+C to stop

# 7. Check for orphaned processes
ps aux | grep -E 'platform[: ](dev|serve)'
```

## Expected result

No orphaned processes

## Actual Result: 

`platform:serve` is still running

```
barry       2665  8.3  0.5 3409416 94800 pts/10  Sl+  02:50   0:02 node /home/barry/repos/nx-orphan-process-repro/node_modules/.bin/nx run platform:dev
barry       2846  0.0  0.0 1857264 7840 pts/9    Sl+  02:50   0:00 ../../dist/platform/platform serve
```
