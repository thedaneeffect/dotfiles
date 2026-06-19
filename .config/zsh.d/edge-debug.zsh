# Headed Edge on the Windows host for GPU-backed browser automation (Playwright MCP).
#
# WHY: WSL2 headless Chromium has no GPU-backed Vulkan ICD, so WebGPU apps (e.g. Loom
# on localhost:6969) die at boot with "No GPU adapter found". Driving the *host* Edge
# over CDP uses the real GPU. WSL mirrored networking makes localhost:PORT shared, so
# Playwright MCP (configured with --cdp-endpoint http://localhost:PORT) reaches it with
# no host-IP juggling or firewall holes.
#
# Isolation guarantees — your daily Edge is never touched:
#   --user-data-dir=<dedicated>  separate profile/cache/history/cookies/logins
#   --disable-extensions         no extensions beyond what you pass in EDGE_DEBUG_EXTRA_ARGS
#   --disable-sync               won't sign into / sync your MS account
#   --no-default-browser-check   no "make Edge default" nag
# It launches as its own window/instance, so it won't hijack the Edge you're using.
#
# Usage:
#   edge-debug          start (idempotent) and print the CDP endpoint
#   edge-debug-status   is it up?
#   edge-debug-stop     kill ONLY the debug-profile Edge (matched by its user-data-dir)
#   edge-debug-restart  stop + start
# Then, in Claude Code:  /mcp -> playwright -> Reconnect
#
# Config (override before calling, e.g. in a local zsh.d snippet or your env):
#   EDGE_DEBUG_PORT=9222
#   EDGE_DEBUG_PROFILE='C:\Users\you\AppData\Local\pw-edge-debug'   # default: %LOCALAPPDATA%\pw-edge-debug
#   EDGE_DEBUG_EXE='/mnt/c/.../msedge.exe'                          # default: auto-detected
#   EDGE_DEBUG_WINDOWSTYLE=Minimized                                # or Normal (live FPS work)
#   EDGE_DEBUG_EXTRA_ARGS=( --load-extension=C:\path\to\ext )       # explicit opt-in extras

# WSL-only: these helpers shell out to Windows (powershell.exe / msedge.exe).
if [[ -n "${WSL_DISTRO_NAME:-}" ]] || grep -qiE '(microsoft|wsl)' /proc/version 2>/dev/null; then

: ${EDGE_DEBUG_PORT:=9222}

# Linux path to msedge.exe (for existence check); honor override.
_edge_debug_exe() {
    emulate -L zsh
    if [[ -n "${EDGE_DEBUG_EXE:-}" ]]; then print -r -- "$EDGE_DEBUG_EXE"; return 0; fi
    local p
    for p in \
        '/mnt/c/Program Files (x86)/Microsoft/Edge/Application/msedge.exe' \
        '/mnt/c/Program Files/Microsoft/Edge/Application/msedge.exe'; do
        [[ -f "$p" ]] && { print -r -- "$p"; return 0; }
    done
    print -u2 "edge-debug: msedge.exe not found; set EDGE_DEBUG_EXE"; return 1
}

# Dedicated Windows profile path. Default lives under %LOCALAPPDATA% so it never
# collides with your real Edge profile and needs no admin to create.
_edge_debug_profile() {
    emulate -L zsh
    if [[ -n "${EDGE_DEBUG_PROFILE:-}" ]]; then print -r -- "$EDGE_DEBUG_PROFILE"; return 0; fi
    local lad; lad="$(cmd.exe /c 'echo %LOCALAPPDATA%' 2>/dev/null | tr -d '\r')"
    print -r -- "${lad}\\pw-edge-debug"
}

edge-debug-status() {
    emulate -L zsh
    if curl -fsS --max-time 3 "http://localhost:${EDGE_DEBUG_PORT}/json/version" >/dev/null 2>&1; then
        print "edge-debug: up -> http://localhost:${EDGE_DEBUG_PORT}"
        return 0
    fi
    print "edge-debug: down"
    return 1
}

