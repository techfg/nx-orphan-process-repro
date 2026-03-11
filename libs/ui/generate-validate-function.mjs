import standaloneCode from "ajv/dist/standalone/index.js"
import Ajv from "ajv"
import * as esbuild from "esbuild"
import { readFile } from "node:fs/promises"

// NOTE - if NODE_ENV is undefined, esbuild will set NODE_ENV to "production" when all minify options
// are true, otherwise set it to "development".  If NODE_ENV is specifically set to development prior
// to calling esbuild, we do not want to minify.
const isDev = process.env.NODE_ENV === "development"

const contents = await readFile(
  "../../dist/ui/types/metadata/view/viewDefinition.schema.json",
)
const viewDefinitionSchema = JSON.parse(contents.toString())

const ajv = new Ajv({
  code: {
    source: true,
    esm: true,
  },
  allowUnionTypes: true,
})

const validate = ajv.compile(viewDefinitionSchema)
const moduleCode = standaloneCode(ajv, validate)

await esbuild.build({
  bundle: true,
  outfile:
    "../../dist/ui/types/metadata/view/viewDefinition.schema.validate.mjs",
  allowOverwrite: true,
  write: true,
  minify: !isDev,
  format: "esm",
  logLevel: isDev ? "debug" : "warning", // defaults to "warning" if not set https://esbuild.github.io/api/#log-level
  sourcemap: true,
  stdin: {
    contents: moduleCode,
    resolveDir: "../../",
    sourcefile: "viewDefinition.schema.validate.input.js", // Virtual source file name
    loader: "js",
  },
})
