# xb-tasks production build
# Mirrors the Taskfile.yml build pipeline using just

app_name := "team-ado-tool"
bin_dir  := "bin"

# Default recipe — build + package
default: package

# Full production build (frontend + bindings + icons + Go binary)
build: _go-mod-tidy _frontend _icons
    CGO_ENABLED=1 \
    GOOS=darwin \
    CGO_CFLAGS="-mmacosx-version-min=10.15" \
    CGO_LDFLAGS="-mmacosx-version-min=10.15" \
    MACOSX_DEPLOYMENT_TARGET="10.15" \
    go build -tags production -trimpath -buildvcs=false -ldflags="-w -s" -o {{bin_dir}}/{{app_name}}

# Create .app bundle and ad-hoc codesign
package: build
    mkdir -p "{{bin_dir}}/{{app_name}}.app/Contents/MacOS"
    mkdir -p "{{bin_dir}}/{{app_name}}.app/Contents/Resources"
    cp build/darwin/icons.icns "{{bin_dir}}/{{app_name}}.app/Contents/Resources"
    @if [ -f build/darwin/Assets.car ]; then \
        cp build/darwin/Assets.car "{{bin_dir}}/{{app_name}}.app/Contents/Resources"; \
    fi
    cp "{{bin_dir}}/{{app_name}}" "{{bin_dir}}/{{app_name}}.app/Contents/MacOS"
    cp build/darwin/Info.plist "{{bin_dir}}/{{app_name}}.app/Contents"
    codesign --force --deep --sign - "{{bin_dir}}/{{app_name}}.app"
    @echo "✓ {{bin_dir}}/{{app_name}}.app ready"

# Clean build artifacts and WebView caches
clean:
    rm -rf {{bin_dir}}/{{app_name}} {{bin_dir}}/{{app_name}}.app
    rm -rf frontend/dist frontend/bindings
    rm -rf ~/Library/WebKit/com.xboxservices.teamadotool/WebsiteData
    rm -rf ~/Library/Caches/com.xboxservices.teamadotool/WebKit

# --- internal recipes ---

_go-mod-tidy:
    go mod tidy

_frontend-deps:
    cd frontend && npm ci

_bindings: _go-mod-tidy
    wails3 generate bindings -f '-tags production -trimpath -buildvcs=false -ldflags="-w -s"' -clean=true -ts

_frontend: _frontend-deps _bindings
    rm -rf frontend/dist
    cd frontend && PRODUCTION=true npm run build -q

_icons:
    cd build && wails3 generate icons \
        -input appicon.png \
        -macfilename darwin/icons.icns \
        -windowsfilename windows/icon.ico \
        -iconcomposerinput appicon.icon \
        -macassetdir darwin
