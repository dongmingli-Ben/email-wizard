set -e

cd "$(dirname "$0")"
rm -rf client/dist/bundle.js
npm run build