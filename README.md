# Wails 3 + Svelte + TypeScript + Tailwind CSS

This template includes a minimal Svelte 5 frontend, Vite, Tailwind CSS 4, and a Go greeting service. Wails generates the platform build assets during initialization and regenerates the bindings during builds.

## Create an app from Git

Run this from the directory where you want the new app. Use the template repository's clone URL with `-t`.

```sh
wails3 init -n appname -t https://github.com/your-account/wails-template.git
```

Wails fetches the remote's default branch. Commit and push template changes before using the URL. HTTPS and SSH clone URLs work as well.

## Create an app from a local checkout

Export the committed template into a clean directory so Wails does not copy the checkout's `.git` directory or local build files.

```sh
repo="/path/to/wails-template"
template_dir="$(mktemp -d)/svelte-tailwind-ts"
mkdir "$template_dir"
git -C "$repo" archive HEAD | tar -xf - -C "$template_dir"
wails3 init -n appname -t "$template_dir"
```

The installed Wails v3.0.0-beta.31 CLI uses the `-ts` suffix to select TypeScript bindings for external templates. For a Git URL, you can append a pushed tag ending in `-ts` (for example, `@v0.0.1-ts`). The Svelte frontend uses TypeScript with either form.

## Work on the generated app

```sh
cd appname
wails3 dev
wails3 build
```

With `-n appname`, the app and window are named `appname` and the Go module defaults to `appname`. Use `-mod example.com/team/appname` to set the module explicitly. Use `-git <new-app-repo-url>` to initialize the generated app's Git repository and derive its module path. To set the display name in packaging metadata, also pass `-productname "App Name"`.

To run Vite directly, first install frontend dependencies and generate the Wails bindings.

```sh
(cd frontend && npm ci)
wails3 generate bindings -ts -i
(cd frontend && npm run check && npm run build)
```
