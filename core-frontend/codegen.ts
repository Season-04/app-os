import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  overwrite: true,
  schema: './schema.graphql',

  // documents: "src/**/*.gql",
  documents: ['src/**/*.tsx'],
  generates: {
    'src/gql/': {
      preset: 'client',
    },
    'src/gql/fragments.ts': {
      plugins: ['fragment-matcher'],
      config: {
        apolloClientVersion: 3,
      },
    },
  },
}

export default config