edge-debug() {
    emulate -L zsh
    if curl -fsS --max-time 3 "http://localhost:${EDGE_DEBUG_PORT}/json/version" >/dev/null 2>&1; then
        print "edge-debug: already running -> http://localhost:${EDGE_DEBUG_PORT}"
        return 0
    fi

    local exe_lin exe_win profile
    exe_lin="$(_edge_debug_exe)"     || return 1
    exe_win="$(wslpath -w "$exe_lin")"
    profile="$(_edge_debug_profile)"

    # Only flags listed here are applied (per requirement). EDGE_DEBUG_EXTRA_ARGS is the
    # one explicit escape hatch for opt-in extras (e.g. --load-extension).
    #
    # The --disable-*backgrounding/occlusion/throttle flags keep WebGPU rendering at full
    # rate even when the window is unfocused or minimized — so you never have to keep it in
    # front. (Playwright screenshots/clicks go through CDP regardless of window state.)
    local -a flags
    flags=(
        "--remote-debugging-port=${EDGE_DEBUG_PORT}"
        "--user-data-dir=${profile}"
        --disable-extensions
        --disable-sync
        --no-first-run
        --no-default-browser-check
        --disable-features=CalculateNativeWinOcclusion
        --disable-backgrounding-occluded-windows
        --disable-renderer-backgrounding
        --disable-background-timer-throttling
        ${EDGE_DEBUG_EXTRA_ARGS[@]}
        about:blank
    )

    # Build PowerShell -ArgumentList:  'a','b','c'
    local arglist; arglist="$(printf "'%s'," "${flags[@]}")"; arglist="${arglist%,}"

    # WHY Start-Process: a `cmd start` from a UNC cwd fizzles; Start-Process reliably
    # detaches a Windows GUI process from this shell so the function returns immediately.
    # -WindowStyle Minimized: launch out of the way, never stealing focus or covering your
    # work (override with EDGE_DEBUG_WINDOWSTYLE=Normal for live animation/FPS work).
    powershell.exe -NoProfile -Command \
        "Start-Process -FilePath '${exe_win}' -WindowStyle ${EDGE_DEBUG_WINDOWSTYLE:-Minimized} -ArgumentList ${arglist}" >/dev/null 2>&1

    # Cold start can take a few seconds; poll the CDP port.
    local i
    for i in {1..20}; do
        if curl -fsS --max-time 2 "http://localhost:${EDGE_DEBUG_PORT}/json/version" >/dev/null 2>&1; then
            print -r -- "edge-debug: up -> http://localhost:${EDGE_DEBUG_PORT}  (profile: ${profile})"
            print "  reconnect Playwright MCP:  /mcp -> playwright -> Reconnect"
            return 0
        fi
        sleep 0.5
    done
    print -u2 "edge-debug: timed out waiting for port ${EDGE_DEBUG_PORT}"
    return 1
}

edge-debug-stop() {
    emulate -L zsh
    local profile; profile="$(_edge_debug_profile)"
    # Kill ONLY msedge processes whose command line references our profile dir.
    # WHY: a blanket `taskkill /im msedge.exe` would also close your personal browser.
    powershell.exe -NoProfile -Command \
        "Get-CimInstance Win32_Process -Filter \"Name='msedge.exe'\" | Where-Object { \$_.CommandLine -like '*${profile}*' } | ForEach-Object { Stop-Process -Id \$_.ProcessId -Force }" \
        >/dev/null 2>&1
    print -r -- "edge-debug: stopped (profile: ${profile})"
}

edge-debug-restart() { edge-debug-stop; sleep 1; edge-debug; }

fi

# Hook entry point — defined unconditionally so it always exists (no-op on non-WSL).
# The Claude Code PreToolUse hook for mcp__playwright__* calls this, keeping the launch
# logic here in the dotfiles instead of inline in settings.json. Always returns 0 so it
# can never block a tool call.
edge-debug-ensure() {
    emulate -L zsh
    whence -w edge-debug >/dev/null 2>&1 || return 0   # not WSL / not loaded -> no-op
    edge-debug >/dev/null 2>&1
    return 0
}
