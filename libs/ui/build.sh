#!/bin/bash

set -e

node ./build.mjs
mkdir -p ../../dist/ui/types/client
mkdir -p ../../dist/ui/types/server
echo -e "declare module \"@acme/bots\" {\n$(cat src/public_types/server/index.d.ts)\n}" > ../../dist/ui/types/server/bots.d.ts
echo -e "declare module \"@acme/bots\" {\n$(cat src/public_types/server/index.d.ts)\n}" > ../../dist/ui/types/client/index.d.ts
echo -e "declare module \"@acme/ui\" {\n$(cat src/public_types/client/index.d.ts)\n}" >> ../../dist/ui/types/client/index.d.ts

npx ts-json-schema-generator -j extended \
    --path "src/definition/definition.ts" \
    --tsconfig "./tsconfig.lib.json" \
    --type 'ViewDefinition' \
    --no-type-check \
    -o ../../dist/ui/types/metadata/view/viewDefinition.schema.json
npx ts-json-schema-generator -j extended \
    --path "src/definition/definition.ts" \
    --tsconfig "./tsconfig.lib.json" \
    --type 'ViewMetadata' \
    --no-type-check \
    -o ../../dist/ui/types/metadata/view/view.schema.json

node ./generate-validate-function.mjs
